package helperstruct

import "github.com/Nishad4140/order_service/entitties"

type OrderAll struct {
	ProductId uint
	Quantity  float64
	Total     uint
}

type GetAllOrder struct {
	OrderId       uint
	AddressId     uint
	PaymentTypeId uint
	OrderStatusId uint
	OrderItems    []entitties.OrderItems
}