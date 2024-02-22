package storeProduct

import (
	"bitbucket.easy.de/users/n.gauche/service-shop-go/internal/product"
	"fmt"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestServiceProductStore_Get(t *testing.T) {
	mockStoreProduct := StoreProduct{
		ProductID: "02",
		StoreID:   "01",
		Price:     20,
		Quantity:  25,
	}
	repoMock := NewRepositoryMock(t)
	repoMock.EXPECT().Get("01", "02").Return(mockStoreProduct, nil).Once()
	sps := serviceProductStore{
		repo: repoMock,
	}
	p, err := sps.Get("01", "02")
	assert.Equal(t, mockStoreProduct.StoreID, p.StoreID)
	assert.Equal(t, mockStoreProduct.ProductID, p.ProductID)
	assert.Nil(t, err)
}

func TestServiceProductStore_GetAll(t *testing.T) {
	storeProducts1 := StoreProduct{
		ProductID: "02",
		StoreID:   "01",
		Price:     20,
		Quantity:  25,
	}
	storeProducts2 := StoreProduct{
		ProductID: "02",
		StoreID:   "01",
		Price:     10,
		Quantity:  50,
	}
	mockStoreProduct := []StoreProduct{storeProducts1, storeProducts2}
	repoMock := NewRepositoryMock(t)
	repoMock.EXPECT().GetAll().Return(mockStoreProduct, nil).Once()
	sps := serviceProductStore{
		repo: repoMock,
	}
	all, err := sps.GetAll()
	assert.Equal(t, 2, len(all))
	assert.Contains(t, all, storeProducts1)
	assert.Contains(t, all, storeProducts2)
	assert.Nil(t, err)
}

func TestServiceProductStore_Post(t *testing.T) {
	mockStoreProduct := StoreProduct{
		ProductID: "02",
		StoreID:   "01",
		Price:     20,
		Quantity:  25,
	}
	repoMock := NewRepositoryMock(t)
	repoMock.EXPECT().Post(mock.AnythingOfType("StoreProduct")).Run(func(p StoreProduct) {
		assert.Equal(t, 20, p.Price)
	}).Return(mockStoreProduct, nil).Once()
	sps := serviceProductStore{
		repo: repoMock,
	}
	ps, err := sps.Post(mockStoreProduct)
	if ps.ProductID == "" {
		ps.ProductID = fmt.Sprintf("%v", xid.New())
	}
	assert.Equal(t, "02", ps.ProductID)
	assert.Equal(t, "01", ps.StoreID)
	assert.Equal(t, 20, ps.Price)
	assert.Equal(t, 25, ps.Quantity)
	assert.Nil(t, err)
}
func TestServiceProductStore_Delete(t *testing.T) {
	repoMock := NewRepositoryMock(t)
	repoMock.EXPECT().Delete("01", "02").Return(nil).Once()
	sps := serviceProductStore{
		repo: repoMock,
	}
	err := sps.Delete("01", "02")
	assert.Nil(t, err)
}
func TestServiceProductStore_GetByCategory(t *testing.T) {
	var mockStoreProducts []StoreProduct
	mockStoreProducts = make([]StoreProduct, 0, 2)
	storeProducts1 := StoreProduct{
		ProductID: "02",
		StoreID:   "01",
		Price:     20,
		Quantity:  25,
	}
	storeProducts2 := StoreProduct{
		ProductID: "02",
		StoreID:   "01",
		Price:     10,
		Quantity:  50,
	}
	mockStoreProducts = append(mockStoreProducts, storeProducts1, storeProducts2)
	repoMock := NewRepositoryMock(t)
	repoMock.EXPECT().GetByCategory(product.Handy).Return(mockStoreProducts, nil).Once()
	sps := serviceProductStore{
		repo: repoMock,
	}
	categoryList, err := sps.GetByCategory(product.Handy)
	assert.Equal(t, 2, len(categoryList))
	assert.Nil(t, err)
}
func TestServiceProductStore_GetByCategory_Error(t *testing.T) {
	mockProductList := []StoreProduct{}
	repoMock := NewRepositoryMock(t)
	repoMock.EXPECT().GetByCategory(product.Handy).Return(mockProductList, nil).Once()
	sS := serviceProductStore{
		repo: repoMock,
	}
	_, err := sS.GetByCategory(product.Handy)
	assert.Nil(t, err)
}

func TestServiceProductStore_Get_Error(t *testing.T) {
	mockProductList := StoreProduct{}
	repoMock := NewRepositoryMock(t)
	repoMock.EXPECT().Get("01", "02").Return(mockProductList, nil).Once()
	sS := serviceProductStore{
		repo: repoMock,
	}
	s, err := sS.Get("01", "02")
	assert.Equal(t, "", s.ProductID)
	assert.Nil(t, err)
}

func TestServiceProductStore_GetAll_Error(t *testing.T) {
	mockStoreList := []StoreProduct{}
	repoMock := NewRepositoryMock(t)
	repoMock.EXPECT().GetAll().Return(mockStoreList, nil).Once()
	sS := serviceProductStore{
		repo: repoMock,
	}
	s, err := sS.GetAll()
	assert.Equal(t, 0, len(s))
	assert.Nil(t, err)
}

func TestServiceProductStore_Delete_Error(t *testing.T) {
	repoMock := NewRepositoryMock(t)
	repoMock.EXPECT().Delete("52", "52").Return(nil).Once()
	sS := serviceProductStore{
		repo: repoMock,
	}
	err := sS.Delete("52", "52")
	assert.Nil(t, err)
}
