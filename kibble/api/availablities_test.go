package api

import (
	"encoding/json"
	models "kibble/models"

	"kibble/utils"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeAvailabilities(t *testing.T) {

	// site
	site := &models.Site{
		Films: models.FilmCollection{
			"/film/103":{ID: 103,
				Slug: "/film/103",
			},
		},
	}

	// setup index
	itemIndex := make(models.ItemIndex)
	itemIndex.Set(site.Films["/film/103"].Slug, site.Films["/film/103"].GetGenericItem())

	from := utils.ParseTimeFromString("2012-04-01T00:00:00.000Z")
	to := utils.ParseTimeFromString("2012-05-01T00:00:00.000Z")
	availabilities := availabilities{
		{
			Slug: "/film/103",
			From: &from,
			To:   &to,
		},
	}

	// act - load the prices
	count, err := processAvailabilities(availabilities, site, itemIndex)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)

	// verify the film entry is updated
	assert.Equal(t, from, *site.Films["/film/103"].Available.From, "film from was not set")

	// // check the itemIndex is updated
	item := itemIndex.Get("/film/103")
	film, ok := item.InnerItem.(models.Film)
	assert.True(t, ok)
	assert.Equal(t, to, *film.Available.To)
}

func TestDeserializeAvailabilities(t *testing.T) {

	body := `[{"slug":"/film/1586","from":null,"ms_from":null,"to":null,"ms_to":null,"rental_duration_minutes":5760,"rental_playback_duration_minutes":11580},{"slug":"/tv/213/season/1","from":null,"ms_from":null,"to":null,"ms_to":null,"rental_duration_minutes":5760,"rental_playback_duration_minutes":11580},{"slug":"/tv/213/season/1/bonus/1","from":null,"ms_from":null,"to":null,"ms_to":null,"rental_duration_minutes":5760,"rental_playback_duration_minutes":11580}]
	`
	var d availabilities
	err := json.Unmarshal([]byte(body), &d)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(d))
}
