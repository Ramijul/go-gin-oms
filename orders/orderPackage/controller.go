package orderPackage

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/Ramijul/go-gin-oms/orders/rabbitmq"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type orderService interface {
	GetAll() (ordersDao *ManyOrdersResponseDao, err error)
	GetOne(id uuid.UUID) (ordersDao *OrderResponseDao, err error)
	Create(requestDao *CreateRequestDao) (ordersDao *OrderResponseDao, err error)
}

type rabbitMQService interface {
	SendMessage(message amqp.Publishing) error
}

type Controller struct {
	Service         orderService
	RabbitMQService rabbitMQService
}

func (c *Controller) GetAll(ctx *gin.Context) {
	orders, err := c.Service.GetAll()

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

	order, orderFetchErr := c.Service.GetOne(id)
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

	order, err := c.Service.Create(&createReqDao)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "Invalid input")
		return
	}

	go sendPaymentRequest(order, c.RabbitMQService)

	ctx.JSON(http.StatusOK, order)
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
