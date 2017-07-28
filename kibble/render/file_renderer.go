package render

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/CloudyKit/jet"
	"github.com/indiereign/shift72-kibble/kibble/models"
	"github.com/indiereign/shift72-kibble/kibble/utils"
)

// FileRenderer - designed to render to the file system for testing
type FileRenderer struct {
	view     *jet.Set
	rootPath string
}

// Initialise - start the rendering process
func (c FileRenderer) Initialise() {
	os.RemoveAll(c.rootPath)

	err := utils.CopyDir(staticFolder, c.rootPath)
	if err != nil {
		log.Warningf("Warn: static folder copy failed %s", err)
	}

	// copy language files too, they are a special file name format
	cd, _ := os.Getwd()
	glob := filepath.Join(cd, "/*.all.json")
	langFiles, err := filepath.Glob(glob)
	if len(langFiles) > 0 {
		for _, file := range langFiles {
			dst := filepath.Join(c.rootPath, filepath.Base(file))
			err := utils.CopyFile(file, dst)
			if err != nil {
				log.Warningf("Warn: language file (%s) copy failed %s", file, err)
			}
		}
	}
}

// Render - render to the console
func (c FileRenderer) Render(route *models.Route, filePath string, data jet.VarMap) (errorCount int) {
	defer func() {
		if r := recover(); r != nil {
			errorCount++
			log.Debug("Error. Recovered from ", r)
		}
	}()

	fullPath := path.Join(c.rootPath, filePath)
	if strings.HasSuffix(filePath, "/") {
		fullPath = path.Join(fullPath, "index.html")
	}

	log.Debugf("FilePath: %s", fullPath)

	w := bytes.NewBufferString("")
	t, err := c.view.GetTemplate(route.TemplatePath)
	if err != nil {
		errorCount++
		log.Error("Template load error", err)
		return
	}

	if err = t.Execute(w, data, nil); err != nil {
		errorCount++
		w.WriteString("<pre>")
		w.WriteString(err.Error())
		w.WriteString("</pre>")

		log.Errorf("Template execute error: %s", err)
	}

	dirPath := filepath.Dir(fullPath)
	os.MkdirAll(dirPath, 0777)

	// optional check
	if _, err := os.Stat(fullPath); err == nil {
		log.Warningf("File exists and will be overwritten: %s", fullPath)
	}

	err = ioutil.WriteFile(fullPath, w.Bytes(), 0777)
	if err != nil {
		errorCount++
		log.Error("File write:", err)
	}

	return
}
