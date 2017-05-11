package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddingBundlesToItemIndex(t *testing.T) {

	itemIndex := make(ItemIndex)

	itemIndex.MapToUnresolvedItems([]string{"/bundle/123"})

	assert.Equal(t, 1, len(itemIndex["bundle"]), "expect item index to include the bundle")
	assert.Equal(t, Unresolved, itemIndex["bundle"]["/bundle/123"], "expect item to be unresolved")
	assert.Equal(t, false, itemIndex["bundle"]["/bundle/123"].IsResolved(), "expect item to be IsResolved() == false")

	assert.Equal(t, Unresolved, itemIndex.Get("/bundle/123"), "expect item to be Unresolved")
}

func TestLinkingBundlesItems(t *testing.T) {

	itemIndex := make(ItemIndex)

	// add bundle with 2 items
	bundle := Bundle{
		ID:          123,
		Title:       "Bundle One",
		Description: "Bundle description",
		Items:       itemIndex.MapToUnresolvedItems([]string{"/film/2", "/film/123"}),
	}

	assert.Equal(t, "/film/2", bundle.Items[0].Slug, "expect the item to be a reference to /film/2")
	assert.Equal(t, false, bundle.Items[0].IsResolved(), "expect the item to be marked as unresolved")

	// add the film to the itemIndex
	film := Film{
		ID:    2,
		Slug:  "/film/2",
		Title: "Lord of the Rings",
	}
	itemIndex.Set(film.Slug, film.GetGenericItem())

	// create a very small site, 1 bundle and 1 film
	site := Site{
		Bundles: make(BundleCollection, 0),
		Pages:   make(Pages, 0),
		Films:   make(FilmCollection, 0),
	}
	site.Bundles = append(site.Bundles, bundle)
	site.Films = append(site.Films, film)

	// link the items
	site.LinkItems(itemIndex)

	assert.Equal(t, 1, len(site.Bundles), "expect item index to include the bundle")
	assert.Equal(t, "/film/2", site.Bundles[0].Items[0].Slug, "expect item index to include the bundle")
	assert.Equal(t, "Lord of the Rings", site.Bundles[0].Items[0].Title, "expect item index to include the bundle")
	assert.Equal(t, true, site.Bundles[0].Items[0].IsResolved(), "expect the item to be marked as resolved")
}