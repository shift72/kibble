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

package sync

import (
	"crypto/md5"
	"fmt"
	"io"

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
	defer func() {
		if err := os.Chdir(wd); err != nil {
			log.Errorf("WARN: error while changing directory back to original %s", err)
		}
	}()

	var fileList []FileRef
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
