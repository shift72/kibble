package sync

import (
	"errors"
	"fmt"
	"path"
	"strings"
	"sync"

	"github.com/indiereign/shift72-kibble/kibble/models"
	"github.com/indiereign/shift72-kibble/kibble/render"
	"github.com/indiereign/shift72-kibble/kibble/utils"
)

const (
	// ADD - addition detected
	ADD = 1
	// REMOVE - removal detected
	REMOVE = 2
)

// FileRef - reference to a file could be local or remote
type FileRef struct {
	path   string
	hash   string
	action int
}

// FileRefCollection - a list of file refs
type FileRefCollection []FileRef

// Config - sync configuration
type Config struct {
	Profile        string
	Region         string
	Bucket         string
	BucketRootPath string
	FileRootPath   string
}

// Store - a place to store the site files
type Store interface {
	List() (FileRefCollection, error)
	Upload(wg *sync.WaitGroup, file FileRef) error
	Delete(wg *sync.WaitGroup, file FileRef) error
}

// Sync -
type Sync struct {
	Config Config
	Store  Store
}

// Execute - start a sync
func Execute(config Config) error {

	s3Store, _ := NewS3Store(config)
	localStore, _ := NewLocalStore(config)

	swDetect := utils.NewStopwatch("detect")
	remote, _ := s3Store.List()
	local, _ := localStore.List()
	changes := compare(local, remote)
	swDetect.Completed()

	swSync := utils.NewStopwatch("sync files")
	err := PerformSync(s3Store, changes)
	swSync.Completed()

	return err
}

// TestIdempotent - run the sync twice and check for differences
func TestIdempotent(config Config, cfg *models.Config) error {

	utils.ConfigureStandardLogging(false)

	local, err := NewLocalStore(config)
	if err != nil {
		log.Error("sync error", err)
		return err
	}

	var sample1Path = path.Join(".kibble", "build-sample-1")

	render.Render(sample1Path, cfg)

	sample1, err := local.List()
	if err != nil {
		log.Error("sync error", err)
		return err
	}

	var sample2Path = path.Join(".kibble", "build-sample-2")

	render.Render(sample2Path, cfg)

	sample2, err := local.List()
	if err != nil {
		log.Error("sync error", err)
		return err
	}

	diff := compare(sample1, sample2)

	if len(diff) > 0 {
		log.Errorf("test failed. files have changed between executions")
		log.Errorf("see folders %s and %s for more detail.", sample1Path, sample2Path)
		diff.Print()
	}

	return nil
}

// compare - compare lists
func compare(local, remote FileRefCollection) (changes FileRefCollection) {

	found := false

	// add
	for l := 0; l < len(local); l++ {
		found = false
		for r := 0; r < len(remote); r++ {
			if local[l] == remote[r] {
				found = true
				break
			}
		}
		if !found {
			local[l].action = ADD
			changes = append(changes, local[l])
		}
	}

	// remove
	for r := 0; r < len(remote); r++ {
		found = false
		for l := 0; l < len(local); l++ {
			// only check path
			if local[l].path == remote[r].path {
				found = true
				break
			}
		}
		if !found {
			remote[r].action = REMOVE
			changes = append(changes, remote[r])
		}
	}

	return
}

func parseFileRef(raw string) FileRef {
	p := strings.Split(raw, "|")
	return FileRef{path: p[0], hash: p[1]}
}

// String - convert file ref to string
func (f *FileRef) String() string {
	return fmt.Sprintf("%s|%s", f.path, f.hash)
}

// Print -
func (c *FileRefCollection) Print() {
	for _, f := range *c {
		fmt.Println(f.String())
	}
}

// PerformSync - sync all files to the remote server
func PerformSync(store Store, changes []FileRef) error {

	concurrency := 20

	var wg sync.WaitGroup
	work := make(chan FileRef, concurrency)
	errorChan := make(chan error, concurrency)

	for w := 0; w < concurrency; w++ {
		go func(id int) {
			for {
				j, more := <-work
				if !more {
					return
				}

				if j.action == ADD {
					err := store.Upload(&wg, j)
					if err != nil {
						errorChan <- err
					}
				} else if j.action == REMOVE {
					err := store.Delete(&wg, j)
					if err != nil {
						errorChan <- err
					}
				}
			}
		}(w)
	}

	added := 0
	removed := 0

	// queue the work
	for _, f := range changes {
		if f.action == ADD {
			log.Debugf("added: %s", f.path)
			added++
		}
		if f.action == REMOVE {
			log.Debugf("removed: %s", f.path)
			removed++
		}
		wg.Add(1)
		work <- f
	}

	// wait for a result
	close(work)
	wg.Wait()
	close(errorChan)

	if len(errorChan) > 0 {
		return errors.New("Unable to sync files")
	}

	log.Infof("sync successful [added: %d][removed: %d]", added, removed)
	return nil
}
