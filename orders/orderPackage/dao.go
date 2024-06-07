package orderPackage

import (
	"time"

	"github.com/Ramijul/go-gin-oms/orders/models"
)

type address struct {
	Street   string `json:"street" binding:"required"`
	Zip      string `json:"zip" binding:"required"`
	City     string `json:"city" binding:"required"`
	Province string `json:"province" binding:"required"`
}

type RequestedItem struct {
	ID    string `json:"id" binding:"required,uuid4"`
	Units int    `json:"units" binding:"required,min=1"`
}

type CreateRequestDao struct {
	UserID   string           `json:"user_id" binding:"required,uuid4"`
	Address  address          `json:"address" binding:"required"`
	Products []*RequestedItem `json:"products" binding:"required,gt=0"`
}

type RequestedByUser struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

type OrderItemsResposeDao struct {
	ProductID        string  `json:"product_id"`
	ProductName      string  `json:"product_name"`
	ProductUnitPrice float64 `json:"product_unit_price"`
	Units            int     `json:"product_units"`
}

type OrderResponseDao struct {
	ID                string                  `json:"id"`
	TotalPrice        float64                 `json:"total_price"`
	PaymentReceivedAt time.Time               `json:"payment_received_at"`
	OrderStatus       string                  `json:"order_status"`
	PaymentStatus     string                  `json:"payment_status"`
	User              RequestedByUser         `json:"user"`
	Address           address                 `json:"address"`
	CreatedAt         time.Time               `json:"created_at"`
	OrderItems        []*OrderItemsResposeDao `json:"order_items,omitempty"`
}

type ManyOrdersResponseDao struct {
	Orders []*OrderResponseDao `json:"orders"`
}

func ToOrderResponseDao(o *models.OrderWithItems) *OrderResponseDao {
	var orderItems []*OrderItemsResposeDao
	for _, item := range o.OrderItems {
		orderItems = append(orderItems, &OrderItemsResposeDao{
			ProductID:        item.ProductID.String(),
			ProductName:      item.ProductName,
			ProductUnitPrice: item.ProductUnitPrice,
			Units:            item.Units,
		})
	}

	return &OrderResponseDao{
		ID:                o.ID.String(),
		TotalPrice:        o.TotalPrice,
		PaymentReceivedAt: o.PaymentReceivedAt,
		OrderStatus:       o.OrderStatus,
		PaymentStatus:     o.PaymentStatus,
		User: RequestedByUser{
			ID:          o.UserID.String(),
			Name:        o.UserName,
			Email:       o.UserEmail,
			PhoneNumber: o.UserPhoneNumber,
		},
		Address:    address(o.Address),
		CreatedAt:  o.CreatedAt,
		OrderItems: orderItems,
	}
}

func ToManyOrdersResponseDao(o []*models.Order) *ManyOrdersResponseDao {
	var orders []*OrderResponseDao
	for _, item := range o {
		orders = append(orders, &OrderResponseDao{
			ID:                item.ID.String(),
			TotalPrice:        item.TotalPrice,
			PaymentReceivedAt: item.PaymentReceivedAt,
			OrderStatus:       item.OrderStatus,
			PaymentStatus:     item.PaymentStatus,
			User: RequestedByUser{
				ID:          item.UserID.String(),
				Name:        item.UserName,
				Email:       item.UserEmail,
				PhoneNumber: item.UserPhoneNumber,
			},
			Address:   address(item.Address),
			CreatedAt: item.CreatedAt,
		})
	}

	return &ManyOrdersResponseDao{
		Orders: orders,
	}
}
