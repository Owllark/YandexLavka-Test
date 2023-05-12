package tests

import (
	"bytes"
	"encoding/json"
	"example.com/src/db/schemas"
	"fmt"
	"github.com/stretchr/testify/require"
	"io"
	"math/rand"
	"net/http"
	"testing"
)

func TestRespondsWithLove(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/couriers"))
	require.NoError(t, err, "HTTP error")
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode, "HTTP status code")

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "failed to read HTTP body")

	// Finally, test the business requirement!
	require.Equal(t, "pong", string(body), "Wrong ping response")
}

func TestCreateCourier(t *testing.T) {

	n := 10
	for i := 0; i < n; i++ {
		testCourier := generateCourier()
		var request schemas.CreateCourierRequest
		request.Couriers = append(request.Couriers, testCourier)

		jsonBody, err := json.Marshal(request)
		if err != nil {
			t.Error("Failed to encode request body:", err)
			continue
		}
		fmt.Println(string(jsonBody))
		resp, err := http.Post(fmt.Sprintf("http://localhost:8080/couriers"), "application/json", bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Error("HTTP error")
			continue
		}

		if resp.StatusCode != http.StatusOK {
			t.Error("HTTP status code ", resp.StatusCode)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Error("failed to read http body")
			continue
		}

		var response schemas.CreateCourierResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			t.Error("failed unmarshall json body")
			continue
		}
		var responseData []schemas.CreateCourierDto
		for j := 0; j < len(response.Couriers); j++ {
			el := response.Couriers[j]
			courier := schemas.CreateCourierDto{
				CourierType:  el.CourierType,
				Regions:      el.Regions,
				WorkingHours: el.WorkingHours,
			}
			responseData = append(responseData, courier)
		}
		jsonResponse, err := json.Marshal(responseData)
		jsonExpected, _ := json.Marshal(request.Couriers)
		if string(jsonResponse) != string(jsonExpected) {
			t.Error("failed to read http body")
			continue
		}
	}

}

func generateCourier() schemas.CreateCourierDto {
	var res schemas.CreateCourierDto
	switch 1 + rand.Int()%3 {
	case 1:
		res.CourierType = "FOOT"
		res.Regions = generateRegions(1)
	case 2:
		res.CourierType = "BIKE"
		res.Regions = generateRegions(2)
	case 3:
		res.CourierType = "AUTO"
		res.Regions = generateRegions(3)
	}
	startTime := 6 + rand.Int()%11
	endTime := startTime + (1 + rand.Int()%8)
	res.WorkingHours = append(res.WorkingHours, fmt.Sprintf("%d:00", startTime))
	res.WorkingHours = append(res.WorkingHours, fmt.Sprintf("%d:00", endTime))
	return res
}

func generateRegions(n int) []int32 {
	var res []int32
	for i := 0; i < n; i++ {
		res = append(res, int32(rand.Int()%10))
	}
	return res
}

func TestCreateOrder(t *testing.T) {

	n := 10
	for i := 0; i < n; i++ {
		testOrder := generateOrder()
		var request schemas.CreateOrderRequest
		request.Orders = append(request.Orders, testOrder)

		jsonBody, err := json.Marshal(request)
		if err != nil {
			t.Error("Failed to encode request body:", err)
			continue
		}
		fmt.Println(string(jsonBody))
		resp, err := http.Post(fmt.Sprintf("http://localhost:8080/orders"), "application/json", bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Error("HTTP error")
			continue
		}

		if resp.StatusCode != http.StatusOK {
			t.Error("HTTP status code ", resp.StatusCode)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Error("failed to read http body")
			continue
		}

		var response []schemas.OrderDto
		err = json.Unmarshal(body, &response)
		if err != nil {
			t.Error("failed to unmarshall json body")
			continue
		}
		var responseData []schemas.CreateOrderDto
		for j := 0; j < len(response); j++ {
			el := response[j]
			order := schemas.CreateOrderDto{
				Weight:        el.Weight,
				Region:        el.Region,
				DeliveryHours: el.DeliveryHours,
				Cost:          el.Cost,
			}
			responseData = append(responseData, order)
		}
		jsonResponse, _ := json.Marshal(responseData)
		jsonExpected, _ := json.Marshal(request.Orders)
		if string(jsonResponse) != string(jsonExpected) {
			t.Error("failed to read http body")
			continue
		}
	}

}

func generateOrder() schemas.CreateOrderDto {
	var res schemas.CreateOrderDto
	res.Cost = int32(50 + (rand.Int()%20)*50)
	res.Weight = 50 * float32(1+rand.Int()%20) / 20
	res.Region = int32(rand.Int() % 10)
	startTime := 6 + rand.Int()%15
	endTime := startTime + (1 + rand.Int()%3)
	res.DeliveryHours = append(res.DeliveryHours, fmt.Sprintf("%d:00", startTime))
	res.DeliveryHours = append(res.DeliveryHours, fmt.Sprintf("%d:00", endTime))
	return res
}
