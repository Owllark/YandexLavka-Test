package routes

import (
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/labstack/echo/v4"
	"time"
	"yandex-team.ru/bstask/db"
)

var (
	database = new(db.LavkaDatabasePG)
)

func ConnectWithDataBase(user, password, dbname, host string) error {
	err := database.Connect(user, password, dbname, host)
	return err
}

func SetupRoutes(e *echo.Echo) {

	// Create a rate newLimiter with a restriction of 10 requests per second (RPS)
	var newLimiter = tollbooth.NewLimiter(10, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Second}) // Apply rate newLimiter middleware to all routes
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			httpError := tollbooth.LimitByRequest(newLimiter, c.Response(), c.Request())
			if httpError != nil {
				return echo.NewHTTPError(429)
			}
			return next(c)
		}
	})

	e.POST("/couriers", createCourier)
	e.GET("/couriers/:courier_id", getCourierById)
	e.GET("/couriers", getCouriers)
	e.GET("/orders", getOrders)
	e.POST("/orders", createOrder)
	e.GET("/orders/:order_id", getOrderById)
	e.POST("/orders/complete", completeOrder)
	e.GET("/couriers/meta-info/:courier_id", getCourierMetaInfo)
}

func shouldExcludeRoute(path string) bool {
	return false
}
