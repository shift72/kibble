package models

import "github.com/CloudyKit/jet"

// Film - represents a film
type Film struct {
	ID       int
	Slug     string
	Title    string
	Synopsis string
}

// Renderer - rendering implementation
type Renderer interface {
	Render(route *Route, filePath string, data jet.VarMap)
}
