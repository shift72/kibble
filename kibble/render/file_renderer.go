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

package render

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"kibble/utils"

	"github.com/CloudyKit/jet"
)

// FileRenderer - designed to render to the file system for testing
type FileRenderer struct {
	view       *jet.Set
	buildPath  string
	sourcePath string
}

// Initialise - start the rendering process
func (c FileRenderer) Initialise() {
	os.RemoveAll(c.buildPath)

	err := utils.CopyDir(filepath.Join(c.sourcePath, staticFolder), c.buildPath)
	if err != nil {
		log.Warningf("Warn: static folder copy failed %s", err)
	}

	// copy language files too, they are a special file name format
	glob := filepath.Join(c.sourcePath, "/*.json")
	langFiles, err := filepath.Glob(glob)
	if len(langFiles) > 0 {
		for _, file := range langFiles {
			dst := filepath.Join(c.buildPath, filepath.Base(file))
			err := utils.CopyFile(file, dst)
			if err != nil {
				log.Warningf("Warn: language file (%s) copy failed %s", file, err)
			}
		}
	}
}

// Render - render to a file
func (c FileRenderer) Render(templatePath string, filePath string, data jet.VarMap) int {
	errorCount := 0
	defer func() {
		if r := recover(); r != nil {
			errorCount++
			log.Errorf("Error. Recovered from %s in %s", r, templatePath)
		}
	}()

	fullPath := path.Join(c.buildPath, filePath)
	if strings.HasSuffix(filePath, "/") {
		fullPath = path.Join(fullPath, "index.html")
	}

	log.Debugf("FilePath: %s", fullPath)

	w := bytes.NewBufferString("")
	t, err := c.view.GetTemplate(templatePath)
	if err != nil {
		errorCount++
		log.Errorf("Template load error: %s in %s", err, templatePath)
		return errorCount
	}

	data.Set("currentUrlPath", filePath)
	if err = t.Execute(w, data, nil); err != nil {
		errorCount++
		w.WriteString("<pre>")
		w.WriteString(err.Error())
		w.WriteString("</pre>")
		log.Errorf("Template execute error: %s in %s", err, templatePath)
		return errorCount
	}

	dirPath := filepath.Dir(fullPath)
	os.MkdirAll(dirPath, 0777)

	// optional check
	if _, err := os.Stat(fullPath); err == nil {
		log.Warningf("File exists and will be overwritten: %s", fullPath)
	}

	err = ioutil.WriteFile(fullPath, w.Bytes(), 0777)
	if err != nil {
		errorCount++
		log.Errorf("File write: %s attempting %s", err, fullPath)
	}

	return errorCount
}
