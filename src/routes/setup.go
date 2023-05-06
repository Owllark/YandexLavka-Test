package routes

import (
	"github.com/labstack/echo/v4"
)

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
