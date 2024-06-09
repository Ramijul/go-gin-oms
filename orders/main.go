package main

/*
TODO: add better logging
*/

import (
	"github.com/Ramijul/go-gin-oms/orders/orderPackage"
	"github.com/Ramijul/go-gin-oms/orders/productPackage"
	"github.com/Ramijul/go-gin-oms/orders/rabbitmq"
	"github.com/Ramijul/go-gin-oms/orders/userPackage"
	db "github.com/Ramijul/go-gin-oms/orders/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	rabbitmq.IntiallizeVariables()
}

func main() {
	session, err := db.Connect()
	if err != nil {
		panic(err)
	}

	// initialize rabbitmq service
	conn, ch, q := rabbitmq.InitializeRabbitMQ(rabbitmq.REQUEST_QUEUE)
	rabbitMQService := &rabbitmq.RabbitMQService{
		Conn: conn,
		Ch:   ch,
		Q:    q,
	}
	defer rabbitMQService.CloseConnection()

	// TODO: move dependency injections and routing elsewhere

	// initialize product
	productRepo := &productPackage.Repository{
		Session: session,
	}
	productService := &productPackage.Service{
		Repository: productRepo,
	}
	productController := productPackage.Controller{
		Service: productService,
	}

	// initialize user
	userRepo := &userPackage.Repository{
		Session: session,
	}
	userService := &userPackage.Service{
		Repository: userRepo,
	}
	userController := userPackage.Controller{
		Service: userService,
	}

	// initialize order
	orderRepo := &orderPackage.Repository{
		Session: session,
	}
	orderService := &orderPackage.Service{
		OrderRepository:   orderRepo,
		ProductRepository: productRepo,
		UserRepository:    userRepo,
	}
	orderController := orderPackage.Controller{
		Service:         orderService,
		RabbitMQService: rabbitMQService,
	}

	// consumer process
	go orderPackage.ConsumePaymentConfirmation(*orderService)

	// app on main thread
	r := gin.Default()
	r.GET("/products", productController.GetAll)
	r.GET("/users", userController.GetAll)
	r.GET("/orders", orderController.GetAll)
	r.GET("/order/:id", orderController.GetOne)
	r.POST("/order", orderController.Create)
	r.Run()
}
