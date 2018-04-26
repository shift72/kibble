package sync

import (
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"strings"
	"sync"
	"time"

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

// Summary of the sync process
type Summary struct {
	FilesRemoved            int           `json:"filesRemoved"`
	FilesAdded              int           `json:"filesAdded"`
	FilesTotal              int           `json:"filesTotal"`
	RenderDuration          time.Duration `json:"renderDuration"`
	ChangesDetectedDuration time.Duration `json:"changesDetectedDuration"`
	UploadDuration          time.Duration `json:"uploadDuration"`
	Errors                  []string      `json:"errors"`
}

// Execute - start a sync
func Execute(config Config) (*Summary, error) {

	s3Store, err := NewS3Store(config)
	if err != nil {
		return nil, err
	}
	localStore, err := NewLocalStore(config)
	if err != nil {
		return nil, err
	}

	swDetect := utils.NewStopwatch("detect")
	remote, err := s3Store.List()
	if err != nil {
		return nil, err
	}

	local, err := localStore.List()
	if err != nil {
		return nil, err
	}
	changes := compare(local, remote)
	detect := swDetect.Completed()

	swSync := utils.NewStopwatch("sync files")
	added, removed, err := PerformSync(s3Store, changes)
	if err != nil {
		return nil, err
	}
	upload := swSync.Completed()

	s := &Summary{
		FilesAdded:              added,
		FilesRemoved:            removed,
		FilesTotal:              len(local),
		ChangesDetectedDuration: detect,
		UploadDuration:          upload,
		Errors:                  make([]string, 0),
	}

	return s, err
}

// ToJSON renders the summary to json
func (s *Summary) ToJSON() string {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		log.Errorf("failed to serialize. %s", err)
		return "{\"error\":\"unable to render json\"}"
	}
	return string(data)
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
	sourcePath := cfg.SourcePath()

	render.Render(sourcePath, sample1Path, cfg)

	sample1, err := local.List()
	if err != nil {
		log.Error("sync error", err)
		return err
	}

	var sample2Path = path.Join(".kibble", "build-sample-2")

	render.Render(sourcePath, sample2Path, cfg)

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
func PerformSync(store Store, changes []FileRef) (added int, removed int, err error) {

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

	added = 0
	removed = 0

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
		return added, removed, errors.New("Unable to sync files")
	}

	log.Infof("sync successful [added: %d][removed: %d]", added, removed)
	return added, removed, nil
}
