package dto

type StoreProductDto struct {
	ProductID string `json:"productID"`
	StoreID   string `json:"storeID"`
	Price     int    `json:"price"`
	Quantity  int    `json:"quantity"`
}
