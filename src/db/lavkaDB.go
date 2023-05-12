package db

import (
	"fmt"
	"log"
	"strconv"
	"yandex-team.ru/bstask/db/schemas"
	"yandex-team.ru/bstask/db/sql_db"
)

// LavkaDatabase provides higher level of database abstraction
// contains sql_db.PostgreSQLDatabase and provides methods for necessary queries
type LavkaDatabase struct {
	db sql_db.PostgreSQLDatabase
}

// Connect takes arguments user, password, dbname, host for connecting to database
// returns error of database connection
func (l *LavkaDatabase) Connect(user, password, dbname, host string) error {
	err := l.db.Connect(user, password, dbname, host)
	return err
}

// GetCourierById returns courier data as schemas.CourierDto for the given courier id, and error
func (l *LavkaDatabase) GetCourierById(id int64) (schemas.CourierDto, error) {
	var res schemas.CourierDto
	row := l.db.QueryRow("SELECT courier_id, type, regions, working_hours FROM couriers WHERE courier_id=$1", strconv.FormatInt(id, 10))
	scanRow := rowScanner{row}
	res, err := scanRow.ScanCourierData()
	if err != nil {
		return res, err
	}

	return res, err
}

// GetCouriers returns array of schemas.CourierDto and error
func (l *LavkaDatabase) GetCouriers() ([]schemas.CourierDto, error) {
	var res []schemas.CourierDto
	rows, err := l.db.Query("SELECT courier_id, type, regions, working_hours FROM couriers")
	for rows.Next() {
		scanRows := rowsScanner{
			rows,
		}
		c, err := scanRows.ScanCourierData()

		if err != nil {
			continue
		}

		res = append(res, c)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return res, err
}

// InsertCourier Inserts into couriers table courier data from taken schemas.CreateCourierDto
// and returns inserted schemas.CourierDto and error
func (l *LavkaDatabase) InsertCourier(c schemas.CreateCourierDto) (schemas.CourierDto, error) {
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
func (l *LavkaDatabase) GetOrderById(id int64) (schemas.OrderDto, error) {
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
func (l *LavkaDatabase) GetOrders() ([]schemas.OrderDto, error) {
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
func (l *LavkaDatabase) InsertOrder(order schemas.CreateOrderDto) (schemas.OrderDto, error) {
	row := l.db.QueryRow(
		"INSERT INTO orders (weight, region, delivery_hours, cost) VALUES($1, $2, $3, $4) RETURNING order_id, weight, region, delivery_hours, cost",
		order.Weight,
		order.Region,
		stringSliceToUint8Array(order.DeliveryHours),
		order.Cost,
	)
	scanRow := rowScanner{row}
	res, err := scanRow.ScanOrderData()
	return res, err
}

func (l *LavkaDatabase) SetOrderCompleteTime(id int64, time string) error {
	_, err := l.db.Exec("UPDATE orders SET completed_time = $1 WHERE order_id = $2", []byte(time), id)
	return err
}

func (l *LavkaDatabase) InsertCompletedOrder(complete schemas.CompleteOrder) error {
	row := l.db.QueryRow(
		"INSERT INTO completed_orders (courier_id, order_id, complete_time) VALUES($1, $2, $3)",
		complete.CourierId,
		complete.OrderID,
		[]byte(complete.CompleteTime),
	)
	err := row.Scan()
	return err
}

func (l *LavkaDatabase) CountCourierEarnings(id int64, startDate string, endDate string) (int32, error) {
	var res int32
	row := l.db.QueryRow("SELECT SUM(orders.cost) * \n"+
		"(SELECT CASE WHEN type = 'FOOT' THEN 2\n"+
		"WHEN type = 'BIKE' THEN 3\n"+
		"WHEN type = 'AUTO' THEN 4\n"+
		"ELSE -1\n"+
		"END\n"+
		"FROM couriers\n"+
		"WHERE courier_id = 1\n"+
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

func (l *LavkaDatabase) CountCourierRating(id int64, startDate string, endDate string) (int32, error) {
	var res float32
	row := l.db.QueryRow("SELECT (COUNT(*)::FLOAT / EXTRACT(EPOCH FROM (MAX(complete_time) - MIN(complete_time))) / 3600) * \n"+
		"CASE WHEN type = 'FOOT' THEN 3\n"+
		"WHEN type = 'BIKE' THEN 2\n"+
		"WHEN type = 'AUTO' THEN 1\n"+
		"ELSE -1\n"+
		"END AS rating\n"+
		"FROM completed_orders\n"+
		"JOIN couriers ON completed_orders.courier_id = couriers.id\n"+
		"WHERE completed_orders.courier_id = $1\n"+
		"WHERE $2 >= 'start_date' AND $3 < 'end_date';",
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
