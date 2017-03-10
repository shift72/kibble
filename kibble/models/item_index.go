package models

import (
	"fmt"
	"strings"
)

// Add - an item
func (itemIndex ItemIndex) Add(slug string, item GenericItem) {

	slugParts := strings.Split(slug, "/")
	slugType := slugParts[1]

	// item[slugType]
	t, ok := itemIndex[slugType]
	if !ok {
		itemIndex[slugType] = make(map[string]GenericItem)
		t = itemIndex[slugType]
	}

	if _, ok := t[slug]; !ok {
		t[slug] = item
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

// IndexItems - add an index
func (site *Site) IndexItems(itemIndex ItemIndex) {

	for _, f := range site.Films {
		for _, slug := range f.Recommendations {
			itemIndex.Add(slug, Empty)
		}
	}

	for _, p := range site.Pages {
		for _, pf := range p.PageFeatures {
			for _, slug := range pf.Items {
				itemIndex.Add(slug, Empty)
			}
		}
	}
}

// LinkItems - link the items to the specific parts
func (site *Site) LinkItems(itemIndex ItemIndex) {

	for i := range site.Films {
		for _, slug := range site.Films[i].Recommendations {
			t := itemIndex.Get(slug)
			if t != Empty {
				site.Films[i].ResolvedRecommendations = append(site.Films[i].ResolvedRecommendations, t)
			}
		}
	}

	for i := range site.Pages {
		for j := range site.Pages[i].PageFeatures {
			for _, slug := range site.Pages[i].PageFeatures[j].Items {
				t := itemIndex.Get(slug)
				if t != Empty {
					site.Pages[i].PageFeatures[j].ResolvedItems = append(site.Pages[i].PageFeatures[j].ResolvedItems, t)
				}
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
				fmt.Printf("%s - %s : nil\n", t, k)
			} else {
				fmt.Printf("%s - %s : set\n", t, k)
			}
		}
	}
}

// PrintStats - print the stats about the index
func (itemIndex ItemIndex) PrintStats() {
	fmt.Println("--- item index stats ---")

	var loadedCount = 0
	var totalCount = 0
	for t, val := range itemIndex {
		var count = 0
		var loaded = 0
		for _, v := range val {
			count++
			totalCount++
			if v != Empty {
				loaded++
				loadedCount++
			}
		}
		fmt.Printf("type: %s\t%d/%d\n", t, loaded, count)
	}

	fmt.Printf("total: \t\t%d/%d\n", loadedCount, totalCount)
}
