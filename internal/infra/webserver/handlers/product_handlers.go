package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ClaudionorJunior/go-expert-api/internal/dto"
	"github.com/ClaudionorJunior/go-expert-api/internal/entity"
	"github.com/ClaudionorJunior/go-expert-api/internal/infra/database"
	entityPkg "github.com/ClaudionorJunior/go-expert-api/pkg/entity"
	"github.com/go-chi/chi/v5"
)

type ProductHandler struct {
	ProductDB database.ProductInterface
}

func NewProductHandler(db database.ProductInterface) *ProductHandler {
	return &ProductHandler{
		ProductDB: db,
	}
}

// Create Product godoc
// @Summary Create product
// @Description Create products
// @Tags products
// @Accept  json
// @Produce  json
// @Param request body dto.CreateProductInput true "product request"
// @Success 201
// @Failure 400 {object} dto.ErrorMessageResponse
// @Failure 500 {object} dto.ErrorMessageResponse
// @Router /products [post]
// @Security ApiKeyAuth
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product dto.CreateProductInput
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorMessageResponse{
			Message: "Invalid product data",
		})
		return
	}

	p, err := entity.NewPoduct(product.Name, product.Price)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorMessageResponse{
			Message: "Errored product data",
		})
		return
	}

	err = h.ProductDB.Create(p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorMessageResponse{
			Message: "Error to create product",
		})
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// List Product godoc
// @Summary Get all product
// @Description Get all products
// @Tags products
// @Accept  json
// @Produce  json
// @Param page query string false "page number"
// @Param limit query string false "limit number"
// @Success 200 {array} entity.Product
// @Failure 500 {object} dto.ErrorMessageResponse
// @Router /products [get]
// @Security ApiKeyAuth
func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 0
	}
	pageLimit, err := strconv.Atoi(limit)
	if err != nil {
		pageLimit = 0
	}

	sort := r.URL.Query().Get("sort")

	products, err := h.ProductDB.FindAll(pageInt, pageLimit, sort)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorMessageResponse{
			Message: "Error to get products",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

// Get Product godoc
// @Summary Get a product
// @Description Get a product
// @Tags products
// @Accept  json
// @Produce  json
// @Param id path string true "product ID" Format(uuid)
// @Success 200 {object} entity.Product
// @Failure 400 {object} dto.ErrorMessageResponse
// @Failure 404 {object} dto.ErrorMessageResponse
// @Router /products/{id} [get]
// @Security ApiKeyAuth
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorMessageResponse{
			Message: "Path id is required",
		})
		return
	}

	products, err := h.ProductDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(dto.ErrorMessageResponse{
			Message: "Product not found",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(products)
}

// Update Product godoc
// @Summary Update a product
// @Description Update a products
// @Tags products
// @Accept  json
// @Produce  json
// @Param id path string true "product ID" Format(uuid)
// @Param request body dto.CreateProductInput true "product request"
// @Success 200
// @Failure 400 {object} dto.ErrorMessageResponse
// @Failure 404 {object} dto.ErrorMessageResponse
// @Failure 500 {object} dto.ErrorMessageResponse
// @Router /products/{id} [put]
// @Security ApiKeyAuth
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorMessageResponse{
			Message: "Path id is required",
		})
		return
	}

	var product entity.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorMessageResponse{
			Message: "Invalid product data",
		})
		return
	}

	product.ID, err = entityPkg.ParseID(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorMessageResponse{
			Message: "Invalid product id",
		})
		return
	}

	_, err = h.ProductDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(dto.ErrorMessageResponse{
			Message: "Product not found",
		})
		return
	}

	err = h.ProductDB.Update(&product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorMessageResponse{
			Message: "Error to update product",
		})
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Delete Product godoc
// @Summary Delete a product
// @Description Delete a product
// @Tags products
// @Accept  json
// @Produce  json
// @Param id path string true "product ID" Format(uuid)
// @Success 201
// @Failure 400 {object} dto.ErrorMessageResponse
// @Failure 404 {object} dto.ErrorMessageResponse
// @Failure 500 {object} dto.ErrorMessageResponse
// @Router /products/{id} [delete]
// @Security ApiKeyAuth
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorMessageResponse{
			Message: "Path id is required",
		})
		return
	}

	_, err := h.ProductDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(dto.ErrorMessageResponse{
			Message: "Product not found",
		})
		return
	}

	err = h.ProductDB.Delete(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorMessageResponse{
			Message: "Error to delete product",
		})
		return
	}
	w.WriteHeader(http.StatusOK)
}
