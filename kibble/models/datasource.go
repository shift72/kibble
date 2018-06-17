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

import "errors"

var (
	// ErrDataSource is returned when no datasource could be found
	ErrDataSource = "!Error"
	// ErrDataSourceMissing raises when the datasource could not be found
	ErrDataSourceMissing = errors.New("Missing")
)

var store map[string]DataSource

// AddDataSource - register a datasource
func AddDataSource(ds DataSource) {

	if store == nil {
		store = make(map[string]DataSource)
	}

	store[ds.GetName()] = ds
}

// FindDataSource - find the data source by name
func FindDataSource(name string) DataSource {
	return store[name]
}
