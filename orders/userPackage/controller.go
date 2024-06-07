package userPackage

import (
	"net/http"

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
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}
