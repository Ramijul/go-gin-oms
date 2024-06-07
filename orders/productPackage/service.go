package productPackage

import (
	"github.com/Ramijul/go-gin-oms/orders/models"
)

type Service struct {
	Repository ProductRepository
}

func (s *Service) GetAll() (products []*models.Product, err error) {
	products, err = s.Repository.GetAll()
	if err != nil {
		return nil, err
	}

	return products, nil
}
