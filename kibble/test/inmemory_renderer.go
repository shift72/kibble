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
	"fmt"

	"github.com/CloudyKit/jet"
)

// InMemoryRenderer - render to memory for testing
type InMemoryRenderer struct {
	View    *jet.Set
	Results []InMemoryResult
}

// InMemoryResult - in memory result
type InMemoryResult struct {
	buffer   *bytes.Buffer
	filePath string
	err      error
}

// Output - return the result
func (r *InMemoryResult) Output() string {
	if r.err == nil {
		return fmt.Sprintf("%s", r.buffer)
	}
	return fmt.Sprintf("error: %s\n", r.err)
}

// ErrorCount - return the number of errors
func (c *InMemoryRenderer) ErrorCount() int {
	i := 0
	for _, r := range c.Results {
		if r.err != nil {
			i++
		}
	}
	return i
}

// DumpErrors - print the errors
func (c *InMemoryRenderer) DumpErrors() {
	for _, r := range c.Results {
		if r.err != nil {
			fmt.Printf("Error found on %s - %s\n", r.filePath, r.err)
		}
	}
}

// DumpResults - print the results
func (c *InMemoryRenderer) DumpResults() {
	for _, r := range c.Results {
		fmt.Printf("---- %s start ----\n", r.filePath)
		fmt.Printf(r.Output())
		fmt.Printf("---- %s end ----\n", r.filePath)
	}
}

// Initialise -
func (c *InMemoryRenderer) Initialise() {

}

// Render - render the pages to memory
func (c *InMemoryRenderer) Render(templatePath string, filePath string, data jet.VarMap) (errCount int) {

	if c.Results == nil {
		c.Results = make([]InMemoryResult, 0, 10)
	}

	result := InMemoryResult{
		buffer:   bytes.NewBufferString(""),
		filePath: filePath,
	}

	c.Results = append(c.Results, result)

	t, err := c.View.GetTemplate(templatePath)
	if err != nil {
		errCount++
		result.err = err
		return
	}

	if err = t.Execute(result.buffer, data, nil); err != nil {
		errCount++
		result.err = err
	}

	return
}
