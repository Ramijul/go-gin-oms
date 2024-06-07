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

	//create/update tables
	session.AutoMigrate(&models.User{}, &models.Product{}, &models.Order{}, &models.OrderDetails{})

	//default user
	user := models.User{Name: "Ramijul Islam", Email: "ramizul127@gmail.com", PhoneNumber: "3654760064"}
	result := session.Create(&user)
	if result.Error != nil {
		log.Fatal("Unable to create default user", result.Error)
	}

	//default product
	product := models.Product{Name: "Huggies Little Movers Baby Diapers", Price: 27.99, InStock: 1000}
	result = session.Create(&product)
	if result.Error != nil {
		log.Fatal("Unable to create default product", result.Error)
	}

}
