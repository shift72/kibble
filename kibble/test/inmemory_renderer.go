package test

import (
	"bytes"
	"fmt"

	"github.com/CloudyKit/jet"
	"github.com/indiereign/shift72-kibble/kibble/models"
)

type InMemoryRenderer struct {
	View    *jet.Set
	Results []InMemoryResult
}

type InMemoryResult struct {
	buffer   *bytes.Buffer
	filePath string
	err      error
}

func (c *InMemoryRenderer) ErrorCount() int {
	i := 0
	for _, r := range c.Results {
		if r.err != nil {
			i++
		}
	}
	return i
}

func (c *InMemoryRenderer) DumpErrors() {
	for _, r := range c.Results {
		if r.err != nil {
			fmt.Printf("Error found on %s - %s\n", r.filePath, r.err)
		}
	}
}

func (c *InMemoryRenderer) DumpResults() {
	for _, r := range c.Results {
		if r.err == nil {
			fmt.Printf("---- %s ----\n%s\n---- -----\n", r.filePath, r.buffer)
		} else {
			fmt.Printf("---- %s ----\n error: %s\n---- -----\n", r.filePath, r.err)
		}
	}
}

// Render - render the pages to memory
func (c *InMemoryRenderer) Render(route *models.Route, filePath string, data jet.VarMap) {

	if c.Results == nil {
		c.Results = make([]InMemoryResult, 1, 10)
	}

	result := InMemoryResult{
		buffer:   bytes.NewBufferString(""),
		filePath: filePath,
	}

	c.Results = append(c.Results, result)

	t, err := c.View.GetTemplate(route.TemplatePath)
	if err != nil {
		result.err = err
		return
	}

	if err = t.Execute(result.buffer, data, nil); err != nil {
		result.err = err
	}
}
