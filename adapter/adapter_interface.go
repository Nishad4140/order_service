package adapter

import helperstruct "github.com/Nishad4140/order_service/helper_struct"

type AdapterInterface interface {
	OrderAll(items []helperstruct.OrderAll, userId uint) (int, error)
	CancelOrder(orderId uint) error
}
