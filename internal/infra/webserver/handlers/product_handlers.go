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

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product dto.CreateProductInput
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	p, err := entity.NewPoduct(product.Name, product.Price)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.ProductDB.Create(p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

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
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	products, err := h.ProductDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var product entity.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	product.ID, err = entityPkg.ParseID(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = h.ProductDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = h.ProductDB.Update(&product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err := h.ProductDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = h.ProductDB.Delete(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
