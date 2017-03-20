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
