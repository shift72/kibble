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

// CollectionCollection - all collections
type CollectionCollection []Collection

// Collection - collection of films and tv seasons / episodes
type Collection struct {
	ID          int
	Slug        string
	Title       string
	TitleSlug   string
	Description string
	DisplayName string
	ItemLayout  string
	ItemsPerRow int
	Images      ImageSet
	Seo         Seo
	Items       []GenericItem
	SearchQuery string
	CreatedAt   string
	UpdatedAt   string
}

// FindCollectionByID - find the page by id
func (collections *CollectionCollection) FindCollectionByID(collectionID int) (*Collection, error) {
	coll := *collections
	for i := 0; i < len(coll); i++ {
		if coll[i].ID == collectionID {
			return &coll[i], nil
		}
	}
	return nil, ErrDataSourceMissing
}

// FindCollectionBySlug - find the collection by the slug
func (collections *CollectionCollection) FindCollectionBySlug(slug string) (*Collection, error) {
	coll := *collections
	for i := 0; i < len(coll); i++ {
		if coll[i].Slug == slug || coll[i].TitleSlug == slug {
			return &coll[i], nil
		}
	}
	return nil, ErrDataSourceMissing
}

// GetGenericItem - returns a generic item
func (collection Collection) GetGenericItem() GenericItem {
	return GenericItem{
		Title:     collection.Title,
		Seo:       collection.Seo,
		Images:    collection.Images,
		InnerItem: collection,
	}
}
