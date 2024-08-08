package domain

import domain "github.com/HaroldoFV/product-service/internal/domain/entity"

type ProductRepositoryInterface interface {
	Create(product *domain.Product) error
	Update(product *domain.Product) error
	GetByID(id string) (*domain.Product, error)
	List() ([]*domain.Product, error)
}
