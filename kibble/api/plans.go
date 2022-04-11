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

package api

import (
	"encoding/json"
	"fmt"
	"time"

	"kibble/models"

	"github.com/gosimple/slug"
)

// LoadAllPlans - loads all active plans
func LoadAllPlans(cfg *models.Config, site *models.Site, itemIndex *models.ItemIndex) error {

	path := fmt.Sprintf("%s/services/pricing/v1/plans", cfg.SiteURL)
	data, err := Get(cfg, path)
	if err != nil {
		return err
	}

	var apiPlans []PlansV1
	err = json.Unmarshal([]byte(data), &apiPlans)
	if err != nil {
		return err
	}

	for _, b := range apiPlans {
		plan := b.mapToModel(site.Config, itemIndex)

		plan.LinkPlanToPage(site, b.PageID)

		site.Plans = append(site.Plans, plan)

		itemIndex.Set(plan.Slug, plan.GetGenericItem())
	}

	return nil
}

func (p PlansV1) mapToModel(serviceConfig models.ServiceConfig, itemIndex *models.ItemIndex) models.Plan {
	m := models.Plan{
		ID:              p.ID,
		Slug:            fmt.Sprintf("/plan/%d", p.ID),
		NameSlug:        slug.Make(p.Name),
		Name:            p.Name,
		Interval:        "",
		IntervalCount:   0,
		TrialPeriodDays: 0,
		PlanType:        "",
		ExpiryDate:      p.ExpiryDate,
		PortraitImage:   serviceConfig.ForceAbsoluteImagePath(p.PortraitImage),
		LandscapeImage:  serviceConfig.ForceAbsoluteImagePath(p.LandscapeImage),
		Description:     p.Description,
		CreatedAt:       p.CreatedAt,
		UpdatedAt:       p.UpdatedAt,
	}

	if p.Interval != nil {
		m.Interval = *p.Interval
	}
	if p.IntervalCount != nil {
		m.IntervalCount = *p.IntervalCount
	}
	if p.TrialPeriodDays != nil {
		m.TrialPeriodDays = *p.TrialPeriodDays
	}
	if p.PlanType != nil {
		m.PlanType = *p.PlanType
	}
	return m
}

// PlansV1 - model
type PlansV1 struct {
	ID              int       `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	ExpiryDate      time.Time `json:"expiry_date"`
	PageID          int       `json:"page_id"`
	Interval        *string   `json:"interval"`
	IntervalCount   *int      `json:"interval_count"`
	PlanType        *string   `json:"plan_type"`
	TrialPeriodDays *int      `json:"trial_period_days"`
	PortraitImage   string    `json:"portrait_image"`
	LandscapeImage  string    `json:"landscape_image"`
}
