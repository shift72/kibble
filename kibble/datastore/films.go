package datastore

import (
	"github.com/indiereign/shift72-kibble/kibble/models"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

// https://www.ozflix.tv/services/meta/v1/bundles

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
