package test

import (
	"bytes"

	"github.com/CloudyKit/jet"
	"github.com/indiereign/shift72-kibble/kibble/models"
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
func (c *InMemoryTemplateRenderer) Render(route *models.Route, filePath string, data jet.VarMap) (errCount int) {

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
