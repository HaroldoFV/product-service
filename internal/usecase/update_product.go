package usecase

import (
	"github.com/HaroldoFV/product-service/internal/domain"
)

type UpdateProductUseCase struct {
	ProductRepository domain.ProductRepositoryInterface
}

func NewUpdateProductUseCase(
	productRepository domain.ProductRepositoryInterface,
) *UpdateProductUseCase {
	return &UpdateProductUseCase{
		ProductRepository: productRepository,
	}
}

func (u *UpdateProductUseCase) Execute(input ProductUpdateInputDTO) (ProductOutputDTO, error) {
	product, err := u.ProductRepository.GetByID(input.ID)
	if err != nil {
		return ProductOutputDTO{}, err
	}

	err = product.Update(input.Name, input.Description)
	if err != nil {
		return ProductOutputDTO{}, err
	}

	err = product.ChangePrice(input.Price)
	if err != nil {
		return ProductOutputDTO{}, err
	}

	err = u.ProductRepository.Update(product)
	if err != nil {
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
