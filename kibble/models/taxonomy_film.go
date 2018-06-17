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

// GetYear - extract the year from the film
func GetYear(f *Film) []string {
	return []string{strconv.Itoa(f.ReleaseDate.Year())}
}

// GetCountries - extract the countries from the film
func GetCountries(f *Film) []string {
	return f.Countries
}

// GetGenres - extract the genres from the film
func GetGenres(f *Film) []string {
	return f.Genres
}

// GetCast - extact the cast members from the film
func GetCast(f *Film) []string {

	cast := make([]string, len(f.Cast))

	for _, i := range f.Cast {
		cast = append(cast, i.Name)
	}

	return cast
}
