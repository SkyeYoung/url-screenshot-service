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

var (
	sess *session.Session
)

func convertToPtrSlice(slice []string) []*string {
	ptrSlice := make([]*string, len(slice))
	for i, val := range slice {
		ptrSlice[i] = &val
	}
	return ptrSlice
}

func SetupSession(cfg *helper.Config) {
	sess = session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(cfg.AccessKey, cfg.SecretKey, ""),
		Endpoint:    aws.String(cfg.EndPoint),
		Region:      aws.String("us-east-1"),
	}))
}

func GetObjectAttributes(cfg *helper.Config, key *string) (*s3.GetObjectOutput, error) {
	svc := s3.New(sess)
	return svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(cfg.Bucket),
		Key:    aws.String(*key),
	})
}

func Upload(cfg *helper.Config, key *string) (*s3manager.UploadOutput, error) {
	file, err := os.Open(*key)
	if err != nil {
		return nil, errors.New("failed to open tmp file " + *key)
	}
	uploader := s3manager.NewUploader(sess)
	res, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(cfg.Bucket),
		Key:    aws.String(*key),
		Body:   file,
	})
	if err != nil {
		return nil, errors.New("failed to upload file " + *key)
	}

	return res, nil
}
