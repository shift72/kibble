package test

import (
	"github.com/CloudyKit/jet"
)

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
