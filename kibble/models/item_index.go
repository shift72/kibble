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
	"fmt"
	"strings"
)

var (
	// Empty - Generic Item
	Empty = GenericItem{Slug: "empty"}
	// Unresolved - a slug to an item that has not been requested
	Unresolved = GenericItem{Slug: "unresolved"}
)

// ItemIndex - an item index
type ItemIndex map[string]map[string]GenericItem

// IsResolved -
func (genericItem GenericItem) IsResolved() bool {
	return genericItem.InnerItem != nil
}

//
// ItemIndex contains
//   GenericItems - has been found and indexed
//   Empty - could not be loaded
//   Unresolved - found but has not been requested
//

// MapToUnresolvedItems - create an array of unresolved items from an array of slugs
func (itemIndex ItemIndex) MapToUnresolvedItems(items []string) GenericItems {

	genericItems := make(GenericItems, len(items))

	for i, slug := range items {
		genericItems[i].Slug = slug
		itemIndex.Set(slug, Unresolved)
	}

	return genericItems
}

func getSlugType(slug string) string {
	slugParts := strings.Split(slug, "/")

	// film bonus
	if len(slugParts) == 5 {
		return slugParts[1] + "-" + slugParts[3]
	}

	// film
	return slugParts[1]
}

// Set - an item
func (itemIndex ItemIndex) Set(slug string, item GenericItem) {

	slugType := getSlugType(slug)

	index, ok := itemIndex[slugType]
	if !ok {
		itemIndex[slugType] = make(map[string]GenericItem)
		index = itemIndex[slugType]
	}

	// unresolved can be overwritten, empty and item can not
	foundItem, ok := index[slug]
	if (ok && (foundItem == Unresolved || foundItem == Empty)) || !ok {
		index[slug] = item
	}
}

// Replace a value in the index
func (itemIndex ItemIndex) Replace(slug string, item GenericItem) {

	slugType := getSlugType(slug)

	index, ok := itemIndex[slugType]
	if !ok {
		itemIndex[slugType] = make(map[string]GenericItem)
		index = itemIndex[slugType]
	}

	index[slug] = item
}

// Get - get the slug
func (itemIndex ItemIndex) Get(slug string) (item GenericItem) {

	slugType := getSlugType(slug)

	t, ok := itemIndex[slugType]
	if ok {
		return t[slug]
	}

	return Empty
}

// FindEmptySlugs - find the slugs that are missing
func (itemIndex ItemIndex) FindEmptySlugs(slugType string) []string {
	return itemIndex.findSlugsOfType(slugType, Empty)
}

// FindUnresolvedSlugs - find unresolved slugs
func (itemIndex ItemIndex) FindUnresolvedSlugs(slugType string) []string {
	return itemIndex.findSlugsOfType(slugType, Unresolved)
}

func (itemIndex ItemIndex) findSlugsOfType(slugType string, itemType GenericItem) []string {
	found := make([]string, 0)
	t, ok := itemIndex[slugType]
	if ok {
		for k, v := range t {
			fmt.Println(" >>>>>", k, v.Title, v.Slug)
		}
		for k, v := range t {

			if v == itemType {
				found = append(found, k)
			}
		}
	}
	return found
}

// LinkItems - link the items to the specific parts
func (site *Site) LinkItems(itemIndex ItemIndex) {

	for _, f := range site.Films {
		if film, ok := site.Films[f.Slug]; ok {
			film.Recommendations = itemIndex.Resolve(f.Recommendations)
			site.Films[f.Slug] = film
		}
	}

	for i := range site.TVSeasons {
		site.TVSeasons[i].Recommendations = itemIndex.Resolve(site.TVSeasons[i].Recommendations)
	}

	for i := range site.Pages {
		for j := range site.Pages[i].PageCollections {
			site.Pages[i].PageCollections[j].Items = itemIndex.Resolve(site.Pages[i].PageCollections[j].Items)
		}
	}

	for i := range site.Bundles {
		site.Bundles[i].Items = itemIndex.Resolve(site.Bundles[i].Items)
	}

	for i := range site.Collections {
		site.Collections[i].Items = itemIndex.Resolve(site.Collections[i].Items)
	}
}

// Resolve - convert an array of generic items to resolved items
func (itemIndex ItemIndex) Resolve(gItems GenericItems) GenericItems {
	resolvedItems := make([]GenericItem, 0)
	for _, item := range gItems {
		t := itemIndex.Get(item.Slug)
		if t.IsResolved() {
			resolvedItems = append(resolvedItems, t)
		}
	}
	return resolvedItems
}

// Print - print the item index
func (itemIndex ItemIndex) Print() {
	for t, val := range itemIndex {
		log.Infof("type: %s", t)
		for k, v := range val {
			if v == Empty {
				log.Infof("%s - %s : empty", t, k)
			} else if v == Unresolved {
				log.Infof("%s - %s : unresolved", t, k)
			} else {
				log.Infof("%s - %s : set", t, k)
			}
		}
	}
}

// PrintStats - print the stats about the index
func (itemIndex ItemIndex) PrintStats() {
	log.Info("item index:")
	var loadedCount = 0
	var totalCount = 0
	for t, val := range itemIndex {
		var count = 0
		var loaded = 0
		for _, v := range val {
			count++
			totalCount++
			if !(v == Empty || v == Unresolved) {
				loaded++
				loadedCount++
			}
		}
		log.Infof("type: %-15s%4d / %d", t, loaded, count)
	}

	log.Infof("total: %18d / %d", loadedCount, totalCount)
}
