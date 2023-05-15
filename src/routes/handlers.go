package routes

import (
	echo "github.com/labstack/echo/v4"
	"strconv"
	"yandex-team.ru/bstask/schemas"
)

// getCourierById handles GET /couriers/:courier_id request
// takes echo.Context as an argument and returns response 200 with json encoded courier data
// or 400 (BadRequest) in case of error
func getCourierById(ctx echo.Context) error {
	courierID, err := strconv.ParseInt(ctx.Param("courier_id"), 10, 64)
	courierData, err := database.GetCourierById(courierID)
	if err != nil {
		return echo.NewHTTPError(400)
	}

	return ctx.JSON(200, courierData)

}

// getCouriers handles GET /couriers request
// takes echo.Context as an argument and returns response 200 with json encoded array of courier data
// or 400 (BadRequest) in case of error
func getCouriers(ctx echo.Context) error {
	var res schemas.GetCouriersResponse
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
	res.Limit = int32(limit)
	res.Offset = int32(offset)
	if offset >= len(couriersData) {
		res.Couriers = nil
	} else if offset+limit >= len(couriersData) {
		res.Couriers = couriersData[offset:]
	} else {
		res.Couriers = couriersData[offset : offset+limit]
	}
	return ctx.JSON(200, res)

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

// getOrderById handles GET /orders/:order_id request
// takes echo.Context as an argument and returns response 200 with json encoded order data
// or 400 (BadRequest) in case of error
func getOrderById(ctx echo.Context) error {

	orderID, err := strconv.ParseInt(ctx.Param("order_id"), 10, 64)
	orderData, err := database.GetOrderById(orderID)
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

// completeOrder handles POST /orders/complete request
// takes echo.Context as an argument and returns response 200 with json encoded array of completed orders
// or 400 (BadRequest) in case of error
func completeOrder(ctx echo.Context) error {
	var request schemas.CompleteOrderRequestDto
	var response []schemas.OrderDto
	err := ctx.Bind(&request)
	if err != nil {
		return echo.NewHTTPError(400)
	}
	for _, complete := range request.CompleteInfo {
		_, err := database.GetOrderById(complete.OrderID)
		if err != nil {
			return echo.NewHTTPError(400)
		}
		_, err = database.GetCourierById(complete.CourierId)
		if err != nil {
			return echo.NewHTTPError(400)
		}
	}
	for _, complete := range request.CompleteInfo {

		err = database.DeleteCompletedOrder(complete.OrderID)
		if err != nil {
			return echo.NewHTTPError(400)
		}
		err = database.SetOrderCompleteTime(complete.OrderID, complete.CompleteTime)
		if err != nil {
			return echo.NewHTTPError(400)
		}
		order, _ := database.GetOrderById(complete.OrderID)
		response = append(response, order)
		database.InsertCompletedOrder(complete)
	}
	return ctx.JSON(200, response)
}

// completeOrder handles GET /couriers/meta-info/:courier_id request
// takes echo.Context as an argument and returns response 200 with json encoded data about courier and his earnings and ratings
// or 400 (BadRequest) in case of error
func getCourierMetaInfo(ctx echo.Context) error {
	var response schemas.GetCourierMetaInfoResponse
	startDate := ctx.QueryParam("start_date")
	endDate := ctx.QueryParam("end_date")
	id, err := strconv.ParseInt(ctx.Param("courier_id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(400)
	}
	courier, err := database.GetCourierById(id)
	if err != nil {
		return echo.NewHTTPError(400)
	}
	earnings, err := database.CountCourierEarnings(id, startDate, endDate)
	if err != nil {
		return echo.NewHTTPError(400)
	}
	rating, err := database.CountCourierRating(id, startDate, endDate)
	if err != nil {
		return echo.NewHTTPError(400)
	}
	response.CourierId = courier.CourierId
	response.CourierType = courier.CourierType
	response.WorkingHours = courier.WorkingHours
	response.Regions = courier.Regions
	response.Earnings = earnings
	response.Rating = rating

	return ctx.JSON(200, response)

}
