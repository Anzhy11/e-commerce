package providers

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	appConfig "github.com/anzhy11/go-e-commerce/internal/config"
	"github.com/anzhy11/go-e-commerce/internal/interfaces"
)

type S3UploadProvider struct {
	client     *s3.Client
	uploader   *manager.Uploader
	bucketName string
	endpoint   string
}

func NewS3UploadProvider(cfg *appConfig.Config) interfaces.Upload {
	awsConfig, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cfg.AWS.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AWS.AccessKeyID,
			cfg.AWS.SecretAccessKey,
			"",
		)),
	)

	if err != nil {
		panic("Failed to load AWS config: " + err.Error())
	}

	client := s3.NewFromConfig(awsConfig, func(o *s3.Options) {
		if cfg.AWS.S3Endpoint != "" {
			o.BaseEndpoint = aws.String(cfg.AWS.S3Endpoint)
			o.UsePathStyle = true
		}
	})

	return &S3UploadProvider{
		client:     client,
		uploader:   manager.NewUploader(client),
		bucketName: cfg.AWS.S3Bucket,
		endpoint:   cfg.AWS.S3Endpoint,
	}
}

func (s *S3UploadProvider) UploadFile(file *multipart.FileHeader, path string) (string, error) {
	log.Println("Uploading file:", file.Filename)

	src, srcErr := file.Open()
	if srcErr != nil {
		return "", srcErr
	}
	defer func() {
		if err := src.Close(); err != nil {
			fmt.Println("Error closing file:", err)
		}
	}()

	result, err := s.uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(path + file.Filename),
		Body:   src,
	})

	if err != nil {
		return "", err
	}

	return *result.Key, nil
}

func (s *S3UploadProvider) DeleteFile(filename string) error {
	_, err := s.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(filename),
	})

	return err
}
