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

import (
	"time"
)

// PlanCollection is a list of published plans
type PlanCollection []Plan

// Plan -
type Plan struct {
	ID              int
	Slug            string
	Name            string
	NameSlug        string
	Description     string
	Interval        string
	IntervalCount   int
	TrialPeriodDays int
	Page            *Page
	PlanType        string
	PortraitImage   string
	LandscapeImage  string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	ExpiryDate      time.Time
	Prices          PriceInfo
}

func (plan *Plan) LinkPlanToPage(site *Site, PageID int) {

	for i := range site.Pages {
		page := &site.Pages[i]
		if page.ID == PageID {
			// link to the page if it exists
			plan.Page = page

			// Conversely, keep track of what plans a page is associated with
			page.Plans = append(page.Plans, *plan)
		}
	}
}

func (plan *Plan) HasExpiryDate() bool {
	return !plan.ExpiryDate.IsZero()
}

// GetGenericItem - returns a generic item
func (plan Plan) GetGenericItem() GenericItem {
	return GenericItem{
		Slug:      plan.Slug,
		ItemType:  "plan",
		InnerItem: plan,
	}
}

// FindPlanBySlug - find the plan by the slug
func (plans *PlanCollection) FindPlanBySlug(slug string) (*Plan, error) {
	coll := *plans
	for i := 0; i < len(coll); i++ {
		if coll[i].Slug == slug {
			return &coll[i], nil
		}
	}
	return nil, ErrDataSourceMissing
}