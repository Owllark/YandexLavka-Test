package db

import (
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"strconv"
	"yandex-team.ru/bstask/db/sql_db"
	"yandex-team.ru/bstask/schemas"
)

// LavkaDatabasePG provides higher level of database abstraction
// contains sql_db.PostgreSQLDatabase and provides methods for necessary queries
type LavkaDatabasePG struct {
	db sql_db.PostgreSQLDatabase
}

// Connect takes arguments user, password, dbname, host for connecting to database
// load migration files and executes migrations
// returns error of database connection
func (l *LavkaDatabasePG) Connect(user, password, dbname, host string) error {
	err := l.db.ConnectToHost(user, password, host)
	if err != nil {
		return err
	}
	err = MigrationsUp(l.db.Conn, "db/create_db_migrations")
	if err != nil {
		return err
	}
	err = l.db.Connect(user, password, dbname, host)
	if err != nil {
		return err
	}
	err = MigrationsUp(l.db.Conn, "db/migrations")
	return err
}

// GetCourierById returns courier data as schemas.CourierDto for the given courier id, and error
func (l *LavkaDatabasePG) GetCourierById(id int64) (schemas.CourierDto, error) {
	var res schemas.CourierDto
	row := l.db.QueryRow("SELECT courier_id, type, regions, working_hours FROM couriers WHERE courier_id=$1", strconv.FormatInt(id, 10))
	scanRow := rowScanner{row}
	res, err := scanRow.ScanCourierData()

	return res, err
}

// GetCouriers returns array of schemas.CourierDto and error
func (l *LavkaDatabasePG) GetCouriers() ([]schemas.CourierDto, error) {
	var res []schemas.CourierDto
	rows, err := l.db.Query("SELECT courier_id, type, regions, working_hours FROM couriers")
	for rows.Next() {
		scanRows := rowsScanner{
			rows,
		}
		c, err := scanRows.ScanCourierData()

		if err != nil {
			log.Println(err)
			continue
		}

		res = append(res, c)
	}
	err = rows.Err()
	return res, err
}

// InsertCourier Inserts into couriers table courier data from taken schemas.CreateCourierDto
// and returns inserted schemas.CourierDto and error
func (l *LavkaDatabasePG) InsertCourier(c schemas.CreateCourierDto) (schemas.CourierDto, error) {
	var res schemas.CourierDto
	switch c.CourierType {
	case "FOOT":
		if len(c.Regions) != 1 {
			return res, errors.New("invalid data")
		}
	case "BIKE":
		if len(c.Regions) != 2 {
			return res, errors.New("invalid data")
		}
	case "AUTO":
		if len(c.Regions) != 3 {
			return res, errors.New("invalid data")
		}

	}
	row := l.db.QueryRow(
		"INSERT INTO couriers (type, regions, working_hours) VALUES($1, $2, $3) RETURNING courier_id, type, regions, working_hours",
		c.CourierType,
		int32SliceToUint8Array(c.Regions),
		stringSliceToUint8Array(c.WorkingHours),
	)
	scanRow := rowScanner{row}
	res, err := scanRow.ScanCourierData()
	return res, err
}

// GetOrderById returns order data as schemas.OrderDto for the given order id, and error
func (l *LavkaDatabasePG) GetOrderById(id int64) (schemas.OrderDto, error) {
	var res schemas.OrderDto
	row := l.db.QueryRow("SELECT order_id, weight, region, delivery_hours, cost, completed_time FROM orders WHERE order_id=$1", strconv.FormatInt(id, 10))
	scanRow := rowScanner{row}
	res, err := scanRow.ScanOrderData()
	if err != nil {
		return res, err
	}

	return res, err
}

// GetOrders returns all couriers data from table as slice of schemas.OrderDto
func (l *LavkaDatabasePG) GetOrders() ([]schemas.OrderDto, error) {
	var res []schemas.OrderDto
	rows, err := l.db.Query(fmt.Sprintf("SELECT order_id, weight, region, delivery_hours, cost, completed_time FROM orders"))
	for rows.Next() {
		scanRows := rowsScanner{
			rows,
		}
		order, err := scanRows.ScanOrderData()

		if err != nil {
			continue
		}
		res = append(res, order)
	}
	err = rows.Err()
	if err != nil {
		return res, err
	}
	return res, err
}

