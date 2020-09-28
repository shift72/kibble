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

// Heavily inspired by Hugo

import (
	"fmt"
	"sort"
)

// Taxonomies - list of taxonomies
type Taxonomies map[string]Taxonomy

// Taxonomy - a grouping of related generic items
// e.g. films from "2005"
type Taxonomy map[string]GenericItems

// GenericItems - a list of generic items
// e.g. Films from "2005"
type GenericItems []GenericItem

// This set of structs are used to rank results

// OrderedEntry - a list of items ordered by keys
// e.g. Key: Horror
type OrderedEntry struct {
	Key   string
	Items GenericItems
}

// OrderedEntries - array of entires
type OrderedEntries []OrderedEntry

// PopulateTaxonomyWithFilms - select the taxonomy and the film attribute
func (s Site) PopulateTaxonomyWithFilms(taxonomy string, finder func(*Film) []string) {

	t, ok := s.Taxonomies[taxonomy]

	if !ok {
		t = make(Taxonomy)
		s.Taxonomies[taxonomy] = t
	}

	for _, f := range s.Films {
		for _, key := range finder(&f) {
			// omit any empty keys
			if key != "" {
				_, ok := t[key]
				if !ok {
					t[key] = make(GenericItems, 0)
				}
				t[key] = append(t[key], f.GetGenericItem())
			}
		}
	}
}

// PopulateTaxonomyWithTVSeasons - select the taxonomy and the tv season attribute
func (s Site) PopulateTaxonomyWithTVSeasons(taxonomy string, finder func(*TVSeason) []string) {

	t, ok := s.Taxonomies[taxonomy]

	if !ok {
		t = make(Taxonomy)
		s.Taxonomies[taxonomy] = t
	}

	for _, f := range s.TVSeasons {
		for _, key := range finder(f) {
			// omit any empty keys
			if key != "" {
				_, ok := t[key]
				if !ok {
					t[key] = make(GenericItems, 0)
				}
				t[key] = append(t[key], f.GetGenericItem())
			}
		}
	}
}

// Alphabetical - sort the keys
func (t Taxonomy) Alphabetical() OrderedEntries {
	keys := t.getKeys()
	sort.Strings(keys)

	set := make(OrderedEntries, len(keys))

	for ki, k := range keys {

		sort.Sort(t[k])

		set[ki].Key = k
		set[ki].Items = t[k]
	}
	return set
}

// Print - ordered items
func (t OrderedEntries) Print() {
	for _, v := range t {
		fmt.Printf("key: %s\n", v.Key)
		for _, vv := range v.Items {
			fmt.Printf("\t%s\t - %s (%s)\n", vv.Slug, vv.Title, vv.ItemType)
		}
	}
}

func (t Taxonomy) getKeys() []string {
	keys := make([]string, 0, len(t))
	for k := range t {
		keys = append(keys, k)
	}
	return keys
}

// Print - print a Taxonomy
func (t Taxonomy) Print() {
	for k, v := range t {
		fmt.Printf("key: %s\n", k)

		for _, vv := range v {
			fmt.Printf("\t%s\t - %s (%s)\n", vv.Slug, vv.Title, vv.ItemType)
		}
	}
}

func (slice GenericItems) Len() int {
	return len(slice)
}

func (slice GenericItems) Less(i, j int) bool {
	return slice[i].Title < slice[j].Title
}

func (slice GenericItems) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
