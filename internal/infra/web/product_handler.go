package web

import (
	"encoding/json"
	"fmt"
	"github.com/HaroldoFV/product-service/internal/domain"
	"github.com/HaroldoFV/product-service/internal/usecase"
	"net/http"
)

type WebProductHandler struct {
	CreateProductUseCase *usecase.CreateProductUseCase
	ProductRepository    domain.ProductRepositoryInterface
}

func NewWebProductHandler(
	createProductUseCase *usecase.CreateProductUseCase,
	productRepository domain.ProductRepositoryInterface,
) *WebProductHandler {
	return &WebProductHandler{
		CreateProductUseCase: createProductUseCase,
		ProductRepository:    productRepository,
	}
}

// Create godoc
// @Summary Create a new product
// @Description Create a new product with the input payload
// @Tags products
// @Accept  json
// @Produce  json
// @Param product body usecase.ProductInputDTO true "Create product"
// @Success 201 {object} usecase.ProductOutputDTO
// @Router /products [post]
func (h *WebProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request to /products")

	var dto usecase.ProductInputDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		fmt.Println("Error decoding request body:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Received product: %+v\n", dto)

	output, err := h.CreateProductUseCase.Execute(dto)
	if err != nil {
		fmt.Println("Error executing create product use case:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		fmt.Println("Error encoding response:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Product created successfully")
}
