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

package utils

import (
	"path/filepath"
	"strings"
)

// FileIgnorer allows specifying groups of files to be ignored based on
// Glob or prefix file path matching.
type FileIgnorer struct {
	rootPath string
	patterns []string
}

var defaultIgnorePaths = []string{
	".git",
	".kibble",
	"node_modules",
	"npm-debug.log",
	"package.json",
	"package-lock.json",
	".DS_Store",
	"README.md",
	".gitignore",
}

// NewFileIgnorer creates a new FileIgnorer and sets the correct paths
func NewFileIgnorer(root string, patterns []string) FileIgnorer {
	fi := FileIgnorer{
		rootPath: root,
		patterns: append(defaultIgnorePaths, patterns...),
	}

	for i, p := range fi.patterns {
		fi.patterns[i] = filepath.Join(fi.rootPath, p)
	}

	return fi
}

// IsIgnored returns whether the specified path should be ignored based
// on the matching the current patterns.
func (fm FileIgnorer) IsIgnored(path string) bool {
	for _, pattern := range fm.patterns {
		isMatch, err := filepath.Match(pattern, path)
		if err != nil {
			log.Errorf("Failed to match %s to %s. %s", path, pattern, err.Error())
		}
		// support both file globs and simple dir names (which the `filepath.Match` command seems to not support).
		if isMatch || strings.HasPrefix(path, pattern) {
			return true
		}
	}

	return false
}
