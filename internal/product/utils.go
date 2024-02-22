package product

const tva float64 = 0.1

func AddTva(product *Product) {
	product.Price += tva * product.Price
}
