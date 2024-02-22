package api

import (
	"bitbucket.easy.de/users/n.gauche/service-shop-go/internal/product"
	StoreProductInternal "bitbucket.easy.de/users/n.gauche/service-shop-go/internal/storeProduct"
	"bitbucket.easy.de/users/n.gauche/service-shop-go/pkg/dto"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
)

type StoreProductHandler struct {
	Sp StoreProductInternal.Service
}

func (sph *StoreProductHandler) Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/{storeId}/product/category/{categoryName}", sph.GetByCategory)
	router.Get("/{storeId}/product/", sph.GetAll)
	router.Get("/{storeId}/product/{productId}", sph.Get)
	router.Post("/{storeId}/product/{productId}", sph.Post)
	router.Delete("/{storeId}/product/{productId}", sph.Delete)
	return router
}

func (sph *StoreProductHandler) GetAll(response http.ResponseWriter, request *http.Request) {
	storesProducts, err := sph.Sp.GetAll()
	if err != nil {
		http.Error(response, "could not retrieve the products from the stores", http.StatusInternalServerError)
		return
	}
	render.JSON(response, request, storesProducts)
}

func (sph *StoreProductHandler) Get(response http.ResponseWriter, request *http.Request) {
	storeId := chi.URLParam(request, "storeId")
	productId := chi.URLParam(request, "productId")
	s, err := sph.Sp.Get(storeId, productId)
	if err != nil {
		http.Error(response, "could not retrieve the product from the store", http.StatusInternalServerError)
		return
	}
	render.JSON(response, request, s)
}

func (sph *StoreProductHandler) GetByCategory(response http.ResponseWriter, request *http.Request) {
	categoryName := chi.URLParam(request, "categoryName")
	// category check
	var category product.Category
	category = category.ByName(categoryName)

	storeProduct, err := sph.Sp.GetByCategory(category)
	if err != nil {
		http.Error(response, "could not retrieve storeProduct by category", http.StatusInternalServerError)
		return
	}
	dtos := make([]dto.StoreProductDto, 0, len(storeProduct))
	for _, p := range storeProduct {
		dto, _ := p.ToDto()
		dtos = append(dtos, dto)
	}
	render.JSON(response, request, dtos)
}

func (sph *StoreProductHandler) decodeBody(request *http.Request, target interface{}) (err error) {
	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(target)
	return
}

func (sph *StoreProductHandler) Post(response http.ResponseWriter, request *http.Request) {
	dto := StoreProductInternal.StoreProduct{}
	if err := sph.decodeBody(request, &dto); err != nil {
		http.Error(response, "could not decode body", http.StatusBadRequest)
		return
	}

	sp, err := sph.Sp.Post(dto)
	if err != nil {
		http.Error(response, "could not Create the product from the store", http.StatusInternalServerError)
		return
	}
	response.WriteHeader(http.StatusCreated)
	response.Header().Add("Location", "/stores/"+sp.StoreID+"/products/"+sp.ProductID)
	render.JSON(response, request, dto)

}

func (sph *StoreProductHandler) Delete(response http.ResponseWriter, request *http.Request) {
	storeId := chi.URLParam(request, "storeId")
	productId := chi.URLParam(request, "productId")

	if storeId == "" {
		http.Error(response, "store id is empty", http.StatusBadRequest)
		return
	}
	if productId == "" {
		http.Error(response, "product id is empty", http.StatusBadRequest)
		return
	}
	err := sph.Sp.Delete(storeId, productId)
	if err != nil {
		http.Error(response, "could not Delete the product from the store", http.StatusInternalServerError)
		return
	}
	render.NoContent(response, request)
}
