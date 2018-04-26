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
