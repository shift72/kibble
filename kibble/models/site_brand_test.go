//    Copyright 2018 SHIFT72
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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetImageWhenSiteBrandImagesIsEmpty(t *testing.T) {

	var site Site

	path := site.SiteBrand.GetImage("logo@1x", "/defaultPath")
	assert.Equal(t, string(path), "/defaultPath")

}

func TestGetLinkWhenSiteBrandLinksIsEmpty(t *testing.T) {

	var site Site

	path := site.SiteBrand.GetLink("css", "/defaultPath")
	assert.Equal(t, string(path), "/defaultPath")

}

func TestGetLinkWhenSiteBrandLinksContainsCSS(t *testing.T) {

	site := Site{
		SiteBrand: SiteBrand{
			Links: map[string]string{
				"css": "{cloudfront}/path/{unix_timestamp}.css",
			},
		},
	}

	path := site.SiteBrand.GetLink("css", "/defaultPath")
	assert.Equal(t, string(path), "{cloudfront}/path/{unix_timestamp}.css")

	path = site.SiteBrand.GetLink("non-existent-link", "/defaultPath")
	assert.Equal(t, string(path), "/defaultPath")

}

func TestGetImageWhenSiteBrandImagesContainsSingleImage(t *testing.T) {

	site := Site{
		SiteBrand: SiteBrand{
			Images: map[string]string{
				"logo@1x": "{cloudfront}/path/logo@1x.png",
			},
		},
	}

	path := site.SiteBrand.GetImage("logo@1x", "/defaultPath")
	assert.Equal(t, string(path), "{cloudfront}/path/logo@1x.png")

	path = site.SiteBrand.GetImage("logo@2x", "/defaultPath")
	assert.Equal(t, string(path), "/defaultPath")

}
