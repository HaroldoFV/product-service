package database

import (
	"database/sql"
	"fmt"
	domain "github.com/HaroldoFV/product-service/internal/domain/entity"
)

type ProductRepository struct {
	Db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	fmt.Printf("Creating ProductRepository with db: %v\n", db)
	return &ProductRepository{Db: db}
}

func (r *ProductRepository) Create(product *domain.Product) error {
	stmt, err := r.Db.Prepare("INSERT INTO products (id, name, description, price, status) VALUES ($1, $2, $3, $4, $5)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(product.GetID(), product.GetName(), product.GetDescription(), product.GetPrice(), product.GetStatus())
	if err != nil {
		return err
	}
	return nil
}
