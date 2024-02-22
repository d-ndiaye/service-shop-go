package api

import (
	"bitbucket.easy.de/users/n.gauche/service-shop-go/internal/product"
	"bitbucket.easy.de/users/n.gauche/service-shop-go/internal/storeProduct"
	"bitbucket.easy.de/users/n.gauche/service-shop-go/pkg/dto"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStoreProductHandler_Get(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/01/product/02", nil)
	mockStoreProductList := storeProduct.StoreProduct{
		ProductID: "02",
		StoreID:   "01",
		Price:     20,
		Quantity:  25,
	}
	repoMock := storeProduct.NewRepositoryMock(t)
	repoMock.EXPECT().Get(mockStoreProductList.StoreID, mockStoreProductList.ProductID).Return(mockStoreProductList, nil).Once()
	rr := httptest.NewRecorder()
	storeProductHandler := StoreProductHandler{
		Sp: repoMock,
	}
	handler := storeProductHandler.Routes()
	handler.ServeHTTP(rr, req)
	var response dto.StoreProductDto
	err := json.Unmarshal([]byte(rr.Body.String()), &response)
	if err != nil {
		fmt.Println("Parse JSON Data Error")
	}
	assert.Equal(t, 200, rr.Code)
	assert.Equal(t, mockStoreProductList.StoreID, response.StoreID)
	assert.Equal(t, mockStoreProductList.ProductID, response.ProductID)
	assert.Nil(t, err)
}

func TestStoreProductHandler_GetAll(t *testing.T) {
	storeProducts1 := storeProduct.StoreProduct{
		ProductID: "02",
		StoreID:   "01",
		Price:     20,
		Quantity:  25,
	}
	storeProducts2 := storeProduct.StoreProduct{
		ProductID: "02",
		StoreID:   "01",
		Price:     10,
		Quantity:  50,
	}
	mockStoreProductList := []storeProduct.StoreProduct{storeProducts1, storeProducts2}
	req := httptest.NewRequest(http.MethodGet, "/52/product/", nil)
	repoMock := storeProduct.NewRepositoryMock(t)
	repoMock.EXPECT().GetAll().Return(mockStoreProductList, nil).Once()
	rr := httptest.NewRecorder()
	storeProductHandler := StoreProductHandler{
		Sp: repoMock,
	}
	handler := storeProductHandler.Routes()
	handler.ServeHTTP(rr, req)
	var response []dto.StoreProductDto
	err := json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Unexpected result returned: Could not parse response to JSON: '%s'", err)
	}
	assert.Equal(t, 200, rr.Code)
	assert.Equal(t, 2, len(response))
	assert.Nil(t, err)

}

func TestStoreProductHandler_GetByCategory(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/52/product/category/handy", nil)
	var mockStoreProducts []storeProduct.StoreProduct
	mockStoreProducts = make([]storeProduct.StoreProduct, 0, 2)
	storeProducts1 := storeProduct.StoreProduct{
		ProductID: "02",
		StoreID:   "01",
		Price:     20,
		Quantity:  25,
	}
	storeProducts2 := storeProduct.StoreProduct{
		ProductID: "02",
		StoreID:   "01",
		Price:     10,
		Quantity:  50,
	}
	mockStoreProducts = append(mockStoreProducts, storeProducts1, storeProducts2)
	repoMock := storeProduct.NewRepositoryMock(t)
	repoMock.EXPECT().GetByCategory(product.Handy).Return(mockStoreProducts, nil).Once()
	rr := httptest.NewRecorder()
	storeProductHandler := StoreProductHandler{
		Sp: repoMock,
	}
	handler := storeProductHandler.Routes()
	handler.ServeHTTP(rr, req)
	var response []dto.StoreProductDto
	err := json.Unmarshal([]byte(rr.Body.String()), &response)
	if err != nil {
		fmt.Println("Parse JSON Data Error")
	}
	assert.Equal(t, 200, rr.Code)
	assert.Equal(t, 2, len(response))
	assert.Nil(t, err)
}

func TestStoreProductHandler_Delete(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/52/product/52", nil)
	repoMock := storeProduct.NewRepositoryMock(t)
	repoMock.EXPECT().Delete("52", "52").Return(nil).Once()
	rr := httptest.NewRecorder()
	storeProductHandler := StoreProductHandler{
		Sp: repoMock,
	}
	handler := storeProductHandler.Routes()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, 204, rr.Code)
}

func TestStoreProductHandler_Post(t *testing.T) {
	mockStoreProductList := storeProduct.StoreProduct{
		ProductID: "02",
		StoreID:   "01",
		Price:     20,
		Quantity:  25,
	}
	sPList := storeProduct.StoreProduct{
		ProductID: "02",
		StoreID:   "01",
	}
	repoMock := storeProduct.NewRepositoryMock(t)
	repoMock.EXPECT().Post(sPList).Return(mockStoreProductList, nil).Once()
	rr := postStoreProduct("01", "02", t, repoMock)
	var response dto.StoreProductDto
	err := json.Unmarshal([]byte(rr.Body.String()), &response)
	if err != nil {
		fmt.Println("Parse JSON Data Error")
	}
	assert.Equal(t, 201, rr.Code)
	assert.Equal(t, mockStoreProductList.ProductID, response.ProductID)
	assert.Equal(t, mockStoreProductList.StoreID, response.StoreID)
	assert.Nil(t, err)
}

func prepareStoreProductTestServer(mockedService storeProduct.Service, method string, target string, body io.Reader) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, target, body)
	rr := httptest.NewRecorder()
	storeProductHandler := StoreProductHandler{Sp: mockedService}
	handler := storeProductHandler.Routes()
	handler.ServeHTTP(rr, req)
	return rr
}

func postStoreProduct(storeID string, productID string, t *testing.T, storeProductService storeProduct.Service) *httptest.ResponseRecorder {
	params := map[string]string{
		"storeID":   storeID,
		"productID": productID,
	}
	jsonReader := parseJSONBody(t, params)
	return prepareStoreProductTestServer(storeProductService, http.MethodPost, "/{storeID}/product/{productID}", jsonReader)
}
