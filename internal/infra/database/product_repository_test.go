package database

import (
	"database/sql"
	"fmt"
	domain "github.com/HaroldoFV/product-service/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "test_product_db"
)

func setupTestDB(t *testing.T) *sql.DB {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", host, port, user, password)
	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbname))
	if err != nil {
		t.Logf("Database creation error (may already exist): %v", err)
	}

	db.Close()

	connStr = fmt.Sprintf("%s dbname=%s", connStr, dbname)
	db, err = sql.Open("postgres", connStr)
	require.NoError(t, err)

	_, err = db.Exec(`
				CREATE TABLE IF NOT EXISTS products (
					id UUID PRIMARY KEY ,
					name VARCHAR(100) NOT NULL,
					description VARCHAR(500),
					price DECIMAL(10,2) NOT NULL,
					status VARCHAR(20) NOT NULL
				)
	`)
	require.NoError(t, err)

	return db
}

func cleanupTestDB(t *testing.T, db *sql.DB) {
	_, err := db.Exec("DElETE FROM products")
	require.NoError(t, err)
	db.Close()

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", host, port, user, password)
	db, err = sql.Open("postgres", connStr)
	require.NoError(t, err)
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbname))
	require.NoError(t, err)
}

func TestProductRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := NewProductRepository(db)

	t.Run("should create a product successfully", func(t *testing.T) {
		product, err := domain.NewProduct("Test Product", "A test product", 10.00)
		require.NoError(t, err)

		err = repo.Create(product)
		assert.NoError(t, err)

		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM products WHERE id = $1", product.GetID()).Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)
	})

	t.Run("should not create invalid product", func(t *testing.T) {
		_, err := domain.NewProduct("", "", -1)
		assert.Error(t, err)
	})

	t.Run("should return error when database operation fails", func(t *testing.T) {
		validProduct, err := domain.NewProduct("Valid Product", "Description", 10.00)
		require.NoError(t, err)

		db.Close()

		err = repo.Create(validProduct)
		assert.Error(t, err)

		db = setupTestDB(t)
		repo = NewProductRepository(db)
	})
}
