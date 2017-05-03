package render

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/CloudyKit/jet"
	"github.com/indiereign/shift72-kibble/kibble/models"
)

// FileRenderer - designed to render to the file system for testing
type FileRenderer struct {
	view     *jet.Set
	rootPath string
}

// Initialise - start the rendering process
func (c FileRenderer) Initialise() {
	os.RemoveAll(c.rootPath)

	err := CopyDir(staticFolder, c.rootPath)
	if err != nil {
		log.Warningf("Warn: static folder copy failed %s", err)
	}
}

// Render - render to the console
func (c FileRenderer) Render(route *models.Route, filePath string, data jet.VarMap) {

	fullPath := path.Join(c.rootPath, filePath)
	if strings.HasSuffix(filePath, "/") {
		fullPath = path.Join(fullPath, "index.html")
	}

	dirPath := filepath.Dir(fullPath)

	log.Debugf("FilePath: %s", fullPath)

	w := bytes.NewBufferString("")
	t, err := c.view.GetTemplate(route.TemplatePath)
	if err != nil {
		fmt.Println("ERROR: Template load error", err)
		return
	}

	if err = t.Execute(w, data, nil); err != nil {
		w.WriteString("<pre>")
		w.WriteString(err.Error())
		w.WriteString("</pre>")

		log.Errorf("Template execute error: %s", err)
	}

	os.MkdirAll(dirPath, 0777)

	// optional check
	if _, err := os.Stat(fullPath); err == nil {
		log.Warningf("File exists and will be overwritten: %s", fullPath)
	}

	err = ioutil.WriteFile(fullPath, w.Bytes(), 0777)
	if err != nil {
		log.Error("File write:", err)
	}
}
