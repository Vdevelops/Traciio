package file

import (
	"mime/multipart"
)

const (
	// MaxFileSize is 10MB
	MaxFileSize = 10 * 1024 * 1024
	// MaxImageWidth is maximum width for compressed images
	MaxImageWidth = 1920
	// MaxImageHeight is maximum height for compressed images
	MaxImageHeight = 1920
	// JPEGQuality for compressed images (1-100, lower = smaller file)
	JPEGQuality = 85
	// PNGQuality for compressed images (compression level 1-9, higher = smaller file)
	PNGCompressionLevel = 6
)

// Service wraps a StorageProvider and provides file operations
type Service struct {
	storage StorageProvider
}

// NewService creates a new Service with the given storage provider
func NewService(storage StorageProvider) *Service {
	return &Service{
		storage: storage,
	}
}

// UploadImage uploads and compresses an image file
func (s *Service) UploadImage(file *multipart.FileHeader) (string, error) {
	return s.storage.UploadImage(file)
}

// DeleteFile deletes a file from storage
func (s *Service) DeleteFile(filename string) error {
	return s.storage.DeleteFile(filename)
}

// GetFileURL returns the public URL for a file
func (s *Service) GetFileURL(filename string) string {
	return s.storage.GetFileURL(filename)
}
