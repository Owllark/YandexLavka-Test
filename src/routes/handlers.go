package routes

import (
	echo "github.com/labstack/echo/v4"
	"strconv"
	"yandex-team.ru/bstask/db/schemas"
)

// getCourierByID handles GET /couriers/:courier_id request
// takes echo.Context as an argument and returns response 200 with json encoded courier data
// or 400 (BadRequest) in case of error
func getCourierByID(ctx echo.Context) error {
	courierID, err := strconv.ParseInt(ctx.Param("courier_id"), 10, 64)
	courierData, err := database.GetCourierByID(courierID)
	if err != nil {
		return echo.NewHTTPError(400)
	}

	return ctx.JSON(200, courierData)

}

func getCourierMetaInfo() {

}

// getCouriers handles GET /couriers request
// takes echo.Context as an argument and returns response 200 with json encoded array of courier data
// or 400 (BadRequest) in case of error
func getCouriers(ctx echo.Context) error {
	couriersData, err := database.GetCouriers()
	if err != nil {
		return echo.NewHTTPError(400)
	}
	return ctx.JSON(200, couriersData)

}

// createCourier handles POST /couriers request
// takes echo.Context as an argument and returns response 200 with json encoded array of added courier data
// or 400 (BadRequest) in case of error
func createCourier(ctx echo.Context) error {
	var request schemas.CreateCourierRequest
	var response schemas.CreateCourierResponse
	err := ctx.Bind(&request)
	if err != nil {
		return echo.NewHTTPError(400)
	}
	for _, courier := range request.Couriers {
		inserted, err := database.InsertCourier(courier)
		if err != nil {
			return echo.NewHTTPError(400)
		}
		response.Couriers = append(response.Couriers, inserted)
	}

	return ctx.JSON(200, response)
}

// getOrders handles GET /orders request
// takes echo.Context as an argument and returns response 200 with json encoded array of order data
// or 400 (BadRequest) in case of error
func getOrders(ctx echo.Context) error {

	ordersData, err := database.GetOrders()
	if err != nil {
		return echo.NewHTTPError(400)
	}
	return ctx.JSON(200, ordersData)
}

// getOrderByID handles GET /orders/:order_id request
// takes echo.Context as an argument and returns response 200 with json encoded order data
// or 400 (BadRequest) in case of error
func getOrderByID(ctx echo.Context) error {

	orderID, err := strconv.ParseInt(ctx.Param("order_id"), 10, 64)
	orderData, err := database.GetOrderByID(orderID)
	if err != nil {
		return err
	}

	return ctx.JSON(200, orderData)
}

// createOrder handles POST /orders request
// takes echo.Context as an argument and returns response 200 with json encoded array of added order data
// or 400 (BadRequest) in case of error
func createOrder(ctx echo.Context) error {
	var request schemas.CreateOrderRequest
	var response []schemas.OrderDto
	err := ctx.Bind(&request)
	if err != nil {
		return echo.NewHTTPError(400)
	}
	for _, order := range request.Orders {
		inserted, err := database.InsertOrder(order)
		if err != nil {
			return echo.NewHTTPError(400)
		}
		response = append(response, inserted)
	}

	return ctx.JSON(200, response)
}

func completeOrder(ctx echo.Context) error {

	return nil
}

func ordersAssign() {

}

func couriersAssignment() {

}
