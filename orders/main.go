package main

import (
	"log"

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

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run()
}
