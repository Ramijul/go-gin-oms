package orderPackage

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/Ramijul/go-gin-oms/orders/models"
	"github.com/Ramijul/go-gin-oms/orders/rabbitmq"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type orderService interface {
	GetAll() (ordersDao *ManyOrdersResponseDao, err error)
	GetOne(id uuid.UUID) (ordersDao *OrderResponseDao, err error)
	Create(requestDao *CreateRequestDao, userData *models.User, productsRequested []*models.Product) (ordersDao *OrderResponseDao, err error)
}

type rabbitMQService interface {
	SendMessage(message amqp.Publishing) error
}

type productService interface {
	GetAll() (products []*models.Product, err error)
	GetMany(ids []uuid.UUID) (products []*models.Product, err error)
}

type userService interface {
	GetAll() (users []*models.User, err error)
	GetOne(id uuid.UUID) (users *models.User, err error)
}

type Controller struct {
	OrderService    orderService
	RabbitMQService rabbitMQService
	ProductService  productService
	UserService     userService
}

func (c *Controller) GetAll(ctx *gin.Context) {
	orders, err := c.OrderService.GetAll()

	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, orders)
}

func (c *Controller) GetOne(ctx *gin.Context) {
	idParam, provided := ctx.Params.Get("id")

	if !provided {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	id, err := uuid.Parse(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	order, orderFetchErr := c.OrderService.GetOne(id)
	if orderFetchErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	ctx.JSON(http.StatusOK, order)
}

func (c *Controller) Create(ctx *gin.Context) {
	var createReqDao CreateRequestDao

	// TODO: fidn a way to make this globally available
	if err := ctx.ShouldBind(&createReqDao); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = ErrorMsg{fe.Field(), getErrorMsg(fe)}
			}
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": out})
		}
		return
	}

	/* TODO:
	- add validation for quantity - must not exceed product.InStock
	- reduce product.InStock after order has been placed
	*/

	// get the user
	userData, err := getUserData(createReqDao.UserID, c.UserService)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "Invalid user")
		return
	}

	// get the products from db
	productsRequested, productsRequestedErr := getProductsRequested(createReqDao.Products, c.ProductService)
	if productsRequestedErr != nil {
		ctx.JSON(http.StatusBadRequest, "Invalid product(s)")
		return
	}

	order, err := c.OrderService.Create(&createReqDao, userData, productsRequested)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "Invalid input")
		return
	}

	go sendPaymentRequest(order, c.RabbitMQService)

	ctx.JSON(http.StatusOK, order)
}

func getProductsRequested(items []*RequestedItem, ProductService productService) ([]*models.Product, error) {
	//extract product ids from the request
	var productIds []uuid.UUID
	for _, requestedItem := range items {
		pid, err := uuid.Parse(requestedItem.ID)
		if err != nil {
			return nil, err
		}

		productIds = append(productIds, pid)
	}

	// get the products from db
	return ProductService.GetMany(productIds)
}

func getUserData(id string, UserService userService) (*models.User, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	return UserService.GetOne(uid)
}

func sendPaymentRequest(order *OrderResponseDao, rabbitMQService rabbitMQService) {
	// message body
	orderCreateEvent, err := json.Marshal(&rabbitmq.OrderCreateEvent{
		OrderID:    order.ID,
		TotalPrice: order.TotalPrice,
	})

	if err != nil {
		log.Print("Failed to Marshall OrderCreateEvent", err)
		return
	}

	// send message
	err = rabbitMQService.SendMessage(amqp.Publishing{
		ContentType: "application/json",
		Body:        orderCreateEvent,
	})

	if err != nil {
		log.Print("Failed to Send Message to Payments Service", err)
		return
	}

	log.Print("Payment request sent ", orderCreateEvent)

}

/*
ref: https://blog.logrocket.com/gin-binding-in-go-a-tutorial-with-examples/
*/
type ErrorMsg struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

/*
ref: https://blog.logrocket.com/gin-binding-in-go-a-tutorial-with-examples/
*/
func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "min":
		return "Should be greater than " + fe.Param()
	case "uuid4":
		return "Incorrect format " + fe.Param()
	case "gt":
		return "Must be greater than " + fe.Param()
	}
	return "Unknown error"
}
