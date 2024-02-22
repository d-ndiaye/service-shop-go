package storeProduct

import (
	"bitbucket.easy.de/users/n.gauche/service-shop-go/pkg/dto"
	"github.com/jinzhu/copier"
)

type StoreProduct struct {
	ProductID string `bson:"productID"`
	StoreID   string `bson:"storeID"`
	Price     int    `bson:"price"`
	Quantity  int    `bson:"quantity"`
}

func (sp *StoreProduct) ToDto() (dto dto.StoreProductDto, err error) {
	err = copier.Copy(&dto, sp)
	return dto, err
}
func (sp *StoreProduct) FromDto(dto dto.StoreProductDto) (err error) {
	err = copier.Copy(sp, &dto)
	return err
}
