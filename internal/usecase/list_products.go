package usecase

import (
	"github.com/HaroldoFV/product-service/internal/domain"
)

type ListProductsUseCase struct {
	ProductRepository domain.ProductRepositoryInterface
}

func NewListProductsUseCase(productRepository domain.ProductRepositoryInterface) *ListProductsUseCase {
	return &ListProductsUseCase{
		ProductRepository: productRepository,
	}
}

func (l *ListProductsUseCase) Execute(page, limit int, sort string) ([]ProductOutputDTO, int, error) {
	products, totalCount, err := l.ProductRepository.List(page, limit, sort)
	if err != nil {
		return nil, 0, err
	}

	var outputProducts []ProductOutputDTO
	for _, product := range products {
		outputProducts = append(outputProducts, ProductOutputDTO{
			ID:          product.GetID(),
			Name:        product.GetName(),
			Description: product.GetDescription(),
			Price:       product.GetPrice(),
			Status:      product.GetStatus(),
		})
	}
	return outputProducts, totalCount, nil
}
