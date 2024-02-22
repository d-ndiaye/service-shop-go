package product

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestServiceProduct_Get(t *testing.T) {
	mockProduct := Product{
		Id:    "5",
		Name:  "apple",
		Price: 10,
		Note:  "black",
	}
	repoMock := NewRepositoryMock(t)
	repoMock.EXPECT().Get("5").Return(mockProduct, nil).Once()
	sp := ServiceProduct{
		repo: repoMock,
	}
	p, err := sp.Get("5")
	assert.Equal(t, "5", p.Id)
	assert.Equal(t, mockProduct.Name, p.Name)
	assert.Nil(t, err)
}

func TestServiceProduct_GetAll(t *testing.T) {
	product1 := Product{
		Id:    "001",
		Name:  "apple",
		Price: 10,
		Note:  "black",
	}
	product2 := Product{
		Id:    "002",
		Name:  "samsung",
		Price: 20,
		Note:  "white",
	}
	mockProductList := []Product{product1, product2}
	repoMock := NewRepositoryMock(t)
	repoMock.EXPECT().GetAll().Return(mockProductList, nil).Once()
	sp := ServiceProduct{
		repo: repoMock,
	}
	all, err := sp.GetAll()
	assert.Equal(t, 2, len(all))
	assert.Contains(t, all, product1)
	assert.Contains(t, all, product2)
	assert.Nil(t, err)
}

func TestServiceProduct_Delete(t *testing.T) {
	repoMock := NewRepositoryMock(t)
	repoMock.EXPECT().Delete("01").Return(nil).Once()
	sp := ServiceProduct{
		repo: repoMock,
	}
	err := sp.Delete("01")
	assert.Nil(t, err)
}

func TestServiceProduct_Post(t *testing.T) {
	mockProduct := Product{
		Id:    "5",
		Name:  "apple",
		Price: 10,
		Note:  "black",
	}
	repoMock := NewRepositoryMock(t)
	repoMock.EXPECT().Post(mock.AnythingOfType("Product")).Run(func(p Product) {
		assert.Equal(t, 11.0, p.Price)
	}).Return(mockProduct, nil).Once()
	sp := ServiceProduct{
		repo: repoMock,
	}
	p, err := sp.Post(mockProduct)
	assert.Equal(t, "5", p.Id)
	assert.Equal(t, "apple", p.Name)
	assert.Equal(t, 10.0, p.Price)
	assert.Nil(t, err)
}

func TestServiceProduct_GetAll_Error(t *testing.T) {
	mockProductList := []Product{}
	repoMock := NewRepositoryMock(t)
	repoMock.EXPECT().GetAll().Return(mockProductList, nil).Once()
	sp := ServiceProduct{
		repo: repoMock,
	}
	p, err := sp.GetAll()
	assert.Equal(t, 0, len(p))
	assert.Nil(t, err)
}

func TestServiceProduct_Delete_Error(t *testing.T) {
	repoMock := NewRepositoryMock(t)
	repoMock.EXPECT().Delete("52").Return(nil).Once()
	sS := ServiceProduct{
		repo: repoMock,
	}
	err := sS.Delete("52")
	assert.Nil(t, err)
}
func TestServiceProduct_Get_Error(t *testing.T) {
	mockProductList := Product{}
	repoMock := NewRepositoryMock(t)
	repoMock.EXPECT().Get("01").Return(mockProductList, nil).Once()
	sS := ServiceProduct{
		repo: repoMock,
	}
	s, err := sS.Get("01")
	assert.Equal(t, "", s.Id)
	assert.Nil(t, err)
}

func TestServiceProduct_GetAllByName(t *testing.T) {
	product1 := Product{
		Id:    "001",
		Name:  "samsung",
		Price: 10,
		Note:  "black",
	}
	product2 := Product{
		Id:    "002",
		Name:  "apple",
		Price: 20,
		Note:  "white",
	}
	mockProductList := []Product{product1, product2}
	repoMock := NewRepositoryMock(t)
	repoMock.EXPECT().GetAllByName("apple").Return([]Product{product2}, nil).Once()
	sp := ServiceProduct{
		repo: repoMock,
	}
	all, err := sp.GetAllByName("apple")
	assert.Equal(t, 2, len(mockProductList))
	assert.Equal(t, 1, len(all))
	assert.Nil(t, err)
}
