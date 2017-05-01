package sync

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3Store - s3 store
type S3Store struct {
	config Config
	svc    *s3.S3
}

// NewS3Store - create a new store
func NewS3Store(config Config) (*S3Store, error) {

	sess := session.Must(session.NewSessionWithOptions(
		session.Options{
			Profile: config.Profile,
		},
	))
	sess.Config.Region = &config.Region

	return &S3Store{
		config: config,
		svc:    s3.New(sess),
	}, nil
}

// List - list all files
func (store *S3Store) List() (FileRefCollection, error) {
	fileList := []FileRef{}

	err := store.svc.ListObjectsPages(&s3.ListObjectsInput{
		Bucket: &store.config.Bucket,
	}, func(p *s3.ListObjectsOutput, last bool) (shouldContinue bool) {
		for _, obj := range p.Contents {
			path := *obj.Key
			if strings.HasPrefix(path, store.config.BucketRootPath) {
				fileList = append(fileList, FileRef{
					path: path[len(store.config.BucketRootPath):len(path)],
					hash: strings.Trim(*obj.ETag, "\""),
				})
			}
		}
		return true
	})

	return fileList, err
}

// Upload - upload file to s3
func (store *S3Store) Upload(wg *sync.WaitGroup, f FileRef) error {
	defer wg.Done()

	src := store.config.FileRootPath + f.path
	dst := store.config.BucketRootPath + f.path

	// fmt.Printf("upload %s to %s\n", src, dst)

	b, err := ioutil.ReadFile(src)
	if err != nil {
		fmt.Println("Upload file read", err)
		return err
	}

	_, err = store.svc.PutObject(
		&s3.PutObjectInput{
			Bucket: &store.config.Bucket,
			Key:    &dst,
			Body:   bytes.NewReader(b),
			//	ACL: // defaults to read only
		},
	)
	if err != nil {
		fmt.Println("Upload err", err)
		return err
	}

	return nil
}

// Delete - delete file from S3
func (store *S3Store) Delete(wg *sync.WaitGroup, f FileRef) error {
	defer wg.Done()

	dst := store.config.BucketRootPath + f.path

	//TODO: bulk delete

	// fmt.Printf("delete remote file %s\n", dst)

	_, err := store.svc.DeleteObject(
		&s3.DeleteObjectInput{
			Bucket: &store.config.Bucket,
			Key:    &dst,
		},
	)
	if err != nil {
		fmt.Println("Upload err", err)
		return err
	}

	return nil
}
