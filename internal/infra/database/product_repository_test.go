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

func (suite *ProductRepositoryTestSuite) TestList() {
	// Create test products
	products := []struct {
		name        string
		description string
		price       float64
	}{
		{"Product A", "Description A", 10.0},
		{"Product B", "Description B", 20.0},
		{"Product C", "Description C", 30.0},
		{"Product D", "Description D", 40.0},
		{"Product E", "Description E", 50.0},
	}

	for _, p := range products {
		product, err := entity.NewProduct(p.name, p.description, p.price)
		assert.NoError(suite.T(), err)
		err = suite.Repository.Create(product)
		assert.NoError(suite.T(), err)
	}

	// Test cases
	testCases := []struct {
		name          string
		page          int
		limit         int
		sort          string
		expectedCount int
		expectedTotal int
	}{
		{"First page, default sort", 1, 3, "id", 3, 5},
		{"Second page, default sort", 2, 3, "id", 2, 5},
		{"All products, sort by price", 1, 10, "price", 5, 5},
		{"Invalid sort field", 1, 5, "invalid", 5, 5}, // Should default to "id"
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			resultProducts, totalCount, err := suite.Repository.List(tc.page, tc.limit, tc.sort)

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCount, len(resultProducts))
			assert.Equal(t, tc.expectedTotal, totalCount)

			if tc.sort == "price" && len(resultProducts) > 1 {
				for i := 1; i < len(resultProducts); i++ {
					assert.GreaterOrEqual(t, resultProducts[i].GetPrice(), resultProducts[i-1].GetPrice())
				}
			}
		})
	}
}

func (suite *ProductRepositoryTestSuite) TestUpdate() {
	initialProduct, err := entity.NewProduct("Test Product", "Test Description", 10.0)
	suite.Require().NoError(err)

	err = suite.Repository.Create(initialProduct)
	suite.Require().NoError(err)

	err = initialProduct.Update("Updated Product", "Updated Description")
	suite.Require().NoError(err)

	err = initialProduct.ChangePrice(20.0)
	suite.Require().NoError(err)

	err = suite.Repository.Update(initialProduct)
	suite.Require().NoError(err)

	updatedProduct, err := suite.Repository.GetByID(initialProduct.GetID())
	suite.Require().NoError(err)

	assert.Equal(suite.T(), "Updated Product", updatedProduct.GetName())
	assert.Equal(suite.T(), "Updated Description", updatedProduct.GetDescription())
	assert.Equal(suite.T(), 20.0, updatedProduct.GetPrice())
}

func (suite *ProductRepositoryTestSuite) TestGetByID() {
	product, err := entity.NewProduct("Test Product", "Test Description", 10.0)
	suite.Require().NoError(err)

	err = suite.Repository.Create(product)
	suite.Require().NoError(err)

	testCases := []struct {
		name          string
		id            string
		expectedError bool
	}{
		{"Existing product", product.GetID(), false},
		{"Non-existing product", "non-existent-id", true},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			retrievedProduct, err := suite.Repository.GetByID(tc.id)

			if tc.expectedError {
				assert.Error(t, err)
				assert.Nil(t, retrievedProduct)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, retrievedProduct)
				assert.Equal(t, product.GetID(), retrievedProduct.GetID())
				assert.Equal(t, product.GetName(), retrievedProduct.GetName())
				assert.Equal(t, product.GetDescription(), retrievedProduct.GetDescription())
				assert.Equal(t, product.GetPrice(), retrievedProduct.GetPrice())
				assert.Equal(t, product.GetStatus(), retrievedProduct.GetStatus())
			}
		})
	}
}

func TestProductRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(ProductRepositoryTestSuite))
}
