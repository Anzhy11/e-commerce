package providers

import (
	"mime/multipart"
	"os"

	"github.com/anzhy11/go-e-commerce/internal/interfaces"
)

type S3UploadProvider struct {
	bucketName string
	region     string
}

func NewS3UploadProvider() interfaces.Upload {
	return &S3UploadProvider{
		bucketName: os.Getenv("AWS_BUCKET_NAME"),
		region:     os.Getenv("AWS_REGION"),
	}
}

func (s *S3UploadProvider) UploadFile(file *multipart.FileHeader, path string) (string, error) {
	return "", nil
}

func (s *S3UploadProvider) DeleteFile(filename string) error {
	return nil
}
