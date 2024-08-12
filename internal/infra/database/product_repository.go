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

func (r *ProductRepository) List(page, limit int, sort string) ([]*domain.Product, int, error) {
	offset := (page - 1) * limit

	// Count total products
	var totalCount int
	err := r.Db.QueryRow("SELECT COUNT(*) FROM products").Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// Validate and sanitize sort field
	validSortFields := map[string]bool{"id": true, "name": true, "price": true}
	if !validSortFields[sort] {
		sort = "id"
	}

	query := fmt.Sprintf("SELECT id, name, description, price, status FROM products ORDER BY %s LIMIT $1 OFFSET $2", sort)

	// Print the query with actual values
	fmt.Printf("Executing main query: %s [LIMIT %d OFFSET %d]\n", query, limit, offset)

	rows, err := r.Db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		var id, name, description, status string
		var price float64

		err := rows.Scan(&id, &name, &description, &price, &status)
		if err != nil {
			return nil, 0, err
		}
		product, err := domain.NewProduct(name, description, price)
		if err != nil {
			return nil, 0, err
		}
		product.SetID(id)

		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}
	return products, totalCount, nil
}

func (r *ProductRepository) Update(product *domain.Product) error {
	stmt, err := r.Db.Prepare("UPDATE products SET name = $1, description = $2, price = $3 WHERE id = $4")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(product.GetName(), product.GetDescription(), product.GetPrice(),
		product.GetID())
	if err != nil {
		return err
	}
	return nil
}

func (r *ProductRepository) GetByID(id string) (*domain.Product, error) {
	query := "SELECT id, name, description, price, status FROM products WHERE id = $1"

	row := r.Db.QueryRow(query, id)

	var product *domain.Product
	var idStr, name, description, status string
	var price float64

	err := row.Scan(&idStr, &name, &description, &price, &status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product with id %s not found", id)
		}
		return nil, err
	}

	product, err = domain.NewProduct(name, description, price)
	if err != nil {
		return nil, err
	}

	product.SetID(idStr)

	if status == domain.ENABLED {
		err = product.Enable()
	} else {
		err = product.Disable()
	}
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *ProductRepository) Delete(id string) error {
	query := "DELETE FROM products WHERE id = $1"
	result, err := r.Db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product with id %s not found", id)
	}

	return nil
}
