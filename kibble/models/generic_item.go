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
	ItemType string
	Slug     string
	Title    string
	Images   ImageSet
	Seo      Seo
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
	}
	return i.Title
}
