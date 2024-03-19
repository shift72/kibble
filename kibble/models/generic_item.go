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

import "github.com/nicksnyder/go-i18n/i18n"

// GenericItem - used to store the common properties
type GenericItem struct {
	// link to the actual item
	InnerItem interface{}
	// film / show / season/ episode / bundle / page
	ItemType      string
	Slug          string
	Title         string
	Images        ImageSet
	Seo           Seo
	CarouselFocus string
}

// GetTitle - returns the title in the current language
// expect to be called as item.GetTitle(i18n) where i18n is the translation function
// for the current language
func (i *GenericItem) GetTitle(T i18n.TranslateFunc) string {
	switch i.ItemType {
	case "tvseason":
		if s, ok := i.InnerItem.(TVSeason); ok {
			return s.GetTitle(T)
		}
	case "tvepisode":
		if e, ok := i.InnerItem.(TVEpisode); ok {
			return e.GetTitle(T)
		}
	}
	return i.Title
}

// GetTranslatedTitle returns an i18n version of a GenericItem title using the specified key as the template
func (i GenericItem) GetTranslatedTitle(T i18n.TranslateFunc, i18nKey string) string {
	switch i.ItemType {
	case "tvseason":
		if s, ok := i.InnerItem.(TVSeason); ok {
			return s.GetTranslatedTitle(T, i18nKey)
		}
	case "tvepisode":
		if e, ok := i.InnerItem.(TVEpisode); ok {
			return e.GetTranslatedTitle(T, i18nKey)
		}
	}
	return i.Title
}

// get carousel focus
func (i GenericItem) GetCarouselImageFocusArea() string {
	switch i.ItemType {
	case "film":
		if f, ok := i.InnerItem.(Film); ok {
			return f.CarouselFocus
		}
	case "tvseason":
		if s, ok := i.InnerItem.(TVSeason); ok {
			return s.CarouselFocus
		}
	case "page":
		if p, ok := i.InnerItem.(Page); ok {
			return p.CustomFields.GetString("carousel_focus", "")
		}
	case "bundle":
		if b, ok := i.InnerItem.(Bundle); ok {
			return b.CustomFields.GetString("carousel_focus", "")
		}
	}
	return i.CarouselFocus
}
