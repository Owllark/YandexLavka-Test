package schemas

type CreateOrderDto struct {
	Weight        float32  `json:"weight,omitempty"`
	Region        int32    `json:"region,omitempty"`
	DeliveryHours []string `json:"delivery___hours,omitempty"`
	Cost          int32    `json:"cost,omitempty"`
}

type CreateOrderRequest struct {
	Orders []CreateOrderDto `json:"orders,omitempty"`
}

type OrderDto struct {
	OrderId       int64    `json:"order_id"`
	Weight        float32  `json:"weight,omitempty"`
	Region        int32    `json:"region,omitempty"`
	DeliveryHours []string `json:"delivery___hours,omitempty"`
	Cost          int32    `json:"cost,omitempty"`
	CompletedTime string   `json:"completed___time,omitempty"`
}

type GroupOrders struct {
	GroupOrderId int64      `json:"group_order_id,omitempty"`
	Orders       []OrderDto `json:"orders,omitempty"`
}

type CouriersGroupOrders struct {
	CourierId int64         `json:"courier_id,omitempty"`
	Orders    []GroupOrders `json:"orders,omitempty"`
}

type OrderAssignResponse struct {
	Date     string              `json:"date,omitempty"`
	Couriers CouriersGroupOrders `json:"couriers"`
}

type BadRequestResponse struct {
}

type CompleteOrder struct {
	CourierId    int64  `json:"courier_id,omitempty"`
	OrderID      int64  `json:"order_id,omitempty"`
	CompleteTime string `json:"complete_time,omitempty"`
}

type CompleteOrderRequestDto struct {
	CompleteInfo []CompleteOrder `json:"complete_info,omitempty"`
}

const (
	FOOT = "FOOT"
	BIKE = "BIKE"
	AUTO = "AUTO"
)

type CreateCourierDto struct {
	CourierType  string   `json:"courier_type,omitempty"`
	Regions      []int32  `json:"regions,omitempty"`
	WorkingHours []string `json:"working_hours,omitempty"`
}

type CreateCourierRequest struct {
	Couriers []CreateCourierDto `json:"couriers,omitempty"`
}

type CreateCourierResponse struct {
	Couriers []CourierDto `json:"couriers,omitempty"`
}

type CourierDto struct {
	CourierId    int64    `json:"courier_id,omitempty"`
	CourierType  string   `json:"courier_type,omitempty"`
	Regions      []int32  `json:"regions,omitempty"`
	WorkingHours []string `json:"working_hours,omitempty"`
}

type NotFoundResponse struct {
}

type GetCouriersResponse struct {
	Couriers []CourierDto `json:"couriers,omitempty"`
	Limit    int32        `json:"limit,omitempty"`
	Offset   int32        `json:"offset,omitempty"`
}

type GetCourierMetaInfoResponse struct {
	CourierId    int64    `json:"courier_id,omitempty"`
	CourierType  string   `json:"courier_type,omitempty"`
	Regions      []int32  `json:"regions,omitempty"`
	WorkingHours []string `json:"working_hours,omitempty"`
	Rating       int32    `json:"rating,omitempty"`
	Earnings     int32    `json:"earnings,omitempty"`
}
