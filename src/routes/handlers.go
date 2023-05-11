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
	offsetStr := ctx.QueryParam("offset")
	limitStr := ctx.QueryParam("limit")
	var offset, limit int
	var err error
	if offsetStr == "" {
		offset = 0
	} else {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			return echo.NewHTTPError(400)
		}
	}
	if limitStr == "" {
		limit = 1
	} else {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 0 {
			return echo.NewHTTPError(400)
		}
	}

	couriersData, err := database.GetCouriers()
	if err != nil {
		return echo.NewHTTPError(400)
	}
	if offset >= len(couriersData) {
		return ctx.NoContent(200)
	}
	if offset+limit > len(couriersData) {
		return ctx.JSON(200, couriersData[offset:])
	}
	return ctx.JSON(200, couriersData[offset:offset+limit])

}

// createCourier handles POST /couriers request
// takes echo.Context as an argument and returns response 200 with json encoded array of added courier data
// or 400 (BadRequest) in case of error
func createCourier(ctx echo.Context) error {
	var request schemas.CreateCourierRequest
	var response schemas.CreateCourierResponse
	err := ctx.Bind(&request.Couriers)
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
	offsetStr := ctx.QueryParam("offset")
	limitStr := ctx.QueryParam("limit")
	var offset, limit int
	var err error
	if offsetStr == "" {
		offset = 0
	} else {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			return echo.NewHTTPError(400)
		}
	}
	if limitStr == "" {
		limit = 1
	} else {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 0 {
			return echo.NewHTTPError(400)
		}
	}
	ordersData, err := database.GetOrders()
	if err != nil {
		return echo.NewHTTPError(400)
	}

	if err != nil {
		return echo.NewHTTPError(400)
	}
	if offset >= len(ordersData) {
		return ctx.NoContent(200)
	}
	if offset+limit > len(ordersData) {
		return ctx.JSON(200, ordersData[offset:])
	}
	return ctx.JSON(200, ordersData[offset:offset+limit])
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
	err := ctx.Bind(&request.Orders)
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
	var request schemas.CompleteOrderRequestDto
	var response []schemas.OrderDto
	err := ctx.Bind(&request.CompleteInfo)
	if err != nil {
		return echo.NewHTTPError(400)
	}
	for _, complete := range request.CompleteInfo {
		err := database.SetOrderCompleteTime(complete.OrderID, complete.CompleteTime)
		if err != nil {
			return echo.NewHTTPError(400)
		}
		order, _ := database.GetOrderByID(complete.OrderID)
		response = append(response, order)
		database.InsertCompletedOrder(complete)
	}
	return ctx.JSON(200, response)
}

func ordersAssign() {

}

func couriersAssignment() {

}
