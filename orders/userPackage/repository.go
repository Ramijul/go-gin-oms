package userPackage

import (
	"github.com/Ramijul/go-gin-oms/orders/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	TableName = "users"
)

type UserRepository interface {
	GetAll() (users []*models.User, err error)
	GetOne(id uuid.UUID) (user *models.User, err error)
}

type Repository struct {
	Session *gorm.DB
}

func (r *Repository) GetAll() (users []*models.User, err error) {
	runner := r.Session.Table(TableName)
	result := runner.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (r *Repository) GetOne(id uuid.UUID) (user *models.User, err error) {
	runner := r.Session.Table(TableName)
	result := runner.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}
