package models

const (
	Rent = "rent"
	Buy  = "buy"
	SD   = "sd"
	HD   = "hd"
)

// PriceInfo - store the price and currency
type PriceInfo struct {
	Currency string
	Prices   PriceCollection
}

// PriceCollection -
type PriceCollection []Price

// Price -
type Price struct {
	Ownership   string
	Quality     string
	Price       float32
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
		if found == nil || found.Price > col[i].Price {
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

// GetLowestPrice will return the lowest price
func (p PriceInfo) GetLowestPrice() string {
	found := p.Prices.findLowestPrice()
	if found != nil {
		return found.PriceString
	}
	return ""
}
