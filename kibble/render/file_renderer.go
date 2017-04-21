package render

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/CloudyKit/jet"
	"github.com/indiereign/shift72-kibble/kibble/models"
)

// FileRenderer - designed to render to the file system for testing
type FileRenderer struct {
	view        *jet.Set
	rootPath    string
	showSummary bool
}

// Initialise - start the rendering process
func (c FileRenderer) Initialise() {
	os.RemoveAll(c.rootPath)
}

// Render - render to the console
func (c FileRenderer) Render(route *models.Route, filePath string, data jet.VarMap) {

	fullPath := c.rootPath + filePath
	if strings.HasSuffix(fullPath, "/") {
		fullPath = fullPath + "index.html"
	}

	dirPath := filepath.Dir(fullPath)

	if c.showSummary {
		fmt.Printf("FilePath: %s\n", fullPath)
	}

	w := bytes.NewBufferString("")
	t, err := c.view.GetTemplate(route.TemplatePath)
	if err != nil {
		if c.showSummary {
			fmt.Println("Template load error", err)
		}
		return
	}

	if err = t.Execute(w, data, nil); err != nil {
		if c.showSummary {
			fmt.Println("Template execute error", err)
		}
	}

	os.MkdirAll(dirPath, 0777)
	err = ioutil.WriteFile(fullPath, w.Bytes(), 0777)
	if err != nil {
		fmt.Println(err)
	}
}
