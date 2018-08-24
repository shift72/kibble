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

import "github.com/nicksnyder/go-i18n/i18n"

// Runtime - Allows us to get accurate measures of hours and minutes.
type Runtime int

// Hours - returns the basic number of hours in the runtime
func (rt Runtime) Hours() int {
	return int(rt) / 60
}

// Minutes - returns the number of minutes (left =over from any hours) in the runtime.
// This is not the total minutes.
func (rt Runtime) Minutes() int {
	return int(rt) % 60
}

// Localise - returns the localised runtime format
func (rt Runtime) Localise(T i18n.TranslateFunc) string {
	arg := map[string]interface{}{
		"Hours":   rt.Hours(),
		"Minutes": rt.Minutes(),
	}

	h := rt.Hours()
	if h == 0 {
		// zero != 0 in languages
		// https://github.com/nicksnyder/go-i18n/issues/58
		return T("runtime_minutes_only", arg)
	}

	return T("runtime", h, arg)
}
