package datastore

import (
	"testing"

	"github.com/nicksnyder/go-i18n/i18n"
	"github.com/stretchr/testify/assert"

	"kibble/models"
	"kibble/test"

	"github.com/CloudyKit/jet/v6"
)

func TestSitePlans(t *testing.T) {

	i18n.MustLoadTranslationFile("../sample_site/en_US.all.json")
	T, _ := i18n.Tfunc("en-US")

	ctx := models.RenderContext{}
	view := models.CreateTemplateView(nil, T, &ctx, "../sample_site/templates/")

	renderer1 := &test.InMemoryRenderer{View: view}

	site := &models.Site{
		Plans: models.PlanCollection{
			models.Plan{
				ID:          123,
				Slug:        "/film/123",
				Name:        "Gold",
				Description: "Gold Plan",
			},
		},
	}

	data := jet.VarMap{}
	data.Set("site", site)
	renderer1.Render("./plans.jet", "output.txt", data)

	assert.Contains(t, renderer1.Results[0].Output(), "TestName:Gold")
	assert.Contains(t, renderer1.Results[0].Output(), "TestDescription:Gold Plan")
	assert.Contains(t, renderer1.Results[0].Output(), "TestInterval:")

	assert.NotContains(t, renderer1.Results[0].Output(), "IntervalOptionalCheck")
}

func TestSitePlansWithSubscriptionDetails(t *testing.T) {

	i18n.MustLoadTranslationFile("../sample_site/en_US.all.json")
	T, _ := i18n.Tfunc("en-US")

	ctx := models.RenderContext{}
	view := models.CreateTemplateView(nil, T, &ctx, "../sample_site/templates/")
	view.AddGlobal("version", "v1.1.145")

	renderer1 := &test.InMemoryRenderer{View: view}

	interval := "week"
	intervalCount := 4
	trialPeriodDays := 7

	site := &models.Site{
		Plans: models.PlanCollection{
			models.Plan{
				ID:              123,
				Slug:            "/film/123",
				Name:            "Gold",
				Description:     "Gold Plan",
				Interval:        interval,
				IntervalCount:   intervalCount,
				TrialPeriodDays: trialPeriodDays,
			},
		},
	}

	data := jet.VarMap{}
	data.Set("site", site)
	renderer1.Render("./plans.jet", "output.txt", data)

	assert.Contains(t, renderer1.Results[0].Output(), "TestName:Gold")
	assert.Contains(t, renderer1.Results[0].Output(), "TestDescription:Gold Plan")
	assert.Contains(t, renderer1.Results[0].Output(), "TestInterval:week")
	assert.Contains(t, renderer1.Results[0].Output(), "IntervalCount:4")
	assert.Contains(t, renderer1.Results[0].Output(), "TrialPeriodDays:7")
	assert.Contains(t, renderer1.Results[0].Output(), "TrialPeriodDays-i18n:Try your free trial period of 7 days now!")

	assert.Contains(t, renderer1.Results[0].Output(), "IntervalOptionalCheck")
}
