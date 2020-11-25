package models

import (
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
