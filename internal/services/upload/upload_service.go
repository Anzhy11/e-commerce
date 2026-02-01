package uploadService

import (
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/anzhy11/go-e-commerce/internal/interfaces"
)

type UploadService struct {
	provider interfaces.Upload
}

func NewUploadService(provider interfaces.Upload) *UploadService {
	return &UploadService{
		provider: provider,
	}
}

func (s *UploadService) UploadProductImage(productID uint, file *multipart.FileHeader) (string, error) {
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !isValidImageExt(ext) {
		return "", errors.New("invalid file type")
	}

	randomName := fmt.Sprintf("%d%s", time.Now().Unix(), ext)

	path := fmt.Sprintf("products/%d/%s", productID, randomName)

	return s.provider.UploadFile(file, path)
}

func (s *UploadService) DeleteFile(filename string) error {
	return s.provider.DeleteFile(filename)
}

// Helper
func isValidImageExt(ext string) bool {
	return ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".gif" || ext == ".webp"
}
