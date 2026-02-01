package providers

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/anzhy11/go-e-commerce/internal/interfaces"
)

type LocalUploadProvider struct {
	basePath string
}

func NewLocalUploadProvider(basePath string) interfaces.Upload {
	return &LocalUploadProvider{
		basePath: basePath,
	}
}

func (l *LocalUploadProvider) UploadFile(file *multipart.FileHeader, path string) (string, error) {
	fullPath := filepath.Join(l.basePath, path)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return "", err
	}

	src, srcErr := file.Open()
	if srcErr != nil {
		return "", srcErr
	}
	defer func() {
		if err := src.Close(); err != nil {
			fmt.Println("Error closing file:", err)
		}
	}()

	dest, destErr := os.Create(fullPath)
	if destErr != nil {
		return "", destErr
	}
	defer func() {
		if err := dest.Close(); err != nil {
			fmt.Println("Error closing file:", err)
		}
	}()

	if _, err := dest.ReadFrom(src); err != nil {
		return "", err
	}

	return fmt.Sprintf("/uploads/%s", path), nil
}

func (l *LocalUploadProvider) DeleteFile(filename string) error {
	fullPath := filepath.Join(l.basePath, filename)
	return os.Remove(fullPath)
}
