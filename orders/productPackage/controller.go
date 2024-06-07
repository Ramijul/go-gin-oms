package productPackage

import (
	"log"

	"github.com/Ramijul/go-gin-oms/orders/models"
	"github.com/gin-gonic/gin"
)

type productService interface {
	GetAll() (products []*models.Product, err error)
}

type Controller struct {
	Service productService
}

func (c *Controller) GetAll(ctx *gin.Context) {
	products, err := c.Service.GetAll()
	if err != nil {
		log.Fatal(err)
		ctx.AbortWithStatus(500)
	}

	ctx.JSON(200, gin.H{
		"products": products,
	})
}
