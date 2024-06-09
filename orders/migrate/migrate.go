package main

import (
	"log"
	"os"
	"strings"

	"github.com/Ramijul/go-gin-oms/orders/models"
	db "github.com/Ramijul/go-gin-oms/orders/utils"
	"github.com/joho/godotenv"
)

func init() {

	// fails when env is passed from docker-compose
	err := godotenv.Load()
	if err != nil {
		// test for an env
		if len(os.Getenv("RABBITMQ_CONN_STRING")) == 0 {
			panic("Error loading .env file")
		}
	}
}

func main() {
	session, err := db.Connect()
	if err != nil {
		panic(err)
	}

	if session != nil {
		log.Println("Database connected, initiating migration")
	}

	//create/update tables
	session.AutoMigrate(&models.User{}, &models.Product{}, &models.Order{}, &models.OrderItem{})

	//default user
	user := models.User{Name: "Ramijul Islam", Email: "ramizul127@gmail.com", PhoneNumber: "3654760064"}
	result := session.Create(&user)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key value violates unique constraint") {
			return
		}
		log.Fatal("Unable to create default user", result.Error)
	}

	//default product
	product := models.Product{Name: "Huggies Little Movers Baby Diapers", Price: 27.99, InStock: 1000}
	result = session.Create(&product)
	if result.Error != nil {
		log.Fatal("Unable to create default product", result.Error)
	}

}
