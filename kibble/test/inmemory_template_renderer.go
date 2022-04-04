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

package test

import (
	"bytes"

	"github.com/CloudyKit/jet/v6"
)

// InMemoryTemplateRenderer - render template in memory
type InMemoryTemplateRenderer struct {
	View     *jet.Set
	Template *jet.Template
	Result   InMemoryResult
}

// Initialise -
func (c *InMemoryTemplateRenderer) Initialise() {

}

// Render - render the pages to memory
func (c *InMemoryTemplateRenderer) Render(templatePath string, filePath string, data jet.VarMap) (errCount int) {

	c.Result = InMemoryResult{
		buffer:   bytes.NewBufferString(""),
		filePath: filePath,
	}

	if err := c.Template.Execute(c.Result.buffer, data, nil); err != nil {
		errCount++
		c.Result.err = err
	}

	return
}
