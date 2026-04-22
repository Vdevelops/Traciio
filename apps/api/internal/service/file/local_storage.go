package file

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/image/draw"
)

// LocalStorage implements StorageProvider for local file system storage
type LocalStorage struct {
	uploadDir string
	baseURL   string
}

// NewLocalStorage creates a new LocalStorage instance
func NewLocalStorage(uploadDir, baseURL string) *LocalStorage {
	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create upload directory: %v", err))
	}

	return &LocalStorage{
		uploadDir: uploadDir,
		baseURL:   baseURL,
	}
}

// UploadImage uploads and compresses an image file
func (s *LocalStorage) UploadImage(file *multipart.FileHeader) (string, error) {
	// Validate file size
	if file.Size > MaxFileSize {
		return "", fmt.Errorf("file size exceeds maximum allowed size of %d bytes", MaxFileSize)
	}

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Detect image format
	contentType := file.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return "", fmt.Errorf("file is not an image: %s", contentType)
	}

	// Decode image directly from the file
	img, format, err := image.Decode(src)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	// Resize image if needed
	img = resizeImage(img)

	// Generate unique filename
	filename := generateFilename(format)
	filePath := filepath.Join(s.uploadDir, filename)

	// Create output file
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	// Encode and compress image
	if err := encodeImage(dst, img, format); err != nil {
		os.Remove(filePath) // Clean up on error
		return "", fmt.Errorf("failed to encode image: %w", err)
	}

	// Return URL
	return s.GetFileURL(filename), nil
}

// DeleteFile deletes a file from storage
func (s *LocalStorage) DeleteFile(filename string) error {
	filePath := filepath.Join(s.uploadDir, filename)
	return os.Remove(filePath)
}

// GetFileURL returns the public URL for a file
func (s *LocalStorage) GetFileURL(filename string) string {
	return fmt.Sprintf("%s/%s", s.baseURL, filename)
}

// resizeImage resizes image if it exceeds maximum dimensions
func resizeImage(img image.Image) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Check if resize is needed
	if width <= MaxImageWidth && height <= MaxImageHeight {
		return img
	}

	// Calculate new dimensions maintaining aspect ratio
	var newWidth, newHeight int
	if width > height {
		newWidth = MaxImageWidth
		newHeight = int(float64(height) * float64(MaxImageWidth) / float64(width))
	} else {
		newHeight = MaxImageHeight
		newWidth = int(float64(width) * float64(MaxImageHeight) / float64(height))
	}

	// Create new image with calculated dimensions
	dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	draw.ApproxBiLinear.Scale(dst, dst.Bounds(), img, bounds, draw.Src, nil)

	return dst
}

// encodeImage encodes image to file with compression
func encodeImage(w io.Writer, img image.Image, format string) error {
	switch format {
	case "jpeg", "jpg":
		return jpeg.Encode(w, img, &jpeg.Options{Quality: JPEGQuality})
	case "png":
		encoder := &png.Encoder{CompressionLevel: PNGCompressionLevel}
		return encoder.Encode(w, img)
	default:
		// For other formats, try to encode as JPEG
		return jpeg.Encode(w, img, &jpeg.Options{Quality: JPEGQuality})
	}
}

// generateFilename generates a unique filename
func generateFilename(format string) string {
	// Normalize format
	format = strings.ToLower(format)
	if format == "jpg" {
		format = "jpeg"
	}

	// Generate UUID and timestamp
	id := uuid.New().String()
	timestamp := time.Now().Format("20060102-150405")
	return fmt.Sprintf("%s-%s.%s", timestamp, id[:8], format)
}



