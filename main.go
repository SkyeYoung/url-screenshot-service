package main

import (
	"errors"
	"os"

	"github.com/SkyeYoung/url-screenshot-service/internal/helper"
	"github.com/SkyeYoung/url-screenshot-service/internal/screenshot"
	"github.com/gofiber/fiber/v2"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Website struct {
	Url string `json:"url"`
}

func main() {
	helper.SetupLogger()
	app := fiber.New()
	logger := helper.GetLogger()

	app.Get("/", func(c *fiber.Ctx) error {
		logger.Info("Hello, World 👋!")
		msg := "Hello, World 👋!"
		return c.SendString(msg)
	})

	app.Post("/", func(c *fiber.Ctx) error {
		url, err := helper.GetValidUrl(func() string {
			d := new(Website)
			if err := c.BodyParser(d); err != nil {
				return "invalid" // will cause error
			}
			return d.Url
		}())
		if err != nil {
			logger.Error(err)
			return err
		}

		sess := session.Must(session.NewSession(&aws.Config{
			Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
			Endpoint:    aws.String(endPoint),
			Region:      aws.String("us-east-1"),
		}))

		svc := s3.New(sess)
		info, err := svc.GetObjectAttributes(&s3.GetObjectAttributesInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(prefix + helper.EncodeImgNameAddExt(url)),
		})

		if err != nil {
		}
		info.LastModified

		logger.Infof("trying to get screeshot of %v", url)
		path, err := screenshot.Screenshot(url)
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return errors.New("failed to open file " + path)
		}
		uploader := s3manager.NewUploader(sess)
		res, err := uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(path),
			Body:   file,
		})
		if err != nil {
			return errors.New("failed to upload file " + path)
		}
		logger.Infof("file uploaded to, %s\n", aws.StringValue(&res.Location))

		return c.SendString(url)
	})

	logger.Fatal(app.Listen(":3004"))
}
