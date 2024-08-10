package database_test

import (
	"database/sql"
	"log"
	"testing"

	"github.com/HaroldoFV/product-service/internal/domain/entity"
	"github.com/HaroldoFV/product-service/internal/infra/database"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ProductRepositoryTestSuite struct {
	suite.Suite
	DB         *sql.DB
	Repository *database.ProductRepository
}

func (suite *ProductRepositoryTestSuite) SetupSuite() {
	connectionString := "host=localhost port=5433 user=root_test password=root_test dbname=test_product_db sslmode=disable"
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	suite.DB = db
	suite.Repository = database.NewProductRepository(db)

	// Create the products table
	_, err = suite.DB.Exec(`
		CREATE TABLE IF NOT EXISTS products (
			id VARCHAR(36) PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			description VARCHAR(500),
			price DECIMAL(10, 2) NOT NULL,
			status VARCHAR(10) NOT NULL
		)
	`)
	if err != nil {
		log.Fatal(err)
	}
}

func (suite *ProductRepositoryTestSuite) TearDownSuite() {
	_, err := suite.DB.Exec("DROP TABLE IF EXISTS products")
	if err != nil {
		log.Fatal(err)
	}
	suite.DB.Close()
}

func (suite *ProductRepositoryTestSuite) SetupTest() {
	_, err := suite.DB.Exec("DELETE FROM products")
	if err != nil {
		log.Fatal(err)
	}
}

func (suite *ProductRepositoryTestSuite) TestCreateProduct() {
	product, err := entity.NewProduct("Test Product", "Test Description", 10.0)
	assert.NoError(suite.T(), err)

	err = suite.Repository.Create(product)
	assert.NoError(suite.T(), err)

	// Verify the product was created
	var count int
	err = suite.DB.QueryRow("SELECT COUNT(*) FROM products WHERE id = $1", product.GetID()).Scan(&count)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, count)
}

func (suite *ProductRepositoryTestSuite) TestCreateProductWithInvalidData() {
	product, err := entity.NewProduct("", "Invalid Product", -5.0)
	assert.Error(suite.T(), err)

	if product != nil {
		err = suite.Repository.Create(product)
		assert.Error(suite.T(), err)
	}
}

func TestProductRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(ProductRepositoryTestSuite))
}
