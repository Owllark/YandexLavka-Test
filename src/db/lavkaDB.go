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

// GetCourierByID returns courier data as schemas.CourierDto for the given courier id, and error
func (l *LavkaDatabase) GetCourierByID(id int64) (schemas.CourierDto, error) {
	var res schemas.CourierDto
	row := l.db.QueryRow("SELECT * FROM couriers WHERE courier_id=$1", strconv.FormatInt(id, 10))
	scanRow := rowScanner{row}
	res, err := scanRow.ScanCourierData()
	if err != nil {
		log.Fatal(err)
	}

	return res, err
}

// GetCouriers returns array of schemas.CourierDto and error
func (l *LavkaDatabase) GetCouriers() ([]schemas.CourierDto, error) {
	var res []schemas.CourierDto
	rows, err := l.db.Query("SELECT * FROM couriers")
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
		"INSERT INTO couriers (type, regions, working_hours) VALUES($1, $2, $3) RETURNING id, type, regions, working_hours",
		c.CourierType,
		int32SliceToUint8Array(c.Regions),
		stringSliceToUint8Array(c.WorkingHours),
	)
	scanRow := rowScanner{row}
	res, err := scanRow.ScanCourierData()
	return res, err
}

// GetOrderByID returns order data as schemas.OrderDto for the given order id, and error
func (l *LavkaDatabase) GetOrderByID(id int64) (schemas.OrderDto, error) {
	var res schemas.OrderDto
	row := l.db.QueryRow("SELECT * FROM orders WHERE order_id=$1", strconv.FormatInt(id, 10))
	scanRow := rowScanner{row}
	res, err := scanRow.ScanOrderData()
	if err != nil {
		log.Fatal(err)
	}

	return res, err
}

// GetOrders returns all couriers data from table as slice of schemas.OrderDto
func (l *LavkaDatabase) GetOrders() ([]schemas.OrderDto, error) {
	var res []schemas.OrderDto
	rows, err := l.db.Query(fmt.Sprintf("SELECT * FROM orders"))
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
		log.Fatal(err)
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

/*

func (l *LavkaDatabase) MarkOrderAsCompleted() error {

}*/
