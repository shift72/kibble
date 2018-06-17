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
	"testing"

	"github.com/stretchr/testify/assert"
)

var cfg = &Config{
	SiteRootPath: "..",
}

func TestEmptyContent(t *testing.T) {

	ConfigureShortcodeTemplatePath(cfg)
	assert.Equal(t, "", ApplyContentTransforms(""))
}
func TestMarkdownContent(t *testing.T) {
	ConfigureShortcodeTemplatePath(cfg)
	assert.Equal(t, "<h1>ONE</h1>\n", ApplyContentTransforms("# ONE"))
}
func TestEchoEmbeddedTemplateContent(t *testing.T) {
	ConfigureShortcodeTemplatePath(cfg)
	assert.Equal(t,
		"<p>aa<div class=\"echo\">slug:/film/1</div>bb</p>\n",
		ApplyContentTransforms("aa{{echo slug=/film/1 }}bb"))
}
func TestS72VideoTemplateContent(t *testing.T) {
	ConfigureShortcodeTemplatePath(cfg)
	assert.Equal(t,
		"<p>aa<video slug=\"/film/1\"></video>bb</p>\n",
		ApplyContentTransforms("aa{{video slug=/film/1 }}bb"))
}
func TestS72VideoTemplateContentMultiple(t *testing.T) {
	ConfigureShortcodeTemplatePath(cfg)
	assert.Equal(t,
		"<p>aa<div class=\"echo\">slug:/film/1</div>bb<div class=\"echo\">slug:/film/2</div>cc</p>\n",
		ApplyContentTransforms("aa{{echo slug=/film/1}}bb{{ echo slug=/film/2 }}cc"))
}

func TestYoutubeTemplateDefault(t *testing.T) {
	ConfigureShortcodeTemplatePath(cfg)
	assert.Equal(t,
		"<p>\n<div style=\"position: relative; padding-bottom: 56.25%; padding-top: 30px; height: 0; overflow: hidden;\" >\n<iframe src=\"//www.youtube.com/embed/aaaa\" style=\"position: absolute; top: 0; left: 0; width: 100%; height: 100%;\" allowfullscreen frameborder=\"0\"></iframe>\n</div></p>\n",
		ApplyContentTransforms("{{youtube id=aaaa}}"))
}

func TestYoutubeTemplateWithClass(t *testing.T) {
	ConfigureShortcodeTemplatePath(cfg)
	assert.Equal(t,
		"<p>\n<div class=\"yt\" >\n<iframe src=\"//www.youtube.com/embed/aaaa\" class=\"yt\" allowfullscreen frameborder=\"0\"></iframe>\n</div></p>\n",
		ApplyContentTransforms("{{youtube id=aaaa class=yt}}"))
}

func TestYoutubeTemplateWithClassAutoplay(t *testing.T) {
	ConfigureShortcodeTemplatePath(cfg)
	assert.Equal(t,
		"<p>\n<div class=\"yt\" >\n<iframe src=\"//www.youtube.com/embed/aaaa\" class=\"yt\" autoplay=1 allowfullscreen frameborder=\"0\"></iframe>\n</div></p>\n",
		ApplyContentTransforms("{{youtube id=aaaa class=yt autoplay=true}}"))
}

func TestYoutubeTemplateWithClassAutoplayOff(t *testing.T) {
	ConfigureShortcodeTemplatePath(cfg)
	assert.Equal(t,
		"<p>\n<div class=\"yt\" >\n<iframe src=\"//www.youtube.com/embed/aaaa\" class=\"yt\" allowfullscreen frameborder=\"0\"></iframe>\n</div></p>\n",
		ApplyContentTransforms("{{youtube id=aaaa class=yt autoplay=false}}"))
}

func TestEvilContent(t *testing.T) {
	ConfigureShortcodeTemplatePath(cfg)
	assert.Equal(t,
		"<p>JS attempt:</p>\n",
		ApplyContentTransforms("JS attempt:<script src=\"https://blah.com/evil.js\" ></script>"))
}
