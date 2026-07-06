package file

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"io"
	"mime/multipart"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// R2Storage implements StorageProvider for Cloudflare R2 storage
type R2Storage struct {
	client    *s3.Client
	bucket    string
	publicURL string
	baseURL   string // For backward compatibility, can be used as prefix
}

// NewR2Storage creates a new R2Storage instance
func NewR2Storage(endpoint, accessKeyID, secretAccessKey, bucket, publicURL, baseURL string) (*R2Storage, error) {
	// Create custom resolver for R2 endpoint
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:           endpoint,
			SigningRegion: "auto", // R2 uses "auto" as region
		}, nil
	})

	// Load AWS config with custom endpoint
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create S3 client
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true // R2 requires path-style addressing
	})

	return &R2Storage{
		client:    client,
		bucket:    bucket,
		publicURL: publicURL,
		baseURL:   baseURL,
	}, nil
}

// UploadImage uploads and compresses an image file to R2
func (s *R2Storage) UploadImage(file *multipart.FileHeader) (string, error) {
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

	key := s.objectKey(filename)

	// Encode image to buffer
	var buf bytes.Buffer
	if err := encodeImage(&buf, img, format); err != nil {
		return "", fmt.Errorf("failed to encode image: %w", err)
	}

	// Determine content type
	contentType = fmt.Sprintf("image/%s", format)
	if format == "jpg" {
		contentType = "image/jpeg"
	}

	// Upload to R2
	_, err = s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(buf.Bytes()),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to R2: %w", err)
	}

	// Return public URL
	return s.GetFileURL(filename), nil
}

// DeleteFile deletes a file from R2 storage
func (s *R2Storage) DeleteFile(filename string) error {
	key := s.objectKey(filename)

	_, err := s.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	return err
}

func (s *R2Storage) OpenFile(filename string) (io.ReadCloser, string, error) {
	key := s.objectKey(filename)

	result, err := s.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, "", err
	}

	contentType := contentTypeFromFilename(key)
	if result.ContentType != nil && *result.ContentType != "" {
		contentType = *result.ContentType
	}

	return result.Body, contentType, nil
}

// GetFileURL returns the public URL for a file
func (s *R2Storage) GetFileURL(filename string) string {
	// If publicURL is set, use it
	if s.publicURL != "" {
		// Remove trailing slash from publicURL
		url := strings.TrimSuffix(s.publicURL, "/")

		// Add baseURL prefix if set (e.g., "uploads")
		if s.baseURL != "" {
			prefix := strings.TrimPrefix(s.baseURL, "/")
			return fmt.Sprintf("%s/%s/%s", url, prefix, filename)
		}

		// Just use filename if no baseURL
		return fmt.Sprintf("%s/%s", url, filename)
	}

	// Fallback to baseURL (for backward compatibility)
	return fmt.Sprintf("%s/%s", s.baseURL, filename)
}

func (s *R2Storage) objectKey(filename string) string {
	key := strings.TrimSpace(filename)
	if parsed, err := url.Parse(key); err == nil && parsed.Scheme != "" && parsed.Host != "" {
		key = parsed.Path
	}

	key = strings.TrimPrefix(key, "/")
	key = strings.TrimPrefix(key, strings.TrimPrefix(s.publicURL, "/"))
	key = strings.TrimPrefix(key, "/")

	if s.baseURL == "" {
		return key
	}

	prefix := strings.Trim(strings.TrimPrefix(s.baseURL, "/"), "/")
	if prefix == "" || key == prefix || strings.HasPrefix(key, prefix+"/") {
		return key
	}

	return fmt.Sprintf("%s/%s", prefix, key)
}
