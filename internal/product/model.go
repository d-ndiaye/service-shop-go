package product

import (
	"bitbucket.easy.de/users/n.gauche/service-shop-go/pkg/dto"
	"github.com/jinzhu/copier"
)

type Product struct {
	Id       string   `bson:"_id"`
	Name     string   `bson:"name"`
	Price    float64  `bson:"price"`
	Note     string   `bson:"note"`
	Category Category `bson:"category"`
}

type Category string

const (
	Undefined Category = ""
	Handy     Category = "handy"
	Laptop    Category = "laptop"
	Monitor   Category = "monitor"
)

func (c Category) ByName(name string) (category Category) {
	switch name {
	case "handy":
		return Handy
	case "laptop":
		return Laptop
	case "monitor":
		return Monitor
	}
	return Undefined
}

func (p *Product) ToDto() (dto dto.Dto, err error) {
	err = copier.Copy(&dto, p)
	return dto, err
}
func (p *Product) FromDto(dto dto.Dto) (err error) {
	err = copier.Copy(p, &dto)
	return err
}
