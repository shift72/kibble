package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryingFromMultiplePrices(t *testing.T) {
	p := &PriceInfo{
		Currency: "nz",
		Prices: PriceCollection{
			{
				Ownership:   Rent,
				Quality:     SD,
				Price:       4.99,
				PriceString: "$4.99",
			},
			{
				Ownership:   Rent,
				Quality:     HD,
				Price:       7.99,
				PriceString: "$7.99",
			},
		},
	}

	assert.Equal(t, "$4.99", p.GetLowestPrice())
	assert.Equal(t, "$4.99", p.GetPrice(Rent, SD))
	assert.Equal(t, "$7.99", p.GetPrice(Rent, HD))
	assert.Equal(t, "", p.GetPrice(Buy, HD))
}
