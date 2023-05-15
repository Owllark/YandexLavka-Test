package db

import "yandex-team.ru/bstask/schemas"

type LavkaDatabase interface {
	Connect(user, password, dbname, host string) error
	GetCourierById(id int64) (schemas.CourierDto, error)
	GetCouriers() ([]schemas.CourierDto, error)
	InsertCourier(c schemas.CreateCourierDto) (schemas.CourierDto, error)
	GetOrderById(id int64) (schemas.OrderDto, error)
	GetOrders() ([]schemas.OrderDto, error)
	InsertOrder(order schemas.CreateOrderDto) (schemas.OrderDto, error)
	DeleteCompletedOrder(id int64) error
	SetOrderCompleteTime(id int64, time string) error
	InsertCompletedOrder(complete schemas.CompleteOrder) error
	CountCourierEarnings(id int64, startDate, endDate string) (int32, error)
	CountCourierRating(id int64, startDate, endDate string) (int32, error)
}
