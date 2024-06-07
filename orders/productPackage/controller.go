package productPackage

import (
	"net/http"

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
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"products": products,
	})
}
