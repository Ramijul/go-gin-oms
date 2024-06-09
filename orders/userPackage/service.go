package userPackage

import (
	"github.com/Ramijul/go-gin-oms/orders/models"
	"github.com/google/uuid"
)

type Service struct {
	Repository UserRepository
}

func (s *Service) GetAll() (users []*models.User, err error) {
	users, err = s.Repository.GetAll()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *Service) GetOne(id uuid.UUID) (users *models.User, err error) {
	users, err = s.Repository.GetOne(id)
	if err != nil {
		return nil, err
	}

	return users, nil
}
