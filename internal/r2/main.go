package r2

import (
	"errors"
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

func (r *r2) UploadObject(key *string) (*s3manager.UploadOutput, error) {
	file, err := os.Open(*key)
	if err != nil {
		return nil, errors.New("failed to open tmp file " + *key)
	}
	uploader := s3manager.NewUploader(r.sess)
	res, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(r.cfg.Bucket),
		Key:    aws.String(*key),
		Body:   file,
	})
	if err != nil {
		return nil, errors.New("failed to upload file " + *key)
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
