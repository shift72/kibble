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

	"github.com/nicksnyder/go-i18n/i18n"
	"github.com/stretchr/testify/assert"
)

func TestTVSeason_ExpectTranslation(t *testing.T) {

	i18n.MustLoadTranslationFile("../sample_site/en_US.all.json")

	T, _ := i18n.Tfunc("en-US")

	tvSeason := &TVSeason{
		SeasonNumber: 2,
		ShowInfo: &TVShow{
			ID:        123,
			Title:     "Breaking Bad",
			TitleSlug: "breaking-bad",
		},
		Slug: "/tv/123/season/2",
	}

	item := tvSeason.GetGenericItem()

	assert.Equal(t, "Breaking Bad - Season Alt - 2", item.GetTitle(T), "GetTitle")

}