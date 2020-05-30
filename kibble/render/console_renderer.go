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
	"fmt"

	"github.com/CloudyKit/jet"
	"kibble/models"
)

// ConsoleRenderer - designed to render to the console for testing
type ConsoleRenderer struct {
	view        *jet.Set
	showOutput  bool
	showSummary bool
}

// Initialise - start the rendering process
func (c ConsoleRenderer) Initialise() {

}

// Render - render to the console
func (c ConsoleRenderer) Render(route *models.Route, filePath string, data jet.VarMap) (errorCount int) {

	if c.showSummary {
		fmt.Printf("FilePath: %s\n", filePath)
	}

	w := bytes.NewBufferString("")
	w.Write([]byte("--------------------\n"))
	w.Write([]byte(fmt.Sprintf("FilePath: %s\n", filePath)))
	w.Write([]byte("--------------------\n"))

	t, err := c.view.GetTemplate(route.TemplatePath)
	if err != nil {
		w.Write([]byte("Template error\n"))
		w.Write([]byte(err.Error()))
		if c.showSummary {
			fmt.Println(err)
		}
		return
	}

	if err = t.Execute(w, data, nil); err != nil {
		w.Write([]byte("Execute error\n"))
		w.Write([]byte(err.Error()))
		if c.showSummary {
			fmt.Println(err)
		}
	}

	if c.showOutput {
		fmt.Println(w)
	}

	return
}
