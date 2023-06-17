package r2

import (
	"fmt"
	"os"

	"github.com/SkyeYoung/url-screenshot-service/internal/helper"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type r2 struct {
	cfg  *helper.Config
	sess *session.Session
}

type R2 interface {
	GetObject(key *string) (*s3.GetObjectOutput, error)
	HeadObject(key *string) (*s3.HeadObjectOutput, error)
	ListObjects(prefix *string, filter func(obj *s3.Object) bool) ([]*s3.Object, error)
	UploadObject(key *string) (*s3manager.UploadOutput, error)
	DeleteObject(key *string) (*s3.DeleteObjectOutput, error)
}

func New(cfg *helper.Config) R2 {
	return &r2{
		cfg: cfg,
		sess: session.Must(session.NewSession(&aws.Config{
			Credentials: credentials.NewStaticCredentials(cfg.AccessKey, cfg.SecretKey, ""),
			Endpoint:    aws.String(cfg.EndPoint),
			Region:      aws.String("us-east-1"),
		})),
	}
}

func (r *r2) GetObject(key *string) (*s3.GetObjectOutput, error) {
	svc := s3.New(r.sess)
	return svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(r.cfg.Bucket),
		Key:    aws.String(*key),
	})
}

func (r *r2) HeadObject(key *string) (*s3.HeadObjectOutput, error) {
	svc := s3.New(r.sess)
	return svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(r.cfg.Bucket),
		Key:    aws.String(*key),
	})
}

func (r *r2) ListObjects(prefix *string, filter func(obj *s3.Object) bool) ([]*s3.Object, error) {
	svc := s3.New(r.sess)
	var res []*s3.Object

	err := svc.ListObjectsV2Pages(&s3.ListObjectsV2Input{
		Bucket: aws.String(r.cfg.Bucket),
		Prefix: aws.String(*prefix),
	}, func(page *s3.ListObjectsV2Output, lastPage bool) bool {
		for _, obj := range page.Contents {
			if filter(obj) {
				res = append(res, obj)
			}
		}
		return !lastPage
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *r2) UploadObject(key *string) (*s3manager.UploadOutput, error) {
	file, err := os.Open(*key)
	if err != nil {
		return nil, fmt.Errorf("failed to open tmp file `%v`, err: %v", *key, err.Error())
	}
	uploader := s3manager.NewUploader(r.sess)
	res, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(r.cfg.Bucket),
		Key:    aws.String(*key),
		Body:   file,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload file `%v`, err: %v", *key, err.Error())
	}

	return res, nil
}

func (r *r2) DeleteObject(key *string) (*s3.DeleteObjectOutput, error) {
	svc := s3.New(r.sess)
	return svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(r.cfg.Bucket),
		Key:    aws.String(*key),
	})
}
