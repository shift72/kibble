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
	"fmt"
	"testing"

	"github.com/indiereign/shift72-kibble/kibble/models"
)

func TestLoadAll(t *testing.T) {

	cfg := &models.Config{
		SiteURL: "https://staging-store.shift72.com",
	}

	itemIndex := make(models.ItemIndex)
	site := &models.Site{}

	AppendAllTVShows(cfg, site, itemIndex)

}

func TestLoadTVSeasons(t *testing.T) {

	// if testing.Short() {
	// 	return
	// }

	cfg := &models.Config{
		SiteURL: "https://staging-store.shift72.com",
	}

	itemIndex := make(models.ItemIndex)
	site := &models.Site{}
	slugs := []string{
		"/tv/4/season/2",
		"/tv/41/season/1",
		"/tv/9/season/1",
	}

	err := AppendTVSeasons(cfg, site, slugs, itemIndex)
	if err != nil {
		t.Error(err)
	}

	if len(itemIndex) == 0 {
		t.Error("Expected some values to be loaded")
	}

	fmt.Printf("here shows %d ", len(site.TVShows))

}
