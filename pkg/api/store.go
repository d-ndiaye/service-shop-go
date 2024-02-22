package api

// StoreInternal why this link and its parameters
import (
	StoreInternal "bitbucket.easy.de/users/n.gauche/service-shop-go/internal/store"
	"bitbucket.easy.de/users/n.gauche/service-shop-go/pkg/dto"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
)

type StoreHandler struct {
	S StoreInternal.Service
}

func (sh *StoreHandler) Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/", sh.GetAll)
	router.Get("/{storeId}", sh.Get)
	router.Post("/", sh.Post)
	router.Delete("/{storeId}", sh.Delete)
	return router
}

func (sh *StoreHandler) GetAll(response http.ResponseWriter, request *http.Request) {
	stores, err := sh.S.GetAll()
	if err != nil {
		http.Error(response, "could not retrieve the stores", http.StatusInternalServerError)
		return
	}
	render.JSON(response, request, stores)
}

func (sh *StoreHandler) Get(response http.ResponseWriter, request *http.Request) {
	storeId := chi.URLParam(request, "storeId")
	s, err := sh.S.Get(storeId)
	if err != nil {
		http.Error(response, "could not retrieve store", http.StatusInternalServerError)
		return
	}
	dto, err := s.ToDto()
	if err != nil {
		http.Error(response, "Failed to convert store model to dto", http.StatusInternalServerError)
		return
	}
	render.JSON(response, request, dto)
}

func (sh *StoreHandler) decodeBody(request *http.Request, target interface{}) (err error) {
	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(target)
	return
}

func (sh *StoreHandler) Post(response http.ResponseWriter, request *http.Request) {
	dto := dto.StoreDto{}
	if err := sh.decodeBody(request, &dto); err != nil {
		http.Error(response, "could not decode body", http.StatusBadRequest)
		return
	}
	s := &StoreInternal.Store{}
	err := s.FromDto(dto)
	if err != nil {
		http.Error(response, "could not Create store", http.StatusInternalServerError)
		return
	}
	storeCreated, err := sh.S.Post(*s)
	if err != nil {
		http.Error(response, "could not Create product", http.StatusInternalServerError)
		return
	}
	response.WriteHeader(http.StatusCreated)
	response.Header().Add("Location", "/stores/"+storeCreated.Id)
	render.JSON(response, request, dto)
}

func (sh *StoreHandler) Delete(response http.ResponseWriter, request *http.Request) {
	storeId := chi.URLParam(request, "storeId")
	if storeId == "" {
		http.Error(response, "store id is empty", http.StatusBadRequest)
		return
	}
	err := sh.S.Delete(storeId)
	if err != nil {
		http.Error(response, "could not Delete store", http.StatusInternalServerError)
		return
	}
	render.NoContent(response, request)
}
