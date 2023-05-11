package routes

import (
	"github.com/labstack/echo/v4"
	"yandex-team.ru/bstask/db"
)

var (
	database = new(db.LavkaDatabase)
)

func ConnectWithDataBase(user, password, dbname, host string) error {
	err := database.Connect(user, password, dbname, host)
	return err
}

func SetupRoutes(e *echo.Echo) {

	e.GET("/ping", ping)
	e.POST("/couriers", createCourier)
	e.GET("/couriers/:courier_id", getCourierByID)
	e.GET("/couriers", getCouriers)
	e.GET("/orders", getOrders)
	e.POST("/orders", createOrder)
	e.GET("/orders/:order_id", getOrderByID)
	e.POST("/orders/complete", completeOrder)
}
