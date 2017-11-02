package sync

import (
	"bytes"
	"io/ioutil"
	"mime"
	"path"
	"strings"
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

	err = store.svc.ListObjectsPages(&s3.ListObjectsInput{
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
