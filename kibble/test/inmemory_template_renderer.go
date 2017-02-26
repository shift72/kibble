package test

import (
	"bytes"

	"github.com/CloudyKit/jet"
	"github.com/indiereign/shift72-kibble/kibble/models"
)

type InMemoryTemplateRenderer struct {
	View     *jet.Set
	Template *jet.Template
	Result   InMemoryResult
}

// Render - render the pages to memory
func (c *InMemoryTemplateRenderer) Render(route *models.Route, filePath string, data jet.VarMap) {

	c.Result = InMemoryResult{
		buffer:   bytes.NewBufferString(""),
		filePath: filePath,
	}

	if err := c.Template.Execute(c.Result.buffer, data, nil); err != nil {
		c.Result.err = err
	}
}
