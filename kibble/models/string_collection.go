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

import "kibble/utils"

// StringCollection - Allows us to add methods to []string for easing UI array usage
type StringCollection []string

// String - overrides String function for []string which joins all items in the array into english readable format.
func (strings StringCollection) String() string {
	return strings.Join(", ")
}

// Join - Joins all items in []string using specified separator.
func (strings StringCollection) Join(separator string) string {
	return utils.Join(separator, strings...)
}
