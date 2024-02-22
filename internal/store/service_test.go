package store

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestServiceStore_Get(t *testing.T) {
	mockProductList := Store{
		Id:   "01",
		Name: "Easy Leipzig",
	}
	repoMock := NewRepositoryMock(t)
	repoMock.EXPECT().Get("01").Return(mockProductList, nil).Once()
	sS := serviceStore{
		repo: repoMock,
	}
	p, err := sS.Get("01")
	assert.Equal(t, "Easy Leipzig", p.Name)
	assert.Equal(t, "01", p.Id)
	assert.Nil(t, err)
}

func TestServiceStore_GetAll(t *testing.T) {
	store1 := Store{
		Id:   "001",
		Name: "Easy Essen",
		Note: "Open",
	}
	store2 := Store{
		Id:   "002",
		Name: "Easy Leipzig",
		Note: "Brief pause",
	}
	mockProductList := []Store{store1, store2}
	repoMock := NewRepositoryMock(t)
	repoMock.EXPECT().GetAll().Return(mockProductList, nil).Once()
	sS := serviceStore{
		repo: repoMock,
	}
	all, err := sS.GetAll()
	assert.Equal(t, 2, len(all))
	assert.Contains(t, all, store1)
	assert.Contains(t, all, store2)
	assert.Nil(t, err)
}

func TestServiceStore_Delete(t *testing.T) {
	repoMock := NewRepositoryMock(t)
	repoMock.EXPECT().Delete("01").Return(nil).Once()
	sS := serviceStore{
		repo: repoMock,
	}
	err := sS.Delete("01")
	assert.Nil(t, err)
}

func TestServiceStore_Post(t *testing.T) {
	mockStoreList := Store{
		Id:   "001",
		Name: "apple",
		Note: "white",
	}
	slist := Store{}
	repoMock := NewRepositoryMock(t)
	repoMock.EXPECT().Post(slist).Return(mockStoreList, nil).Once()
	sS := serviceStore{
		repo: repoMock,
	}
	s, err := sS.Post(slist)
	assert.Equal(t, "001", s.Id)
	assert.Equal(t, "apple", s.Name)
	assert.Equal(t, "white", s.Note)
	assert.Nil(t, err)
}

func TestServiceStore_GetAll_Error(t *testing.T) {
	mockStoreList := []Store{}
	repoMock := NewRepositoryMock(t)
	repoMock.EXPECT().GetAll().Return(mockStoreList, nil).Once()
	sS := serviceStore{
		repo: repoMock,
	}
	s, err := sS.GetAll()
	assert.Equal(t, 0, len(s))
	assert.Nil(t, err)
}

func TestServiceStore_Delete_Error(t *testing.T) {
	repoMock := NewRepositoryMock(t)
	repoMock.EXPECT().Delete("52").Return(nil).Once()
	sS := serviceStore{
		repo: repoMock,
	}
	err := sS.Delete("52")
	assert.Nil(t, err)
}
func TestServiceStore_Get_Error(t *testing.T) {
	mockProductList := Store{}
	repoMock := NewRepositoryMock(t)
	repoMock.EXPECT().Get("01").Return(mockProductList, nil).Once()
	sS := serviceStore{
		repo: repoMock,
	}
	s, err := sS.Get("01")
	assert.Equal(t, "", s.Id)
	assert.Nil(t, err)
}
