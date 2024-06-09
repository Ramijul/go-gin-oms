package orderPackage

import (
	"errors"
	"math"

	"github.com/Ramijul/go-gin-oms/orders/models"
	product "github.com/Ramijul/go-gin-oms/orders/productPackage"
	"github.com/Ramijul/go-gin-oms/orders/rabbitmq"
	user "github.com/Ramijul/go-gin-oms/orders/userPackage"
	"github.com/Ramijul/go-gin-oms/orders/utils"
	"github.com/google/uuid"
)

type Service struct {
	OrderRepository   OrderRepository
	ProductRepository product.ProductRepository
	UserRepository    user.UserRepository
}

func (s *Service) GetAll() (ordersDao *ManyOrdersResponseDao, err error) {
	orders, err := s.OrderRepository.GetAll()
	if err != nil {
		return nil, err
	}

	return ToManyOrdersResponseDao(orders), nil
}

func (s *Service) GetOne(id uuid.UUID) (ordersDao *OrderResponseDao, err error) {
	orderWithDetails, err := s.OrderRepository.GetOne(id)
	if err != nil {
		return nil, err
	}

	return ToOrderResponseDao(orderWithDetails), nil
}

func (s *Service) HandlePaymentConfirmation(paymentConfirmation rabbitmq.PaymentProcessEvent) error {
	if paymentConfirmation.PaymentStatus != string(utils.PAYMENT_STATUS_FAILED) && paymentConfirmation.PaymentStatus != string(utils.PAYMENT_STATUS_PAID) {
		return errors.New("unknown payment status detected")
	}
	orderId, err := uuid.Parse(paymentConfirmation.OrderID)
	if err != nil {
		return err
	}

	return s.OrderRepository.UpdateStatus(
		orderId,
		utils.PAYMENT_STATUS(paymentConfirmation.PaymentStatus),
		getOrderStatusBasedOnPayment(utils.PAYMENT_STATUS(paymentConfirmation.PaymentStatus)),
	)
}

/*
TODO: add transaction
*/
func (s *Service) Create(requestDao *CreateRequestDao) (ordersDao *OrderResponseDao, err error) {
	// get the user
	uid, uerr := uuid.Parse(requestDao.UserID)
	if uerr != nil {
		return nil, uerr
	}

	userData, userErr := s.UserRepository.GetOne(uid)
	if userErr != nil {
		return nil, userErr
	}

	//extract product ids from the request
	var productIds []uuid.UUID
	for _, requestedItem := range requestDao.Products {
		pid, perr := uuid.Parse(requestedItem.ID)
		if perr != nil {
			return nil, perr
		}

		productIds = append(productIds, pid)
	}

	// get the products from db
	productsRequested, productsRequestedErr := s.ProductRepository.GetMany(productIds)
	if productsRequestedErr != nil {
		return nil, productsRequestedErr
	}

	// generate a map [productID]:models.product
	productMap := getProductMap(productsRequested)

	// CREATE ORDER
	totalPrice := getTotalOrderPrice(productMap, requestDao.Products)
	order := &models.Order{
		TotalPrice:      totalPrice, //inject calculated total price
		OrderStatus:     string(utils.ORDER_STATUS_PAYMENT_PENDING),
		PaymentStatus:   string(utils.PAYMENT_STATUS_PENDING),
		UserID:          userData.ID,
		UserName:        userData.Name,
		UserEmail:       userData.Email,
		UserPhoneNumber: userData.PhoneNumber,
		Address:         models.Address(requestDao.Address),
	}

	orderId, orderCreateErr := s.OrderRepository.CreateOrder(order)
	if orderCreateErr != nil {
		return nil, orderCreateErr
	}

	// CREATE ORDER ITEMS
	var orderItems []*models.OrderItem
	for _, item := range requestDao.Products {
		productRequested := productMap[item.ID]

		orderItems = append(orderItems, &models.OrderItem{
			OrderID:          orderId, // inject the created order id
			ProductID:        productRequested.ID,
			ProductName:      productRequested.Name,
			ProductUnitPrice: productRequested.Price,
			Units:            item.Units,
		})
	}

	success, createErr := s.OrderRepository.CreateOrderItems(orderItems)
	if !success {
		return nil, createErr
	}

	return s.GetOne(orderId)
}

func getOrderStatusBasedOnPayment(paymentStatus utils.PAYMENT_STATUS) utils.ORDER_STATUS {
	switch paymentStatus {
	case utils.PAYMENT_STATUS_PAID:
		return utils.ORDER_STATUS_PROCESSING
	default:
		return utils.ORDER_STATUS_PAYMENT_PENDING
	}
}

func getProductMap(products []*models.Product) map[string]models.Product {
	productMap := make(map[string]models.Product)
	for _, elem := range products {
		productMap[elem.ID.String()] = *elem
	}
	return productMap
}

/*
Calculate total order price = price * quantity
*/
func getTotalOrderPrice(productMap map[string]models.Product, requestedItems []*RequestedItem) float64 {
	total := 0.0
	for _, item := range requestedItems {
		total += float64(item.Units) * productMap[item.ID].Price
	}

	// round to 2 decimal places
	return roundFloat(total, 2)
}

/*
Round a floating number with given precision
*/
func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
