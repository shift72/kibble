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
	"testing"

	"kibble/models"
	"kibble/utils"

	"github.com/stretchr/testify/assert"
)

func TestPlansMapping(t *testing.T) {

	itemIndex := make(models.ItemIndex)

	serviceConfig := commonServiceConfig()

	apiPlan := PlansV1{
		Name:        "Bronze Plan",
		Description: "Plan description",
		Status:      "active",
	}

	model := apiPlan.mapToModel(serviceConfig, itemIndex)

	assert.Equal(t, "Bronze Plan", model.Name)
	assert.Equal(t, "bronze-plan", model.NameSlug)
	assert.Equal(t, "Plan description", model.Description)
	assert.Equal(t, "", model.Interval)
}

func TestRecurringPlanWithExpiryDate(t *testing.T) {

	itemIndex := make(models.ItemIndex)

	serviceConfig := commonServiceConfig()

	interval := "week"
	intervalCount := 4
	trialPeriodDays := 7
	planType := "recurring"
	expiryDate := utils.ParseTimeFromString("2021-04-01T03:02:17.000Z")

	apiPlan := PlansV1{
		Name:            "Bronze Plan",
		Description:     "Plan description",
		Status:          "active",
		Interval:        &interval,
		IntervalCount:   &intervalCount,
		TrialPeriodDays: &trialPeriodDays,
		PlanType:        &planType,
		ExpiryDate:      expiryDate,
	}

	model := apiPlan.mapToModel(serviceConfig, itemIndex)

	assert.Equal(t, "Bronze Plan", model.Name)
	assert.Equal(t, "bronze-plan", model.NameSlug)
	assert.Equal(t, "Plan description", model.Description)
	assert.Equal(t, "week", model.Interval)
	assert.Equal(t, 4, model.IntervalCount)
	assert.Equal(t, 7, model.TrialPeriodDays)
	assert.Equal(t, "recurring", model.PlanType)
	assert.True(t, model.HasExpiryDate())
}

func TestOneOffPlanWithNoExpiryDate(t *testing.T) {
	itemIndex := make(models.ItemIndex)

	serviceConfig := commonServiceConfig()

	planType := "one_off"

	apiPlan := PlansV1{
		Name:        "One Off Plan",
		Description: "Plan description 123",
		Status:      "active",
		PlanType:    &planType,
	}

	model := apiPlan.mapToModel(serviceConfig, itemIndex)
	assert.Equal(t, "One Off Plan", model.Name)
	assert.Equal(t, "one-off-plan", model.NameSlug)
	assert.Equal(t, "Plan description 123", model.Description)
	assert.Equal(t, "", model.Interval)
	assert.Equal(t, 0, model.IntervalCount)
	assert.Equal(t, 0, model.TrialPeriodDays)
	assert.Equal(t, "one_off", model.PlanType)
	assert.False(t, model.HasExpiryDate())
}
