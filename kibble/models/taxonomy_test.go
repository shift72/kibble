package models

import (
	"testing"
)

func create2FilmSite() Site {
	return Site{
		Films: FilmCollection{
			Film{
				ID:        123,
				Slug:      "/film/123",
				TitleSlug: "the-big-lebowski",
				Genres:    []string{"comedy", "drama"},
			},
			Film{
				ID:        124,
				Slug:      "/film/124",
				TitleSlug: "fargo",
				Genres:    []string{"comedy", "mystery"},
			},
		},
		Taxonomies: make(Taxonomies),
	}
}

func TestBuildGenreTaxonomy(t *testing.T) {

	s := Site{
		Films: FilmCollection{
			Film{
				ID:        123,
				Slug:      "/film/123",
				TitleSlug: "the-big-lebowski",
				Genres:    []string{"comedy", "drama"},
			},
		},
		Taxonomies: make(Taxonomies),
	}

	// act
	s.PopulateTaxonomyWithFilms("genres", GetGenres)

	// expect
	if len(s.Taxonomies["genres"]) != 2 {
		t.Errorf("expected 2 genres found %d", len(s.Taxonomies["genres"]))
	}

	if len(s.Taxonomies["genres"]["comedy"]) != 1 {
		t.Errorf("expected 1 film found %d", len(s.Taxonomies["genres"]))
	}
}

func TestGiven2FilmsBuildGenreTaxonomy(t *testing.T) {

	s := create2FilmSite()

	// act
	s.PopulateTaxonomyWithFilms("genres", GetGenres)

	// expect
	if len(s.Taxonomies["genres"]) != 3 {
		t.Errorf("expected 3 genres found %d", len(s.Taxonomies["genres"]))
	}

	if len(s.Taxonomies["genres"]["comedy"]) != 2 {
		t.Errorf("expected 2 films found %d", len(s.Taxonomies["genres"]))
	}
}

func TestGiven2FilmsSortGenreTaxonomy(t *testing.T) {

	s := create2FilmSite()

	// act
	s.PopulateTaxonomyWithFilms("genres", GetGenres)

	orderedGenres := s.Taxonomies["genres"].Alphabetical()

	// expect
	if orderedGenres[0].Key != "comedy" {
		t.Errorf("expected 'comedy' genres found %s", orderedGenres[0].Key)
	}

	if orderedGenres[1].Key != "drama" {
		t.Errorf("expected 'comedy' genres found %s", orderedGenres[1].Key)
	}

	if orderedGenres[2].Key != "mystery" {
		t.Errorf("expected 'comedy' genres found %s", orderedGenres[2].Key)
	}
}
