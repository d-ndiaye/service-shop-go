package storeProduct

import "bitbucket.easy.de/users/n.gauche/service-shop-go/internal/product"

type serviceProductStore struct {
	repo Repository
}

type Service interface {
	GetByCategory(categoryName product.Category) ([]StoreProduct, error)
	Get(idStore string, idProduct string) (StoreProduct, error)
	Delete(idStore string, idProduct string) error
	Post(sp StoreProduct) (StoreProduct, error)
	GetAll() ([]StoreProduct, error)
}

func New(repository Repository) Service {
	sp := serviceProductStore{
		repo: repository,
	}
	return sp
}
func (sps serviceProductStore) GetByCategory(categoryName product.Category) ([]StoreProduct, error) {
	category, err := sps.repo.GetByCategory(categoryName)
	if err != nil {
		return []StoreProduct{}, err
	}
	return category, nil
}

func (sps serviceProductStore) Get(idStore string, idProduct string) (StoreProduct, error) {
	s, err := sps.repo.Get(idStore, idProduct)
	if err != nil {
		return StoreProduct{}, err
	}
	return s, nil
}

func (sps serviceProductStore) Delete(idStore string, idProduct string) error {
	err := sps.repo.Delete(idStore, idProduct)
	if err != nil {
		return err
	}
	return nil
}

func (sps serviceProductStore) Post(sp StoreProduct) (StoreProduct, error) {
	stoPro, err := sps.repo.Post(sp)
	if err != nil {
		return StoreProduct{}, err
	}
	return stoPro, nil
}

func (sps serviceProductStore) GetAll() ([]StoreProduct, error) {
	sp, err := sps.repo.GetAll()
	if err != nil {
		return sp, err
	}
	return sp, nil
}
