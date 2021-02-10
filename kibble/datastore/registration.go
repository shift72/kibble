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

package datastore

import (
	"kibble/models"
)

// Init -
func Init() {
	models.AddDataSource(&FilmDataSource{})
	models.AddDataSource(&FilmIndexDataSource{})

	models.AddDataSource(&PageDataSource{})
	models.AddDataSource(&PageIndexDataSource{})

	models.AddDataSource(&BundleDataSource{})
	models.AddDataSource(&BundleIndexDataSource{})

	models.AddDataSource(&CollectionDataSource{})
	models.AddDataSource(&CollectionIndexDataSource{})

	models.AddDataSource(&TVShowDataSource{})
	models.AddDataSource(&TVShowIndexDataSource{})

	models.AddDataSource(&TVSeasonDataSource{})
	models.AddDataSource(&TVSeasonIndexDataSource{})

	models.AddDataSource(&TVEpisodeDataSource{})

	models.AddDataSource(&FileSystemDataSource{})
}