// InsertOrder Inserts into orders table data from taken schemas.CreateOrderDto
// and returns inserted schemas.CourierDto and error
func (l *LavkaDatabasePG) InsertOrder(order schemas.CreateOrderDto) (schemas.OrderDto, error) {
	row := l.db.QueryRow(
		"INSERT INTO orders (weight, region, delivery_hours, cost) VALUES($1, $2, $3, $4) RETURNING order_id, weight, region, delivery_hours, cost, completed_time",
		order.Weight,
		order.Region,
		stringSliceToUint8Array(order.DeliveryHours),
		order.Cost,
	)
	scanRow := rowScanner{row}
	res, err := scanRow.ScanOrderData()
	return res, err
}

// DeleteCompletedOrder deletes order by id from completed_orders table
func (l *LavkaDatabasePG) DeleteCompletedOrder(id int64) error {
	_, err := l.db.Exec("DELETE FROM completed_orders WHERE order_id = $1", id)
	return err
}

// SetOrderCompleteTime sets order with given id completed_time equal to given time
func (l *LavkaDatabasePG) SetOrderCompleteTime(id int64, time string) error {
	_, err := l.db.Exec("UPDATE orders SET completed_time = $1 WHERE order_id = $2", []byte(time), id)
	return err
}

// InsertCompletedOrder Inserts into completed_orders table data from taken schemas.CompleteOrder
func (l *LavkaDatabasePG) InsertCompletedOrder(complete schemas.CompleteOrder) error {
	row := l.db.QueryRow(
		"INSERT INTO completed_orders (courier_id, order_id, complete_time) VALUES($1, $2, $3)",
		complete.CourierId,
		complete.OrderID,
		[]byte(complete.CompleteTime),
	)
	err := row.Scan()
	return err
}

// CountCourierEarnings returns earnings of courier with given id from startDate to endDate
func (l *LavkaDatabasePG) CountCourierEarnings(id int64, startDate string, endDate string) (int32, error) {
	var res int32
	row := l.db.QueryRow("SELECT SUM(orders.cost) * \n"+
		"(SELECT CASE WHEN type = 'FOOT' THEN 2\n"+
		"WHEN type = 'BIKE' THEN 3\n"+
		"WHEN type = 'AUTO' THEN 4\n"+
		"ELSE -1\n"+
		"END\n"+
		"FROM couriers\n"+
		"WHERE courier_id = $1)\n"+
		"AS earnings\n"+
		"FROM completed_orders\n"+
		"JOIN orders ON completed_orders.order_id = orders.order_id\n"+
		"WHERE completed_orders.courier_id = $1\n"+
		"AND completed_orders.complete_time >= $2\n"+
		"AND completed_orders.complete_time < $3;",
		id,
		[]byte(startDate),
		[]byte(endDate),
	)
	err := row.Scan(&res)
	// Here must be more smart error handling to separate database error and error when reading NULL value
	if err != nil {
		res = 0
		err = nil
	}
	return res, err
}

// CountCourierRating returns rating of courier with given id from startDate to endDate
func (l *LavkaDatabasePG) CountCourierRating(id int64, startDate string, endDate string) (int32, error) {
	var res float32
	row := l.db.QueryRow("SELECT (COUNT(*)::FLOAT / EXTRACT(EPOCH FROM ($3::timestamp - $2::timestamp)) / 3600) * \n"+
		"(SELECT CASE WHEN type = 'FOOT' THEN 2\n"+
		"WHEN type = 'BIKE' THEN 3\n"+
		"WHEN type = 'AUTO' THEN 4\n"+
		"ELSE -1\n\t\t\t\tEND\n"+
		"FROM couriers\n"+
		"WHERE courier_id = $1)\n"+
		"AS rating\n"+
		"FROM completed_orders\n"+
		"JOIN couriers ON completed_orders.courier_id = couriers.id\n"+
		"WHERE completed_orders.courier_id = $1\n"+
		"AND completed_orders.complete_time >= $2\n"+
		"AND completed_orders.complete_time <= $3;",
		id,
		[]byte(startDate),
		[]byte(endDate),
	)
	err := row.Scan(&res)
	// Here must be more smart error handling to separate database error and error when reading NULL value
	if err != nil {
		res = 0
		err = nil
	}
	return int32(res), err
}
