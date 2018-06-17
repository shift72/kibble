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

import "time"

// BundleCollection - all bundles
type BundleCollection []Bundle

// Bundle - model
type Bundle struct {
	ID            int
	Slug          string
	Title         string
	TitleSlug     string
	Tagline       string
	Description   string
	Status        string
	Seo           Seo
	Images        ImageSet
	PromoURL      string
	ExternalID    string
	Items         GenericItems
	PublishedDate time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// FindBundleByID - find the page by id
func (bundles BundleCollection) FindBundleByID(bundleID int) (*Bundle, error) {
	for _, b := range bundles {
		if b.ID == bundleID {
			return &b, nil
		}
	}
	return nil, ErrDataSourceMissing
}

// FindBundleBySlug - find the bundle by the slug
func (bundles BundleCollection) FindBundleBySlug(slug string) (*Bundle, error) {
	for _, p := range bundles {
		if p.Slug == slug || p.TitleSlug == slug {
			return &p, nil
		}
	}
	return nil, ErrDataSourceMissing
}

// GetGenericItem - returns a generic item
func (bundle Bundle) GetGenericItem() GenericItem {
	return GenericItem{
		Title:     bundle.Title,
		Seo:       bundle.Seo,
		Images:    bundle.Images,
		InnerItem: bundle,
	}
}
