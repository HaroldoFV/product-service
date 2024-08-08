package entity_test

import (
	entity "github.com/HaroldoFV/product-service/internal/domain/entity"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewProduct(t *testing.T) {
	t.Run("Valid Product", func(t *testing.T) {
		product, err := entity.NewProduct("Product 1", "description", 99.99)
		require.Nil(t, err)
		require.NotNil(t, product)
		require.Equal(t, "Product 1", product.GetName())
		require.Equal(t, "description", product.GetDescription())
		require.Equal(t, entity.DISABLED, product.GetStatus())
	})

	t.Run("Invalid Name", func(t *testing.T) {
		_, err := entity.NewProduct("", "description", 99.99)
		require.EqualError(t, err, "name cannot be empty")
	})

	t.Run("Name Too Long", func(t *testing.T) {
		longName := string(make([]byte, 101))
		_, err := entity.NewProduct(longName, "description", 99.99)
		require.EqualError(t, err, "name cannot be longer than 100 characters")
	})

	t.Run("Description Too Long", func(t *testing.T) {
		longDesc := string(make([]byte, 501))
		_, err := entity.NewProduct("Product 1", longDesc, 99.99)
		require.EqualError(t, err, "description cannot be longer than 500 characters")
	})

	t.Run("Invalid Price", func(t *testing.T) {
		_, err := entity.NewProduct("Product 1", "description", -1)
		require.EqualError(t, err, "price must be greater or equal zero")
	})
}

func TestProduct_Enable(t *testing.T) {
	t.Run("Enable with Valid Price", func(t *testing.T) {
		product, _ := entity.NewProduct("Product 1", "Product 1 description", 99.90)
		err := product.Enable()
		require.Nil(t, err)
		require.Equal(t, entity.ENABLED, product.GetStatus())
	})
}

func TestProduct_Disable(t *testing.T) {
	t.Run("Disable Product", func(t *testing.T) {
		product, _ := entity.NewProduct("Product 1", "Product 1 description", 99.90)
		err := product.Enable()
		require.Nil(t, err)
		err = product.Disable()
		require.Nil(t, err)
		require.Equal(t, entity.DISABLED, product.GetStatus())
	})
}

func TestProduct_ChancePrice(t *testing.T) {
	t.Run("Change to Valid Price", func(t *testing.T) {
		product, _ := entity.NewProduct("Product 1", "Product 1 description", 99.90)
		err := product.ChancePrice(150.00)
		require.Nil(t, err)
		require.Equal(t, 150.00, product.GetPrice())
	})

	t.Run("Change to Invalid Price", func(t *testing.T) {
		product, _ := entity.NewProduct("Product 1", "Product 1 description", 99.90)
		err := product.ChancePrice(-10.00)
		require.EqualError(t, err, "price must be greater or equal zero")
	})
}

func TestProduct_GetMethods(t *testing.T) {
	product, _ := entity.NewProduct("Product 1", "Product 1 description", 99.90)

	t.Run("GetID", func(t *testing.T) {
		require.NotEmpty(t, product.GetID())
	})

	t.Run("GetName", func(t *testing.T) {
		require.Equal(t, "Product 1", product.GetName())
	})

	t.Run("GetDescription", func(t *testing.T) {
		require.Equal(t, "Product 1 description", product.GetDescription())
	})

	t.Run("GetStatus", func(t *testing.T) {
		require.Equal(t, entity.DISABLED, product.GetStatus())
	})

	t.Run("GetPrice", func(t *testing.T) {
		require.Equal(t, 99.90, product.GetPrice())
	})
}

func TestProduct_Update(t *testing.T) {
	t.Run("Update with valid data", func(t *testing.T) {
		product, _ := entity.NewProduct("Original Product", "Original Description", 100.00)
		err := product.Update("Updated Product", "New Description")
		require.Nil(t, err)
		require.Equal(t, "Updated Product", product.GetName())
		require.Equal(t, "New Description", product.GetDescription())
		require.Equal(t, 100.00, product.GetPrice())
	})

	t.Run("Update with empty name", func(t *testing.T) {
		product, _ := entity.NewProduct("Original Product", "Original Description", 100.00)
		err := product.Update("", "New Description")
		require.EqualError(t, err, "name cannot be empty")
	})

	t.Run("Update with too long name", func(t *testing.T) {
		product, _ := entity.NewProduct("Original Product", "Original Description", 100.00)
		longName := string(make([]byte, 101))
		err := product.Update(longName, "New Description")
		require.EqualError(t, err, "name cannot be longer than 100 characters")
	})

	t.Run("Update with too long description", func(t *testing.T) {
		product, _ := entity.NewProduct("Original Product", "Original Description", 100.00)
		longDesc := string(make([]byte, 501))
		err := product.Update("Updated Product", longDesc)
		require.EqualError(t, err, "description cannot be longer than 500 characters")
	})

	t.Run("Update maintaining status", func(t *testing.T) {
		product, _ := entity.NewProduct("Original Product", "Original Description", 100.00)
		err := product.Enable()
		require.Nil(t, err)
		err = product.Update("Updated Product", "New Description")
		require.Nil(t, err)
		require.Equal(t, entity.ENABLED, product.GetStatus())
	})
}
