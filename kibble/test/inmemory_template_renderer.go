package test

import (
	"bytes"

	"github.com/CloudyKit/jet"
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
