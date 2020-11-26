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

// PageCollection - part of a page
type PageCollection struct {
	ID          int
	Layout      string
	ItemsPerRow int
	ItemLayout  string
	Slug        string
	TitleSlug   string
	DisplayName string
	Description string
	Items       []GenericItem
}

// Page - page structure
type Page struct {
	ID              int
	Slug            string
	Title           string
	TitleSlug       string
	Content         string
	Tagline         string
	Seo             Seo
	Images          ImageSet
	PageCollections []PageCollection
	PageType        string
	URL             string
}

// Pages -
type Pages []Page

// FindPageByID - find the page by id
func (pages Pages) FindPageByID(pageID int) (*Page, bool) {
	for _, p := range pages {
		if p.ID == pageID {
			return &p, true
		}
	}
	return nil, false
}

// FindPageBySlug - find the page by the slug
func (pages Pages) FindPageBySlug(slug string) (*Page, bool) {
	for _, p := range pages {
		if p.Slug == slug {
			return &p, true
		}
	}
	return nil, false
}

// GetGenericItem - returns a generic item
func (page Page) GetGenericItem() GenericItem {
	return GenericItem{
		Slug:      page.Slug,
		Title:     page.Title,
		Images:    page.Images,
		ItemType:  "page",
		InnerItem: page,
	}
}
