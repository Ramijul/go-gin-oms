package orderPackage

import (
	"github.com/Ramijul/go-gin-oms/orders/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	OrdersTableName     = "orders"
	OrderItemsTableName = "order_details"
)

type Repository struct {
	Session *gorm.DB
}

type OrderRepository interface {
	GetAll() (orders []*models.Order, err error)
	GetOne(id uuid.UUID) (order *models.OrderWithItems, err error)
	CreateOrder(order *models.Order) (id uuid.UUID, err error)
	CreateOrderItems(orderItems []*models.OrderItem) (success bool, err error)
}

func (r *Repository) GetAll() (orders []*models.Order, err error) {
	// Fetching multiple order with order details would be expensive. Returning
	// just orders should suffice.
	// Alternative option is to write custom query that joins the two tables, and
	// use resultset extractor to generate an array of models.OrderWithDetails

	runner := r.Session.Table(OrdersTableName)

	// TODO: add pagination
	result := runner.Find(&orders)

	if result.Error != nil {
		return nil, result.Error
	}

	return orders, nil
}

func (r *Repository) GetOne(id uuid.UUID) (orderWithDetails *models.OrderWithItems, err error) {
	// get order
	runner := r.Session.Table(OrdersTableName)
	order := models.Order{}
	result := runner.First(&order, id)
	if result.Error != nil {
		return nil, result.Error
	}

	// get order details
	runner = r.Session.Table(OrderItemsTableName)
	var orderItems []*models.OrderItem
	result = runner.Where("order_id = ?", order.ID).Find(&orderItems)
	if result.Error != nil {
		return nil, result.Error
	}

	// aggregate resutls into one object
	return createOrderWithDetails(&order, orderItems), nil

}

func (r *Repository) CreateOrder(order *models.Order) (id uuid.UUID, err error) {
	// insert order
	runner := r.Session.Table(OrdersTableName)
	result := runner.Create(&order)
	if result.Error != nil {
		return uuid.Nil, result.Error
	}

	return order.ID, nil
}

func (r *Repository) CreateOrderItems(orderItems []*models.OrderItem) (success bool, err error) {
	// insert order items
	runner := r.Session.Table(OrderItemsTableName)
	result := runner.Create(&orderItems)
	if result.Error != nil {
		return false, result.Error
	}

	return true, nil
}

/*
Aggregate models.Order and models.OrderDetails into one model
*/
func createOrderWithDetails(order *models.Order, orderItems []*models.OrderItem) *models.OrderWithItems {
	return &models.OrderWithItems{
		Order: models.Order{
			Model: gorm.Model{
				CreatedAt: order.CreatedAt,
				UpdatedAt: order.UpdatedAt,
				DeletedAt: order.DeletedAt},
			ID:                order.ID,
			TotalPrice:        order.TotalPrice,
			PaymentReceivedAt: order.PaymentReceivedAt,
			OrderStatus:       order.OrderStatus,
			PaymentStatus:     order.PaymentStatus,
			UserID:            order.UserID,
			UserName:          order.UserName,
			UserEmail:         order.UserEmail,
			UserPhoneNumber:   order.UserPhoneNumber,
			Address:           models.Address(order.Address),
		},
		OrderItems: orderItems,
	}
}
