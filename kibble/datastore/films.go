package datastore

import (
	"github.com/indiereign/shift72-kibble/kibble/models"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

func applyContentTransforms(data string) string {

	//TODO: apply shortcodes

	// apply mark down
	unsafe := blackfriday.MarkdownCommon([]byte(data))

	// return string(unsafe)
	// apply sanitization
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)

	return string(html)
}

// Init -
func Init() {
	//TODO: needs to be fixed
	models.AddDataSource(&FilmDataSource{})
	models.AddDataSource(&FilmCollectionDataSource{})

	models.AddDataSource(&PageDataSource{})
	models.AddDataSource(&PageCollectionDataSource{})
}

// dummy data
var allFilms = []models.Film{
	{
		ID:    1,
		Slug:  "/film/1",
		Title: "Forrest Gump",
		Synopsis: ` syn
# one
# two
`,
	},
	{
		ID:    2,
		Slug:  "/film/2",
		Title: "Angel at my table",
		Synopsis: `## header
    `,
	},
}
