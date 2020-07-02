//    Copyright 2018 SHIFT72
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package publish

import (
	"archive/zip"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"kibble/api"
	"kibble/config"
	"kibble/models"
	"kibble/utils"
)

// Execute the publish process by
func Execute(sourcePath string, buildPath string, cfg *models.Config, zipOnly bool) error {

	err := Validate(sourcePath, cfg)
	if err != nil {
		return err
	}

	err = config.CheckVersion(cfg)
	if err != nil {
		return err
	}

	target := path.Join(buildPath, "kibble-nibble.zip")

	// unzipping wigs out on the build server unless we specifically ignore "/",
	// for some reason an entry is added for "/" (even though its empty) in the zip file unless we ignore it
	ignorePatterns := append(cfg.LiveReload.IgnoredPaths, "/")

	ignoredPaths := utils.NewFileIgnorer(sourcePath, ignorePatterns)

	err = createArchive(target, sourcePath, ignoredPaths)
	if err != nil {
		return err
	}

	if zipOnly {
		log.Info("zip only successful")
		return nil
	}

	api.CheckAdminCredentials(cfg)

	extraParams := map[string]string{
		"name":           cfg.Name,
		"version":        cfg.Version,
		"builderVersion": cfg.BuilderVersion,
	}

	log.Infof("uploading name: %s@%s built with %s", cfg.Name, cfg.Version, cfg.BuilderVersion)
	err = api.Upload(cfg, cfg.SiteURL+"/services/users/v1/site_templates", extraParams, target)
	if err != nil {
		return err
	}
	log.Info("upload successful")
	return nil
}

func createArchive(target string, sourcePath string, ignoredPaths utils.FileIgnorer) error {

	targetPath := filepath.Dir(target)

	err := os.MkdirAll(targetPath, 0777)
	if err != nil {
		return err
	}

	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	err = zipit(sourcePath, archive, ignoredPaths)
	if err != nil {
		return err
	}

	return err
}

func zipit(source string, archive *zip.Writer, ignoredPaths utils.FileIgnorer) error {

	info, err := os.Stat(source)
	if err != nil && os.IsNotExist(err) {
		os.MkdirAll(source, 0777)
	}

	info, err = os.Stat(source)
	if err != nil {
		return err
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == "/" || ignoredPaths.IsIgnored(path) {
			return nil
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = strings.TrimPrefix(path, source)
		}

		if info.IsDir() {
			header.Name += string(os.PathSeparator)
		} else {
			header.Method = zip.Deflate
		}

		if header.Name == "/" {
			return nil
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})

	return err
}

func ignorePath(ignorePaths []string, name string) bool {
	for _, c := range ignorePaths {
		if strings.HasPrefix(name, c) || strings.HasSuffix(name, c) {
			return true
		}
	}
	return false
}
