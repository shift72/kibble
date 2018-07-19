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

func TestStringCollectionString(t *testing.T) {
	strings := StringCollection{"first", "second", "third"}
	assert.Equal(t, "first, second, third", strings.String())
}

func TestStringCollectionJoin(t *testing.T) {
	strings := StringCollection{"first", "second", "third"}
	assert.Equal(t, "first|second|third", strings.Join("|"))
}

func TestRuntimeHours(t *testing.T) {
	assert.Equal(t, 1, Runtime(100).Hours())
}

func TestRuntimeHoursLessThanOne(t *testing.T) {
	assert.Equal(t, 0, Runtime(3).Hours())
}

func TestRuntimeHoursZero(t *testing.T) {
	assert.Equal(t, 0, Runtime(0).Hours())
}

func TestRuntimeMinutesLessThanAnHour(t *testing.T) {
	assert.Equal(t, 2, Runtime(2).Minutes())
}

func TestRuntimeMinutesExactlyAnHour(t *testing.T) {
	assert.Equal(t, 0, Runtime(60).Minutes())
}

func TestRuntimeMinutesZero(t *testing.T) {
	assert.Equal(t, 0, Runtime(0).Minutes())
}

func TestRuntimeFormat_ExpectMinutesOnly(t *testing.T) {
	assert.Equal(t, "runtime_minutes_only", Runtime(0).Localise(i18n.IdentityTfunc()))
}

func TestRuntimeFormat_ExpectRuntime(t *testing.T) {
	assert.Equal(t, "runtime", Runtime(60).Localise(i18n.IdentityTfunc()))
}

func TestRuntimeFormat_ExpectTranslation(t *testing.T) {

	i18n.MustLoadTranslationFile("../sample_site/en_US.all.json")

	T, _ := i18n.Tfunc("en-US")

	assert.Equal(t, "1h 0m", Runtime(60).Localise(T), "runtime 60")
	assert.Equal(t, "1h 1m", Runtime(61).Localise(T), "runtime 61")
	assert.Equal(t, "2h 0m", Runtime(120).Localise(T), "runtime 120")

	assert.Equal(t, "0m", Runtime(0).Localise(T), "runtime 0")
}

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
