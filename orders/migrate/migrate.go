package main

import (
	"log"

	"github.com/Ramijul/go-gin-oms/orders/models"
	db "github.com/Ramijul/go-gin-oms/orders/utils"
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
		log.Println("Database connected, inititing migration")
	}

	session.AutoMigrate(&models.User{}, &models.Product{}, &models.Order{}, &models.OrderDetails{})
}
