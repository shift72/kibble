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

	path = site.SiteBrand.GetLink("non-existant-link", "/defaultPath")
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

func TestGetImageWhenSiteBrandImagesContainsMultipleImages(t *testing.T) {

	site := Site{
		SiteBrand: SiteBrand{
			Images: map[string]string{
				"logo@1x":  "{cloudfront}/path/logo@1x.png",
				"logo@2x":  "{cloudfront}/path/logo@2x.png",
				"app-logo": "{cloudfront}/path/app-logo.gif",
			},
		},
	}

	path := site.SiteBrand.GetImage("logo@1x", "/defaultPath")
	assert.Equal(t, string(path), "{cloudfront}/path/logo@1x.png")

	path = site.SiteBrand.GetImage("logo@2x", "/defaultPath")
	assert.Equal(t, string(path), "{cloudfront}/path/logo@2x.png")

	path = site.SiteBrand.GetImage("app-logo", "/defaultPath")
	assert.Equal(t, string(path), "{cloudfront}/path/app-logo.gif")

	path = site.SiteBrand.GetImage("favicon", "/defaultPath")
	assert.Equal(t, string(path), "/defaultPath")

}

func TestGetImageWhenSiteBrandImagesContainsAllImagesTypes(t *testing.T) {

	site := Site{
		SiteBrand: SiteBrand{
			Images: map[string]string{
				"logo@1x":             "{cloudfront}/path/logo@1x.png",
				"logo@2x":             "{cloudfront}/path/logo@2x.png",
				"favicon":             "{cloudfront}/path/favicon.ico",
				"icons-192":           "{cloudfront}/path/icons-192.png",
				"icons-512":           "{cloudfront}/path/icons-512.png",
				"facebook-image":      "{cloudfront}/path/facebook-image.png",
				"app-logo":            "{cloudfront}/path/app-logo.gif",
				"email-logo":          "{cloudfront}/path/email-logo.png",
				"sponsor-banner-xxxs": "{cloudfront}/path/sponsor-banner-xxxs.png",
				"sponsor-banner-xs":   "{cloudfront}/path/sponsor-banner-xs.png",
				"sponsor-banner-md":   "{cloudfront}/path/sponsor-banner-md.png",
			},
		},
	}

	path := site.SiteBrand.GetImage("logo@1x", "/defaultPath")
	assert.Equal(t, string(path), "{cloudfront}/path/logo@1x.png")

	path = site.SiteBrand.GetImage("logo@2x", "/defaultPath")
	assert.Equal(t, string(path), "{cloudfront}/path/logo@2x.png")

	path = site.SiteBrand.GetImage("favicon", "/defaultPath")
	assert.Equal(t, string(path), "{cloudfront}/path/favicon.ico")

	path = site.SiteBrand.GetImage("icons-192", "/defaultPath")
	assert.Equal(t, string(path), "{cloudfront}/path/icons-192.png")

	path = site.SiteBrand.GetImage("icons-512", "/defaultPath")
	assert.Equal(t, string(path), "{cloudfront}/path/icons-512.png")

	path = site.SiteBrand.GetImage("facebook-image", "/defaultPath")
	assert.Equal(t, string(path), "{cloudfront}/path/facebook-image.png")

	path = site.SiteBrand.GetImage("app-logo", "/defaultPath")
	assert.Equal(t, string(path), "{cloudfront}/path/app-logo.gif")

	path = site.SiteBrand.GetImage("email-logo", "/defaultPath")
	assert.Equal(t, string(path), "{cloudfront}/path/email-logo.png")

	path = site.SiteBrand.GetImage("sponsor-banner-xxxs", "/defaultPath")
	assert.Equal(t, string(path), "{cloudfront}/path/sponsor-banner-xxxs.png")

	path = site.SiteBrand.GetImage("sponsor-banner-xs", "/defaultPath")
	assert.Equal(t, string(path), "{cloudfront}/path/sponsor-banner-xs.png")

	path = site.SiteBrand.GetImage("sponsor-banner-md", "/defaultPath")
	assert.Equal(t, string(path), "{cloudfront}/path/sponsor-banner-md.png")

	path = site.SiteBrand.GetImage("non-existant-image", "/defaultPath")
	assert.Equal(t, string(path), "/defaultPath")

}
