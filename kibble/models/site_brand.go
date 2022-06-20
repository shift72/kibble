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

type SiteBrand struct {
	Images map[string]string
	Links  map[string]string
}

func (s *SiteBrand) GetImage(assetType string, defaultURL string) string {

	return getItem(s.Images, assetType, defaultURL)
}

func (s *SiteBrand) GetLink(assetType string, defaultURL string) string {

	return getItem(s.Links, assetType, defaultURL)
}

//Get URL for SiteBrand item, return a passed in default URL if non-existant
func getItem(assetMap map[string]string, assetType string, defaultURL string) string {

	if url, found := assetMap[assetType]; found {
		return url
	}
	return defaultURL
}
