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
					PageCollection{
						ID: 11,
					},
					PageCollection{
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
					PageCollection{
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
			Language{
				Code:               "en",
				Name:               "English",
				Locale:             "nz",
				DefinitionFilePath: "/wut/wut",
				IsDefault:          true,
			},
			Language{
				Code:               "fr",
				Name:               "French",
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

	assert.Equal(t, string(b), "[{\"code\":\"en\",\"name\":\"English\",\"locale\":\"nz\"},{\"code\":\"fr\",\"name\":\"French\",\"locale\":\"FR\"}]")
}
