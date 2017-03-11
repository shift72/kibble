package models

import (
	"fmt"
	"strings"
)

//
// ItemIndex contains
//   GenericItems - has been found and indexed
//   Empty - could not be loaded
//   Unresolved - found but has not been requested
//

// Set - an item
func (itemIndex ItemIndex) Set(slug string, item GenericItem) {

	// fmt.Println("adding", slug, item == Empty, item == Unresolved)

	slugParts := strings.Split(slug, "/")
	slugType := slugParts[1]

	index, ok := itemIndex[slugType]
	if !ok {
		itemIndex[slugType] = make(map[string]GenericItem)
		index = itemIndex[slugType]
	}

	// unresolved can be overwritten, empty and item can not
	foundItem, ok := index[slug]
	if (ok && (foundItem == Unresolved || foundItem == Empty)) || !ok {

		// fmt.Println("setting", slug)

		index[slug] = item
	}
}

// Get - get the slug
func (itemIndex ItemIndex) Get(slug string) (item GenericItem) {

	slugParts := strings.Split(slug, "/")
	slugType := slugParts[1]

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
			if v == itemType {
				found = append(found, k)
			}
		}
	}
	return found
}

// LinkItems - link the items to the specific parts
func (site *Site) LinkItems(itemIndex ItemIndex) {

	//TODO: still needed?

	for i := range site.Films {
		for _, slug := range site.Films[i].Recommendations {
			t := itemIndex.Get(slug)
			if t != Unresolved {
				site.Films[i].ResolvedRecommendations = append(site.Films[i].ResolvedRecommendations, t)
			}
		}
	}

	for i := range site.Pages {
		for j := range site.Pages[i].PageFeatures {
			for _, slug := range site.Pages[i].PageFeatures[j].Items {
				t := itemIndex.Get(slug)
				if t != Unresolved {
					site.Pages[i].PageFeatures[j].ResolvedItems = append(site.Pages[i].PageFeatures[j].ResolvedItems, t)
				}
			}
		}
	}

	for i := range site.Bundles {
		for _, slug := range site.Bundles[i].Items {
			t := itemIndex.Get(slug)
			if t != Unresolved {
				site.Bundles[i].ResolvedItems = append(site.Bundles[i].ResolvedItems, t)
			}
		}
	}

}

// Print - print the item index
func (itemIndex ItemIndex) Print() {
	for t, val := range itemIndex {
		fmt.Printf("type: %s\n", t)
		for k, v := range val {
			if v == Empty {
				fmt.Printf("%s - %s : empty\n", t, k)
			} else if v == Unresolved {
				fmt.Printf("%s - %s : unresolved\n", t, k)
			} else {
				fmt.Printf("%s - %s : set\n", t, k)
			}
		}
	}
}

// PrintStats - print the stats about the index
func (itemIndex ItemIndex) PrintStats() {
	fmt.Println("item index:")
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
		fmt.Printf("type: %s\t%d/%d\n", t, loaded, count)
	}

	fmt.Printf("total: \t\t%d/%d\n", loadedCount, totalCount)
}
