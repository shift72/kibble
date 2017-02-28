package render

import (
	"bytes"
	"fmt"

	"github.com/CloudyKit/jet"
	"github.com/indiereign/shift72-kibble/kibble/models"
)

// ConsoleRenderer - designed to render to the console for testing
type ConsoleRenderer struct {
	view *jet.Set
}

// Render - render to the console
func (c ConsoleRenderer) Render(route *models.Route, filePath string, data jet.VarMap) {

	w := bytes.NewBufferString("")
	w.Write([]byte("--------------------\n"))
	w.Write([]byte(fmt.Sprintf("FilePath: %s\n", filePath)))
	w.Write([]byte("--------------------\n"))

	t, err := c.view.GetTemplate(route.TemplatePath)
	if err != nil {
		w.Write([]byte("Template error\n"))
		w.Write([]byte(err.Error()))
		fmt.Println(w)
		return
	}

	if err = t.Execute(w, data, nil); err != nil {
		w.Write([]byte("Execute error\n"))
		w.Write([]byte(err.Error()))
	}

	fmt.Println(w)
}