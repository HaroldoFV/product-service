package usecase

import (
	"github.com/HaroldoFV/product-service/internal/domain"
)

type GetProductUseCase struct {
	ProductRepository domain.ProductRepositoryInterface
}

func NewGetProductUseCase(productRepository domain.ProductRepositoryInterface) *GetProductUseCase {
	return &GetProductUseCase{
		ProductRepository: productRepository,
	}
}

func (l *GetProductUseCase) Execute(id string) (ProductOutputDTO, error) {
	product, err := l.ProductRepository.GetByID(id)
	if err != nil {
		return ProductOutputDTO{}, err
	}

	var outputProduct = ProductOutputDTO{
		ID:          product.GetID(),
		Name:        product.GetName(),
		Description: product.GetDescription(),
		Price:       product.GetPrice(),
		Status:      product.GetStatus(),
	}
	return outputProduct, nil
}
