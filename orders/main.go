package main

/*
TODO: add better logging
*/

import (
	"log"

	"github.com/Ramijul/go-gin-oms/orders/orderPackage"
	"github.com/Ramijul/go-gin-oms/orders/productPackage"
	"github.com/Ramijul/go-gin-oms/orders/userPackage"
	db "github.com/Ramijul/go-gin-oms/orders/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	session, err := db.Connect()
	if err != nil {
		panic(err)
	}
	if session != nil {
		log.Println("Database connected")
	}

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
		Service: orderService,
	}

	r := gin.Default()
	r.GET("/products", productController.GetAll)
	r.GET("/users", userController.GetAll)
	r.GET("/orders", orderController.GetAll)
	r.GET("/order/:id", orderController.GetOne)
	r.POST("/order", orderController.Create)
	r.Run()
}
