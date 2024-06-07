package main

import (
	"log"

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

	productRepo := &productPackage.Repository{
		Session: session,
	}

	productService := &productPackage.Service{
		Repository: productRepo,
	}

	productController := productPackage.Controller{
		Service: productService,
	}

	userRepo := &userPackage.Repository{
		Session: session,
	}

	userService := &userPackage.Service{
		Repository: userRepo,
	}

	userController := userPackage.Controller{
		Service: userService,
	}

	r := gin.Default()
	r.GET("/products", productController.GetAll)
	r.GET("/users", userController.GetAll)
	r.Run()
}
