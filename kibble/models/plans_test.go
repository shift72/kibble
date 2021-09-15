package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinkPlanToPage(t *testing.T) {

	site := Site{
		Pages: Pages{
			Page{
				ID:        123,
				Slug:      "/film/123",
				TitleSlug: "the-big-lebowski",
				Plans:     nil,
			},
		},
	}

	plan := Plan{
		ID:   99,
		Page: nil,
	}

	plan.LinkPlanToPage(&site, 123)

	page := site.Pages[0]

	assert.Equal(t, 1, len(page.Plans))
	assert.Equal(t, plan, page.Plans[0])
	assert.Equal(t, &page, plan.Page)
}
