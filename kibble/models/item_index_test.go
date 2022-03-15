//    Copyright 2018 SHIFT72
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

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

func TestAddingTVSeasonToItemIndex(t *testing.T) {

	itemIndex := make(ItemIndex)

	itemIndex.MapToUnresolvedItems([]string{"/tv/123/season/456"})

	assert.Equal(t, 1, len(itemIndex["tv-season"]), "expect item index to include the tv season")
	assert.Equal(t, Unresolved, itemIndex["tv-season"]["/tv/123/season/456"], "expect item to be unresolved")
	assert.Equal(t, false, itemIndex["tv-season"]["/tv/123/season/456"].IsResolved(), "expect item to be IsResolved() == false")
	assert.Equal(t, Unresolved, itemIndex.Get("/tv/123/season/456"), "expect item to be Unresolved")
}

func TestAddingTVShowToItemIndex(t *testing.T) {

	itemIndex := make(ItemIndex)

	itemIndex.MapToUnresolvedItems([]string{"/tv/123"})

	assert.Equal(t, 1, len(itemIndex["tv"]), "expect item index to include the tv show")
	assert.Equal(t, Unresolved, itemIndex["tv"]["/tv/123"], "expect item to be unresolved")
	assert.Equal(t, false, itemIndex["tv"]["/tv/123"].IsResolved(), "expect item to be IsResolved() == false")
	assert.Equal(t, Unresolved, itemIndex.Get("/tv/123"), "expect item to be Unresolved")
}

func TestAddingPlanToItemIndex(t *testing.T) {

	itemIndex := make(ItemIndex)

	itemIndex.MapToUnresolvedItems([]string{"/plan/123"})

	assert.Equal(t, 1, len(itemIndex["plan"]), "expect item index to include the plan")
	assert.Equal(t, Unresolved, itemIndex["plan"]["/plan/123"], "expect item to be unresolved")
	assert.Equal(t, false, itemIndex["plan"]["/plan/123"].IsResolved(), "expect item to be IsResolved() == false")
	assert.Equal(t, Unresolved, itemIndex.Get("/plan/123"), "expect item to be Unresolved")
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


	site.Films["/film/103"] = &Film{
		ID: 103,
		Slug: "/film/103",
	}

	// link the items
	site.LinkItems(itemIndex)

	assert.Equal(t, 1, len(site.Bundles), "expect item index to include the bundle")
	assert.Equal(t, "/film/2", site.Bundles[0].Items[0].Slug, "expect item index to include the bundle")
	assert.Equal(t, "Lord of the Rings", site.Bundles[0].Items[0].Title, "expect item index to include the bundle")
	assert.Equal(t, true, site.Bundles[0].Items[0].IsResolved(), "expect the item to be marked as resolved")
}
