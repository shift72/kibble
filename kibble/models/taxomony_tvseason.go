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

import "strconv"

// GetTVSeasonYear - extract the year from the tv season
func GetTVSeasonYear(t *TVSeason) []string {
	return []string{strconv.Itoa(t.ShowInfo.ReleaseDate.Year())}
}

// GetTVShowGenres - extract the genres from the tv show
func GetTVShowGenres(t *TVSeason) []string {
	return t.ShowInfo.Genres
}

// GetTVShowCountries - extract the countries from the film
func GetTVShowCountries(t *TVSeason) []string {
	return t.ShowInfo.Countries
}

// GetTVShowCast - extact the cast members from the tv show
func GetTVShowCast(t *TVSeason) []string {

	cast := make([]string, len(t.Cast))

	for _, i := range t.Cast {
		cast = append(cast, i.Name)
	}

	return cast
}
