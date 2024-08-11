package web

import (
	"encoding/json"
	"fmt"
	"github.com/HaroldoFV/product-service/internal/domain"
	usecase "github.com/HaroldoFV/product-service/internal/usecase"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
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

// List Products godoc
// @Summary List Products
// @Description List Products
// @Tags products
// @Accept json
// @Produce json
// @Param page query int false "page number" default(1)
// @Param limit query int false "limit" default(10)
// @Param sort query string false "sort field" default("id")
// @Success 200 {object} PaginatedProductResponse
// @Failure 404 {object} Error
// @Failure 500 {object} Error
// @Router /products [get]
func (h *WebProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}

	sort := r.URL.Query().Get("sort")
	if sort == "" {
		sort = "id"
	}

	listProductsUseCase := usecase.NewListProductsUseCase(h.ProductRepository)
	output, totalCount, err := listProductsUseCase.Execute(page, limit, sort)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := PaginatedProductResponse{
		Products:   output,
		TotalCount: totalCount,
		Page:       page,
		Limit:      limit,
		TotalPages: (totalCount + limit - 1) / limit,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		fmt.Println("Error encoding response:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

// Update Product godoc
// @Summary Update Product
// @Description Update Product
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID" Format(uuid)
// @Param request body usecase.ProductUpdateInputDTO true "product Request"
// @Success 200 {object} usecase.ProductOutputDTO
// @Failure 404 {object} Error
// @Failure 500 {object} Error
// @Router /products/{id} [put]
func (h *WebProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(Error{Message: "missing product ID"})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	var dto usecase.ProductUpdateInputDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(Error{Message: err.Error()})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	dto.ID = id

	updateProductUseCase := usecase.NewUpdateProductUseCase(h.ProductRepository)
	output, err := updateProductUseCase.Execute(dto)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "product not found" {
			status = http.StatusNotFound
		}
		w.WriteHeader(status)
		err := json.NewEncoder(w).Encode(Error{Message: err.Error()})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)
}

// GetProduct godoc
// @Summary Get Product
// @Description Get Product
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID" Format(uuid)
// @Success 200 {object} usecase.ProductOutputDTO
// @Failure 404 {object} Error
// @Failure 500 {object} Error
// @Router /products/{id} [get]
func (h *WebProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(Error{Message: "missing product ID"})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	getProductUseCase := usecase.NewGetProductUseCase(h.ProductRepository)
	output, err := getProductUseCase.Execute(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == fmt.Sprintf("product with id %s not found", id) {
			status = http.StatusNotFound
		}
		w.WriteHeader(status)
		err := json.NewEncoder(w).Encode(Error{Message: err.Error()})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	response := usecase.ProductOutputDTO{
		ID:          output.ID,
		Name:        output.Name,
		Description: output.Description,
		Price:       output.Price,
		Status:      output.Status,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		fmt.Println("Error encoding response:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

type PaginatedProductResponse struct {
	Products   []usecase.ProductOutputDTO `json:"products"`
	TotalCount int                        `json:"total_count"`
	Page       int                        `json:"page"`
	Limit      int                        `json:"limit"`
	TotalPages int                        `json:"total_pages"`
}

// Error represents an error response
type Error struct {
	Message string `json:"message"`
}
