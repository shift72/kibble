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
	"github.com/CloudyKit/jet/v6"
)

// MockRenderer - testable renderer
type MockRenderer struct {
	InitialisedCalled bool
	RenderCalled      bool
	TemplatePath      string
	FilePath          string
	Data              jet.VarMap
}

// Initialise - start the rendering process
func (c *MockRenderer) Initialise() {
	c.InitialisedCalled = true
}

// Render - render to the console
func (c *MockRenderer) Render(templatePath string, filePath string, data jet.VarMap) (errCount int) {
	c.RenderCalled = true
	c.TemplatePath = templatePath
	c.FilePath = filePath
	c.Data = data
	return
}
