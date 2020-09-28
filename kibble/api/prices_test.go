package api

import (
	"encoding/json"
	"kibble/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func str(s string) *string {
	return &s
}

func TestMappingApiPrices(t *testing.T) {

	p := pricesV2{
		Item:     "/film/103",
		Currency: "NZD",
		Rent: &qualityPriceV2{
			Hd:       str("5.0"),
			HdString: str("$5.00"),
			Sd:       str("3.0"),
			SdString: str("$3.00"),
		},
		Buy: &qualityPriceV2{
			Hd:       str("9.0"),
			HdString: str("$9.00"),
			Sd:       str("6.0"),
			SdString: str("$6.00"),
		},
	}

	pp := p.getPrices()

	assert.Equal(t, "$3.00", pp.GetLowestPrice())
}

func TestMergePrices(t *testing.T) {

	// site
	site := &models.Site{
		Films: []models.Film{
			{ID: 103,
				Slug: "/film/103",
			},
		},
	}

	// setup index
	itemIndex := make(models.ItemIndex)
	itemIndex.Set(site.Films[0].Slug, site.Films[0].GetGenericItem())

	prices := prices{
		Prices: []pricesV2{
			{
				Item:     "/film/103",
				Currency: "NZD",
				Rent: &qualityPriceV2{
					Hd:       str("5.0"),
					HdString: str("$5.00"),
					Sd:       str("3.0"),
					SdString: str("$3.00"),
				},
				Buy: &qualityPriceV2{
					Hd:       str("9.0"),
					HdString: str("$9.00"),
				},
			},
		},
	}

	// act - load the prices
	count, err := processPrices(prices, site, itemIndex)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)

	// verify the film entry is updated
	assert.Equal(t, "$3.00", site.Films[0].Prices.GetLowestPrice(), "film price was not updated")

	// check the itemIndex is updated
	item := itemIndex.Get("/film/103")
	film, ok := item.InnerItem.(models.Film)
	assert.True(t, ok)
	assert.Equal(t, "$3.00", film.Prices.GetLowestPrice())
}

func TestDeserializePrices(t *testing.T) {

	body := `{"prices":[{"item":"/film/103","currency":"NZD","rent":{"hd":null,"hd_string":null,"sd":null,"sd_string":null},"buy":{"hd":"3.0","hd_string":"$3.00","sd":null,"sd_string":null}},{"item":"/film/104","currency":"NZD","rent":{"hd":"5.0","hd_string":"$5.00","sd":"2.0","sd_string":"$2.00"},"buy":{"hd":"3.0","hd_string":"$3.00","sd":"10.0","sd_string":"$10.00"}}],"plans":[]}`

	var details prices
	err := json.Unmarshal([]byte(body), &details)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(details.Prices))
}
