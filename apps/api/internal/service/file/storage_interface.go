package file

import "mime/multipart"

// StorageProvider defines the interface for storage implementations
type StorageProvider interface {
	// UploadImage uploads and compresses an image file, returns the public URL
	UploadImage(file *multipart.FileHeader) (string, error)
	// DeleteFile deletes a file from storage by filename
	DeleteFile(filename string) error
	// GetFileURL returns the public URL for a file
	GetFileURL(filename string) string
}



