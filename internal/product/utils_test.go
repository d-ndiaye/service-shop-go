package product

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddTva(t *testing.T) {
	p := Product{
		Id:    "001",
		Price: 2,
	}
	AddTva(&p)
	assert.Equal(t, 2.2, p.Price)

}
