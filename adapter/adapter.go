package adapter

import (
	"fmt"

	helperstruct "github.com/Nishad4140/order_service/helper_struct"
	"gorm.io/gorm"
)

type OrderAdapter struct {
	DB *gorm.DB
}

func NewOrderAdapter(db *gorm.DB) *OrderAdapter {
	return &OrderAdapter{
		DB: db,
	}
}

func (order *OrderAdapter) OrderAll(items []helperstruct.OrderAll, userId uint) (int, error) {
	var orderId int
	tx := order.DB.Begin()

	query := "INSERT INTO orders (user_id, payment_type, total) VALUES ($1, $2, 0) RETURNING id"
	if err := tx.Raw(query, userId, 1).Scan(&orderId).Error; err != nil {
		tx.Rollback()
		return -1, err
	}
	if orderId == 0 {
		return -1, fmt.Errorf("order not found")
	}
	for _, item := range items {
		queryItemsInsert := "INSERT INTO order_items (product_id, quantity, total, order_id) VALUES ($1, $2, $3, $4)"
		if err := tx.Exec(queryItemsInsert, item.ProductId, item.Quantity, item.Total, orderId).Error; err != nil {
			tx.Rollback()
			return -1, err
		}
		queryUpdateTotal := "UPDATE orders SET total = total + $1 WHERE id = $2"
		if err := tx.Exec(queryUpdateTotal, item.Total, orderId).Error; err != nil {
			tx.Rollback()
			return -1, err
		}
	}
	if err := tx.Commit().Error; err != nil {
		return -1, fmt.Errorf("error while transaction")
	}
	return orderId, nil
}

func (order *OrderAdapter) CancelOrder(orderId uint) error {
	tx := order.DB.Begin()
	queryDelete := "DELETE FROM order_items WHERE order_id = ?"
	if err := tx.Exec(queryDelete, orderId).Error; err != nil {
		tx.Rollback()
		return err
	}
	deleteOrder := "UPDATE orders SET order_status_id = $1 WHERE id = $2"
	if err := tx.Exec(deleteOrder, 5, orderId).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}
