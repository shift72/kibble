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
)

func TestLoadConfig(t *testing.T) {

	if testing.Short() {
		return
	}

	cfg := &models.Config{
		SiteURL: "https://staging-store.shift72.com",
	}

	serviceConfig, err := LoadConfig(cfg)
	if err != nil {
		t.Error(err)
	}

	if len(serviceConfig) == 0 {
		t.Error("Expected some values to be loaded")
	}
}

func TestLoadFeatureToggles(t *testing.T) {

	if testing.Short() {
		return
	}

	cfg := &models.Config{
		SiteURL: "https://staging-store.shift72.com",
	}

	toggles, err := LoadFeatureToggles(cfg)
	if err != nil {
		t.Error(err)
	}

	if len(toggles) == 0 {
		t.Error("Expected some values to be loaded")
	}
}
