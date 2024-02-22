package api

import (
	"bitbucket.easy.de/users/n.gauche/service-shop-go/internal/product"
	"bitbucket.easy.de/users/n.gauche/service-shop-go/pkg/dto"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	http "net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestProductHandler_GetAll(t *testing.T) {
	product1 := product.Product{
		Id:    "001",
		Name:  "apple",
		Price: 10,
		Note:  "black",
	}
	product2 := product.Product{
		Id:    "002",
		Name:  "samsung",
		Price: 20,
		Note:  "white",
	}
	mockProductList := []product.Product{product1, product2}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	repoMock := product.NewRepositoryMock(t)
	repoMock.EXPECT().GetAll().Return(mockProductList, nil).Once()
	rr := httptest.NewRecorder()
	productHandler := ProductHandler{
		S: repoMock,
	}
	handler := productHandler.Routes()
	handler.ServeHTTP(rr, req)
	var response []dto.Dto
	err := json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Unexpected result returned: Could not parse response to JSON: '%s'", err)
	}
	assert.Equal(t, 200, rr.Code)
	assert.Equal(t, 2, len(response))
	assert.Nil(t, err)
}

func TestProductHandler_Get(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/5", nil)
	mockProductList := product.Product{
		Id:    "5",
		Name:  "apple",
		Price: 10,
		Note:  "black",
	}
	repoMock := product.NewRepositoryMock(t)
	repoMock.EXPECT().Get("5").Return(mockProductList, nil).Once()
	rr := httptest.NewRecorder()
	productHandler := ProductHandler{
		S: repoMock,
	}
	handler := productHandler.Routes()
	handler.ServeHTTP(rr, req)
	var dto dto.Dto
	err := json.Unmarshal([]byte(rr.Body.String()), &dto)
	if err != nil {
		fmt.Println("Parse JSON Data Error")
	}
	assert.Equal(t, 200, rr.Code)
	assert.Equal(t, mockProductList.Price, dto.Price)
	assert.Equal(t, mockProductList.Name, dto.Name)
	assert.Equal(t, mockProductList.Note, dto.Note)
}

func TestProductHandler_Delete(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/5", nil)
	rr := httptest.NewRecorder()
	repoMock := product.NewRepositoryMock(t)
	repoMock.EXPECT().Delete("5").Return(nil).Once()
	productHandler := ProductHandler{
		S: repoMock,
	}
	handler := productHandler.Routes()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, 204, rr.Code)
}

func TestProductHandler_Post(t *testing.T) {
	mockProduct := product.Product{
		Id:    "001",
		Name:  "apple",
		Price: 22,
		Note:  "white",
	}
	plist := product.Product{
		Id:   "001",
		Name: "apple",
	}
	repoMock := product.NewRepositoryMock(t)
	repoMock.EXPECT().Post(plist).Return(mockProduct, nil).Once()
	rr := postProduct("001", "apple", t, repoMock)
	var dto dto.Dto
	err := json.Unmarshal([]byte(rr.Body.String()), &dto)
	if err != nil {
		fmt.Println("Parse JSON Data Error")
	}
	assert.Equal(t, 201, rr.Code)
	assert.Equal(t, mockProduct.Id, dto.Id)
	assert.Equal(t, mockProduct.Name, dto.Name)
	assert.Nil(t, err)
}

func parseJSONBody(t *testing.T, input interface{}) (reader *bytes.Reader) {
	jsonBody, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}
	reader = bytes.NewReader(jsonBody)
	return
}

func prepareProductTestServer(mockedService product.Service, method string, target string, body io.Reader) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, target, body)
	rr := httptest.NewRecorder()
	productHandler := ProductHandler{S: mockedService}
	handler := productHandler.Routes()
	handler.ServeHTTP(rr, req)
	return rr
}

func postProduct(productId string, productName string, t *testing.T, productService product.Service) *httptest.ResponseRecorder {
	params := map[string]string{"id": productId, "name": productName}
	jsonReader := parseJSONBody(t, params)
	return prepareProductTestServer(productService, http.MethodPost, "/", jsonReader)
}

func assertBodyEquality(t *testing.T, expectedObj interface{}, rr *httptest.ResponseRecorder) {
	responseBody := decode(t, rr.Body.Bytes())
	expBytes, _ := json.Marshal(expectedObj)

	if !reflect.DeepEqual(responseBody, decode(t, expBytes)) {
		t.Errorf("\nGetProperties: \nexpected: '%v' \nreceived: '%v'", expectedObj, rr.Body.String())
	}
}

func decode(t *testing.T, body []byte) *interface{} {
	var res interface{}

	err := json.Unmarshal(body, &res)
	if err != nil {
		t.Fatalf("could not parse response to JSON: '%s'", err)
	}
	return &res
}
