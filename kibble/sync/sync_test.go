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
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSave(t *testing.T) {

	saved := FileRefCollection{
		parseFileRef("ccc.html|ccc"),
		parseFileRef("ddd.html|ddd"),
	}

	reader := saved.GetReader()

	assert.Equal(t, 26, reader.Len())
}

func TestSaveAndLoad(t *testing.T) {

	loaded := FileRefCollection{}
	saved := FileRefCollection{
		parseFileRef("ccc.html|ccc"),
		parseFileRef("ddd.html|ddd"),
	}

	reader := saved.GetReader()

	err := loaded.Parse(reader)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(loaded))
}

func TestCompareNoChanges(t *testing.T) {
	local := []FileRef{
		parseFileRef("file.html|aaa"),
		parseFileRef("file2.html|bbb"),
	}

	remote := []FileRef{
		parseFileRef("file.html|aaa"),
		parseFileRef("file2.html|bbb"),
	}

	changes := compare(local, remote)

	if count(changes, ADD) != 0 {
		t.Error("expected add to be empty got ", count(changes, ADD))
	}

	if count(changes, REMOVE) != 0 {
		t.Error("expected remove to be empty", count(changes, REMOVE))
	}
}

func count(fileRefs []FileRef, action int) (count int) {
	for _, f := range fileRefs {
		if f.action == action {
			count++
		}
	}
	return
}

func TestCompareAddChanges(t *testing.T) {
	local := []FileRef{
		parseFileRef("file.html|aaa"),
		parseFileRef("file2.html|bbb"),
	}

	var remote []FileRef

	changes := compare(local, remote)

	if count(changes, ADD) != 2 {
		t.Error("expected add to be 2 got ", count(changes, ADD))
	}

	if count(changes, REMOVE) != 0 {
		t.Error("expected remove to be empty", count(changes, REMOVE))
	}
}

func TestCompareRemoveChanges(t *testing.T) {

	local := []FileRef{
		parseFileRef("file.html|aaa"),
	}

	remote := []FileRef{
		parseFileRef("file.html|aaa"),
		parseFileRef("file2.html|bbb"),
	}

	changes := compare(local, remote)

	if count(changes, ADD) != 0 {
		t.Error("expected add to be empty got ", count(changes, ADD))
	}

	if count(changes, REMOVE) != 1 {
		t.Error("expected remove to be 1 got ", count(changes, REMOVE))
	}
}

func TestCompareUpdateChanges(t *testing.T) {

	local := []FileRef{
		parseFileRef("file.html|ccc"),
		parseFileRef("file2.html|ddd"),
	}

	remote := []FileRef{
		parseFileRef("file.html|aaa"),
		parseFileRef("file2.html|bbb"),
	}

	changes := compare(local, remote)

	if count(changes, ADD) != 2 {
		t.Error("expected add to be 2 got ", count(changes, ADD))
	}

	if count(changes, REMOVE) != 0 {
		t.Error("expected remove to be 0 got ", count(changes, REMOVE))
	}
}

func TestCompareAddAndRemoveChanges(t *testing.T) {
	local := []FileRef{
		parseFileRef("file.html|aaa"),
	}

	remote := []FileRef{

		parseFileRef("file2.html|bbb"),
	}

	changes := compare(local, remote)

	if count(changes, ADD) != 1 {
		t.Error("expected add to be 1 got ", count(changes, ADD))
	}

	if count(changes, REMOVE) != 1 {
		t.Error("expected remove to be 1 got ", count(changes, REMOVE))
	}
}

type mockStore struct {
	returnErrors bool
}

func (store mockStore) Upload(wg *sync.WaitGroup, f FileRef) error {
	defer wg.Done()
	// fmt.Println("uploaded ", f)
	if store.returnErrors {
		return errors.New("bang")
	}
	return nil
}

func (store mockStore) UploadFileIndex(FileRefCollection) error {
	return nil
}

func (store mockStore) Delete(wg *sync.WaitGroup, f FileRef) error {
	defer wg.Done()
	// fmt.Println("deleted ", f)
	if store.returnErrors {
		return errors.New("bang")
	}
	return nil
}

func (store mockStore) List() (FileRefCollection, error) {
	fmt.Println("list ")
	return nil, nil
}

func TestSync(t *testing.T) {

	var store mockStore
	var local, remote []FileRef

	for i := 0; i < 50; i++ {
		local = append(local, add(fmt.Sprintf("file%d.html|ccc", i)))
	}

	PerformSync(store, local, remote)
}

func TestSyncWithErrors(t *testing.T) {

	store := mockStore{
		returnErrors: true,
	}

	var local, remote []FileRef

	for i := 0; i < 10; i++ {
		local = append(local, add(fmt.Sprintf("file%d.html|ccc", i)))
	}

	_, _, err := PerformSync(store, local, remote)
	if err == nil {
		t.Error("Expected errors")
	}
}

func add(raw string) (f FileRef) {
	f = parseFileRef(raw)
	f.action = ADD
	return
}

func remove(raw string) (f FileRef) {
	f = parseFileRef(raw)
	f.action = REMOVE
	return
}
