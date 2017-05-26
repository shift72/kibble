package publish

import (
	"archive/zip"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Execute the publish process by
func Execute(rootPath string) error {

	os.RemoveAll(rootPath)

	target := path.Join(rootPath, "kibble-nibble.zip")

	err := createArchive(target)
	if err != nil {
		return err
	}

	//TODO: upload

	return err
}

func createArchive(target string) error {

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

	err = zipit(".kibble/dist/", archive, []string{"dist", "kibble-nibble.zip"})
	if err != nil {
		return err
	}
	return err
}

func zipit(source string, archive *zip.Writer, ignorePaths []string) error {

	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if ignorePath(ignorePaths, path) {
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
