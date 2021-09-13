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

package exit

const (
	ok                       int = 0
	one                          = 1
	minusOne                     = -1
	minusTwo                     = -2
	FailedToLoadTranslations     = -3
)

// If we can be sure that changing the -1/-2 doesn't affect the surface API of exit values,
// then using iota seems to be idiomatic. It starts from zero and auto-increments.
// const (
// 	ok = iota
// 	minusOne
// 	minusTwo
// 	failedToLoadTranslations
// )
