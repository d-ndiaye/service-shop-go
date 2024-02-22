package store

import (
	"bitbucket.easy.de/users/n.gauche/service-shop-go/pkg/dto"
	"github.com/jinzhu/copier"
)

type Store struct {
	Id   string `bson:"_id"`
	Name string `bson:"name"`
	Note string `bson:"note"`
}

func (s *Store) ToDto() (dto dto.StoreDto, err error) {
	err = copier.Copy(&dto, s)
	return dto, err
}

func (s *Store) FromDto(dto dto.StoreDto) (err error) {
	err = copier.Copy(s, &dto)
	return err
}
