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
	"os"
	"path/filepath"
)

// SourcePath is an absolute path to where the Site source files are located.
// Can be configured in `site.json` - `"siteRootPath"`
func (cfg Config) SourcePath() string {
	wd, _ := os.Getwd()

	if cfg.SiteRootPath == "" || cfg.SiteRootPath == "." {
		return wd
	}

	src := filepath.Join(wd, cfg.SiteRootPath)

	// make sure its a directory
	fi, err := os.Stat(src)
	if err != nil {
		log.Fatalf("'%s' is not a directory.", src)
	}

	switch mode := fi.Mode(); {
	case mode.IsDir():
		return src
	default:
		log.Fatalf("'%s' is not a directory", src)
	}

	// We shouldnt get here, maybe this method should also return an error?
	return wd
}

// ShortCodePath is the path where the short code templates will be stored
func (cfg Config) ShortCodePath() string {
	return filepath.Join(cfg.SourcePath(), "templates/shortcodes")
}
