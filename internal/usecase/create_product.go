package usecase

import (
	"github.com/HaroldoFV/product-service/internal/domain"
	"github.com/HaroldoFV/product-service/internal/domain/entity"
)

type CreateProductUseCase struct {
	ProductRepository domain.ProductRepositoryInterface
}

func NewCreateProductUseCase(
	productRepository domain.ProductRepositoryInterface,
) *CreateProductUseCase {
	return &CreateProductUseCase{
		ProductRepository: productRepository,
	}
}

func (c *CreateProductUseCase) Execute(input ProductInputDTO) (ProductOutputDTO, error) {
	product, _ := entity.NewProduct(
		input.Name,
		input.Description,
		input.Price,
	)

	if err := c.ProductRepository.Create(product); err != nil {
		return ProductOutputDTO{}, err
	}
	dto := ProductOutputDTO{
		ID:          product.GetID(),
		Name:        product.GetName(),
		Description: product.GetDescription(),
		Price:       product.GetPrice(),
		Status:      product.GetStatus(),
	}

	return dto, nil
}
