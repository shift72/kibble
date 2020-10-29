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
	"bytes"
	"io/ioutil"
	"mime"
	"path"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
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

	var sess *session.Session
	var s3Config *aws.Config

	if config.Profile != "" {
		sess = session.Must(session.NewSessionWithOptions(
			session.Options{
				SharedConfigState: session.SharedConfigEnable,
				Profile:           config.Profile,
			},
		))

		s3Config = &aws.Config{
			Region: aws.String(config.Region),
		}
	} else {
		sess = session.Must(session.NewSession())
		s3Config = &aws.Config{
			Credentials: ec2rolecreds.NewCredentials(sess),
			Region:      aws.String(config.Region),
		}
	}

	return &S3Store{
		config: config,
		svc:    s3.New(sess, s3Config),
	}, nil
}

// List - list all files
func (store *S3Store) List() (FileRefCollection, error) {
	var fileList FileRefCollection

	// attempt to download an index file
	foundIndexObject, err := store.svc.GetObject(&s3.GetObjectInput{
		Bucket: &store.config.Bucket,
		Key:    aws.String(store.config.BucketRootPath + "index.kibble"),
	})

	if err == nil {
		defer foundIndexObject.Body.Close()
		err = fileList.Parse(foundIndexObject.Body)
		if err == nil {
			return fileList, nil
		}
	}

	// if the index file doesnt exist, then this must be a new site
	// so just return an empty list, so everything is uploaded
	return fileList, err
}

// Upload - upload file to s3
func (store *S3Store) Upload(wg *sync.WaitGroup, f FileRef) error {
	defer wg.Done()

	src := store.config.FileRootPath + f.path
	dst := store.config.BucketRootPath + f.path

	b, err := ioutil.ReadFile(src)
	if err != nil {
		log.Error("upload file read", err)
		return err
	}

	_, err = store.svc.PutObject(
		&s3.PutObjectInput{
			Bucket:      &store.config.Bucket,
			Key:         &dst,
			Body:        bytes.NewReader(b),
			ACL:         aws.String("public-read"),
			ContentType: aws.String(mime.TypeByExtension(path.Ext(f.path))),
		},
	)
	if err != nil {
		log.Error("upload err", err)
		return err
	}

	return nil
}

// UploadFileIndex - the file index
func (store *S3Store) UploadFileIndex(f FileRefCollection) error {
	_, err := store.svc.PutObject(
		&s3.PutObjectInput{
			Bucket:      &store.config.Bucket,
			Key:         aws.String(store.config.BucketRootPath + "index.kibble"),
			Body:        f.GetReader(),
			ACL:         aws.String("public-read"),
			ContentType: aws.String("text/plain"),
		},
	)
	return err
}

// Delete - delete file from S3
func (store *S3Store) Delete(wg *sync.WaitGroup, f FileRef) error {
	defer wg.Done()

	dst := store.config.BucketRootPath + f.path

	//TODO: bulk delete

	_, err := store.svc.DeleteObject(
		&s3.DeleteObjectInput{
			Bucket: &store.config.Bucket,
			Key:    &dst,
		},
	)
	if err != nil {
		log.Error("upload err", err)
		return err
	}

	return nil
}
