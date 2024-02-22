package api

import (
	"bitbucket.easy.de/users/n.gauche/service-shop-go/internal/store"
	"bitbucket.easy.de/users/n.gauche/service-shop-go/pkg/dto"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStoreHandler_Delete(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/5", nil)
	rr := httptest.NewRecorder()
	repoMock := store.NewRepositoryMock(t)
	repoMock.EXPECT().Delete("5").Return(nil).Once()
	storeHandler := StoreHandler{
		S: repoMock,
	}
	handler := storeHandler.Routes()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, 204, rr.Code)
}

func TestStoreHandler_Get(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/05", nil)
	mockProductList := store.Store{
		Id:   "05",
		Name: "Easy Essen",
		Note: "Open",
	}
	repoMock := store.NewRepositoryMock(t)
	repoMock.EXPECT().Get("05").Return(mockProductList, nil).Once()
	rr := httptest.NewRecorder()
	recordHandler := StoreHandler{S: repoMock}
	handler := recordHandler.Routes()
	handler.ServeHTTP(rr, req)
	var dto dto.StoreDto
	err := json.Unmarshal([]byte(rr.Body.String()), &dto)
	if err != nil {
		fmt.Println("Parse JSON Data Error")
	}
	assert.Equal(t, 200, rr.Code)
	assert.Equal(t, mockProductList.Id, dto.Id)
	assert.Equal(t, mockProductList.Name, dto.Name)
	assert.Equal(t, mockProductList.Note, dto.Note)
}

func TestStoreHandler_GetAll(t *testing.T) {
	store1 := store.Store{
		Id:   "001",
		Name: "Easy Essen",
		Note: "Open",
	}
	store2 := store.Store{
		Id:   "002",
		Name: "Easy Leipzig",
		Note: "Brief pause",
	}
	mockStoreList := []store.Store{store1, store2}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	repoMock := store.NewRepositoryMock(t)
	repoMock.EXPECT().GetAll().Return(mockStoreList, nil).Once()
	rr := httptest.NewRecorder()
	storeHandler := StoreHandler{
		S: repoMock,
	}

	handler := storeHandler.Routes()
	handler.ServeHTTP(rr, req)
	var response []dto.StoreDto
	err := json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Unexpected result returned: Could not parse response to JSON: '%s'", err)
	}
	assert.Equal(t, 200, rr.Code)
	assert.Equal(t, 2, len(response))
	assert.Nil(t, err)
}

func TestStoreHandler_Post(t *testing.T) {
	mockStoreList := store.Store{
		Id:   "001",
		Name: "Easy Essen",
		Note: "Open",
	}
	storeList := store.Store{
		Id:   "001",
		Name: "Easy Essen",
		Note: "Open",
	}
	repoMock := store.NewRepositoryMock(t)
	repoMock.EXPECT().Post(storeList).Return(mockStoreList, nil).Once()
	rr := postStore("001", "Easy Essen", "Open", t, repoMock)
	var response dto.StoreDto
	err := json.Unmarshal([]byte(rr.Body.String()), &response)
	if err != nil {
		fmt.Println("Parse JSON Data Error")
	}
	assert.Equal(t, 201, rr.Code)
	assert.Equal(t, mockStoreList.Id, response.Id)
	assert.Equal(t, mockStoreList.Name, response.Name)
	assert.Equal(t, mockStoreList.Note, response.Note)
	assert.Nil(t, err)
}

func prepareStoreTestServer(mockedService store.Service, method string, target string, body io.Reader) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, target, body)
	rr := httptest.NewRecorder()
	storeHandler := StoreHandler{S: mockedService}
	handler := storeHandler.Routes()
	handler.ServeHTTP(rr, req)
	return rr
}

func postStore(storeId string, storeName string, storeNote string, t *testing.T, storeService store.Service) *httptest.ResponseRecorder {
	params := map[string]string{"id": storeId, "name": storeName, "note": storeNote}
	jsonReader := parseJSONBody(t, params)
	return prepareStoreTestServer(storeService, http.MethodPost, "/", jsonReader)
}
