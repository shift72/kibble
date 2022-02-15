package models

import (
	decimal "github.com/shopspring/decimal"
)

const (
	// Rent - const
	Rent = "rent"
	// Buy - const
	Buy = "buy"
	// SD - const
	SD = "sd"
	// HD - const
	HD = "hd"
)

// PriceInfo - store the price and currency
type PriceInfo struct {
	Currency   string
	Prices     PriceCollection
	PlanPrices PlanPriceCollection
}

// PriceCollection -
type PriceCollection []Price

// PlanPriceCollection -
type PlanPriceCollection map[string]PriceCollection

// Price -
type Price struct {
	Ownership   string
	Quality     string
	Price       decimal.Decimal
	PriceString string
}

func (col PriceCollection) find(ownership, quality string) *Price {
	for i := 0; i < len(col); i++ {
		if col[i].Ownership == ownership && col[i].Quality == quality {
			return &col[i]
		}
	}

	return nil
}

func (col PriceCollection) findLowestPrice() *Price {
	var found *Price
	for i := 0; i < len(col); i++ {
		if found == nil || found.Price.GreaterThan(col[i].Price) {
			found = &col[i]
		}
	}

	return found
}

// HasPrice will return true if a ownership and quality match is found
func (p PriceInfo) HasPrice(ownership, quality string) bool {
	return p.Prices.find(ownership, quality) != nil
}

// GetPrice will return the price string
func (p PriceInfo) GetPrice(ownership, quality string) string {
	found := p.Prices.find(ownership, quality)
	if found != nil {
		return found.PriceString
	}
	return ""
}

// GetValue will return the price value
func (p PriceInfo) GetValue(ownership, quality string) decimal.Decimal {
	found := p.Prices.find(ownership, quality)
	if found != nil {
		return found.Price
	}
	return decimal.Zero
}

// GetLowestPrice will return the lowest price
func (p PriceInfo) GetLowestPrice() string {
	found := p.Prices.findLowestPrice()
	if found != nil {
		return found.PriceString
	}
	return ""
}

// GetLowestValue will return the price as value
func (p PriceInfo) GetLowestValue() decimal.Decimal {
	found := p.Prices.findLowestPrice()
	if found != nil {
		return found.Price
	}
	return decimal.Zero
}
