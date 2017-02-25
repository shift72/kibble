package datastore

import "github.com/indiereign/shift72-kibble/kibble/models"

// dummy data
var allFilms = []models.Film{
	{
		ID:    1,
		Title: "Forrest Gump",
		Synopsis: ` syn
# one
# two
`,
	},
	{
		ID:    2,
		Title: "Angel at my table",
		Synopsis: `## header
    `,
	},
}
