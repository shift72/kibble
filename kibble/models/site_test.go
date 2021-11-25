package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdatePageCollectionDescriptions(t *testing.T) {

	site := Site{
		Pages: Pages{
			Page{
				ID:        123,
				Slug:      "/film/123",
				TitleSlug: "the-big-lebowski",
				PageCollections: []PageCollection{
					{
						ID: 11,
					},
					{
						ID: 12,
					},
				},
			},
		},
		Collections: CollectionCollection{
			Collection{
				ID:          11,
				Description: "collection eleven",
			},
			Collection{
				ID:          12,
				Description: "collection twelve",
			},
		},
	}

	site.UpdatePageCollections()

	assert.Equal(t, "collection eleven", site.Pages[0].PageCollections[0].Description)
	assert.Equal(t, "collection twelve", site.Pages[0].PageCollections[1].Description)
}

func TestUpdatePageCollectionMissing(t *testing.T) {

	site := Site{
		Pages: Pages{
			Page{
				ID:        123,
				Slug:      "/film/123",
				TitleSlug: "the-big-lebowski",
				PageCollections: []PageCollection{
					{
						ID: 11,
					},
				},
			},
		},
		Collections: CollectionCollection{},
	}

	site.UpdatePageCollections()

	assert.Equal(t, "", site.Pages[0].PageCollections[0].Description)
}

func TestLanguagesConvertToJSON(t *testing.T) {

	site := Site{
		Languages: []Language{
			{
				Code:  "en",
				Name:  "English",
				Label: "nz",
				//Deprecated
				Locale:             "nz",
				DefinitionFilePath: "/wut/wut",
				IsDefault:          true,
			},
			{
				Code:  "fr",
				Name:  "French",
				Label: "FR",
				//Deprecated
				Locale:             "FR",
				DefinitionFilePath: "/oi/oi",
				IsDefault:          false,
			},
		},
	}

	languages := &site.Languages
	b, err := json.Marshal(languages)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, string(b), `[{"code":"en","name":"English","label":"nz","locale":"nz"},{"code":"fr","name":"French","label":"FR","locale":"FR"}]`)
}
