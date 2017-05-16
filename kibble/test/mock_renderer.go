package test

import (
	"github.com/CloudyKit/jet"
	"github.com/indiereign/shift72-kibble/kibble/models"
)

type MockRenderer struct {
	InitialisedCalled bool
	RenderCalled      bool
	Route             *models.Route
	FilePath          string
	Data              jet.VarMap
}

// Initialise - start the rendering process
func (c *MockRenderer) Initialise() {
	c.InitialisedCalled = true
}

// Render - render to the console
func (c *MockRenderer) Render(route *models.Route, filePath string, data jet.VarMap) (errCount int) {
	c.RenderCalled = true
	c.Route = route
	c.FilePath = filePath
	c.Data = data
	return
}
