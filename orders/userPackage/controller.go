package userPackage

import (
	"log"

	"github.com/Ramijul/go-gin-oms/orders/models"
	"github.com/gin-gonic/gin"
)

type userService interface {
	GetAll() (users []*models.User, err error)
}

type Controller struct {
	Service userService
}

func (c *Controller) GetAll(ctx *gin.Context) {
	users, err := c.Service.GetAll()

	if err != nil {
		log.Fatal(err)
		ctx.AbortWithStatus(500)
	}

	ctx.JSON(200, gin.H{
		"users": users,
	})
}
