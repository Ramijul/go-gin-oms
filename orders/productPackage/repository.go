package productPackage

import (
	"github.com/Ramijul/go-gin-oms/orders/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	TableName = "products"
)

type ProductRepository interface {
	GetMany(ids []uuid.UUID) (products []*models.Product, err error)
	GetAll() (products []*models.Product, err error)
}

type Repository struct {
	Session *gorm.DB
}

func (r *Repository) GetMany(ids []uuid.UUID) (products []*models.Product, err error) {
	runner := r.Session.Table(TableName)

	// TODO: add pagination
	result := runner.Find(&products, ids)

	if result.Error != nil {
		return nil, result.Error
	}

	return products, nil
}

func (r *Repository) GetAll() (products []*models.Product, err error) {
	runner := r.Session.Table(TableName)
	result := runner.Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}
