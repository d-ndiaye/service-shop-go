package api

import (
	productInternal "bitbucket.easy.de/users/n.gauche/service-shop-go/internal/product"
	"bitbucket.easy.de/users/n.gauche/service-shop-go/pkg/dto"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
)

type ProductHandler struct {
	S productInternal.Service
}

func (ph *ProductHandler) Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/", ph.GetAll)
	router.Get("/{productId}", ph.Get)
	router.Post("/", ph.Post)
	router.Delete("/{productId}", ph.Delete)
	return router
}

func (ph *ProductHandler) GetAll(response http.ResponseWriter, request *http.Request) {
	products, err := ph.S.GetAll()
	if err != nil {
		http.Error(response, "could not retrieve the products", http.StatusInternalServerError)
		return
	}
	dtos := make([]dto.Dto, 0, len(products))
	for _, p := range products {
		dto, _ := p.ToDto()
		dtos = append(dtos, dto)
	}
	render.JSON(response, request, dtos)
}

func (ph *ProductHandler) Get(response http.ResponseWriter, request *http.Request) {
	productId := chi.URLParam(request, "productId")
	p, err := ph.S.Get(productId)
	if err != nil {
		http.Error(response, "could not retrieve product", http.StatusInternalServerError)
		return
	}
	dto, _ := p.ToDto()
	render.JSON(response, request, dto)
}

func (ph *ProductHandler) decodeBody(request *http.Request, target interface{}) (err error) {
	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(target)
	return
}

func (ph *ProductHandler) Post(response http.ResponseWriter, request *http.Request) {
	dto := dto.Dto{}
	if err := ph.decodeBody(request, &dto); err != nil {
		http.Error(response, "could not decode body", http.StatusBadRequest)
		return
	}
	p := &productInternal.Product{}
	err := p.FromDto(dto)
	if err != nil {
		http.Error(response, "could not convert dto to product", http.StatusInternalServerError)
		return
	}
	productCreated, err := ph.S.Post(*p)
	if err != nil {
		http.Error(response, "could not Create product", http.StatusInternalServerError)
		return
	}
	response.WriteHeader(http.StatusCreated)
	response.Header().Add("Location", "/products/"+productCreated.Id)
	dto, _ = productCreated.ToDto()
	render.JSON(response, request, dto)
}
func (ph *ProductHandler) Delete(response http.ResponseWriter, request *http.Request) {
	productId := chi.URLParam(request, "productId")
	if productId == "" {
		http.Error(response, "product id is empty", http.StatusBadRequest)
		return
	}
	err := ph.S.Delete(productId)
	if err != nil {
		http.Error(response, "could not Delete product", http.StatusInternalServerError)
		return
	}
	render.NoContent(response, request)

}
