package sync

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

// LocalStore -
type LocalStore struct {
	config Config
}

// NewLocalStore - create a new store
func NewLocalStore(config Config) (*LocalStore, error) {
	return &LocalStore{
		config: config,
	}, nil
}

// List - list the local files and their hashes
func (store *LocalStore) List() (FileRefCollection, error) {
	// get the current wd
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// change to the search dir
	err = os.Chdir(store.config.FileRootPath)
	if err != nil {
		return nil, err
	}
	defer os.Chdir(wd)

	fileList := []FileRef{}
	err = filepath.Walk(".", func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			fileList = append(fileList, calcMd5(store.config.FileRootPath, path))
		}
		return nil
	})

	return fileList, err
}

func calcMd5(rootPath, filePath string) FileRef {

	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	return FileRef{
		path: filePath,
		hash: fmt.Sprintf("%x", h.Sum(nil))}
}
