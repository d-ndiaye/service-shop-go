package dto

type Dto struct {
	Id       string  `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Note     string  `json:"note"`
	Category string  `json:"category"`
}
