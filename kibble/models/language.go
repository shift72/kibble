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

// LanguageCollection - all languages
type LanguageCollection []Language

// Language - instance of a language
type Language struct {
	Code  string `json:"code"`
	Name  string `json:"name"`
	Label string `json:"label"`
	// Deprecated: Use Language.Code instead.
	Locale             string `json:"locale"`
	DefinitionFilePath string `json:"-"`
	IsDefault          bool   `json:"-"`
}
