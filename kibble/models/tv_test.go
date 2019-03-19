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

func TestTVSeason_ExpectKeyedTranslation(t *testing.T) {

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

	assert.Equal(t, "Breaking Bad Ceasar N 2", item.GetTranslatedTitle(T, "season_custom_title"), "GetTranslatedTitle")
	assert.Equal(t, "Breaking Bad - Season Alt - 2", item.GetTranslatedTitle(T, ""), "GetTranslatedTitle")
}

func TestTVEpisodeGetTitle(t *testing.T) {
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
		Episodes: []TVEpisode{{
			Slug:          "/tv/123/season/2/episode/1",
			Title:         "First Episode",
			EpisodeNumber: 1,
		}},
	}
	tvSeason.Episodes[0].Season = *tvSeason
	item := tvSeason.Episodes[0].GetGenericItem()

	assert.Equal(t, "Breaking Bad - S2E1 - First Episode", item.GetTitle(T))
	assert.Equal(t, "Breaking Bad - S2E1 - First Episode", tvSeason.Episodes[0].GetTitle(T))
}

func TestTVEpsiodeGetTranslatedTitle(t *testing.T) {
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
		Episodes: []TVEpisode{{
			Slug:          "/tv/123/season/2/episode/1",
			Title:         "First Episode",
			EpisodeNumber: 1,
		}},
	}
	tvSeason.Episodes[0].Season = *tvSeason
	item := tvSeason.Episodes[0].GetGenericItem()

	assert.Equal(t, "Breaking Bad - Season 2 Episode 1: First Episode", item.GetTranslatedTitle(T, "tvepisode_custom_title"))
	assert.Equal(t, "Breaking Bad - Season 2 Episode 1: First Episode", tvSeason.Episodes[0].GetTranslatedTitle(T, "tvepisode_custom_title"))
}
