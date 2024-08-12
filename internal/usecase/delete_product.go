package usecase

import (
	"github.com/HaroldoFV/product-service/internal/domain"
)

type DeleteProductUseCase struct {
	ProductRepository domain.ProductRepositoryInterface
}

func NewDeleteProductUseCase(
	productRepository domain.ProductRepositoryInterface,
) *DeleteProductUseCase {
	return &DeleteProductUseCase{
		ProductRepository: productRepository,
	}
}

func (u *DeleteProductUseCase) Execute(id string) error {
	err := u.ProductRepository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
