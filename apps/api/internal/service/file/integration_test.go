package file

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"mime/multipart"
	"net/textproto"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helper to create a temporary image file in memory
func createTestImage(t *testing.T) *multipart.FileHeader {
	// Create a simple image
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	// Set some color
	for x := 0; x < 100; x++ {
		for y := 0; y < 100; y++ {
			img.Set(x, y, color.RGBA{255, 0, 0, 255})
		}
	}

	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	assert.NoError(t, err)

	// Create multipart file header
	// We need to write to a temp file first because FileHeader.Open() reads from disk usually if created via FormFile,
	// BUT here we can manually construct it if we had a way to provide content.
	// `multipart.FileHeader` usually comes from `ParseMultipartForm`.
	// For testing `UploadImage(file *multipart.FileHeader)`, `file.Open()` is called.
	// `Open()` method of FileHeader opens the file.
	// We have to mock the behavior or create a real request and parse it.
	// Easiest is to write bytes to a temp file and simulate FileHeader pointing to it?
	// Actually, easier to use a buffer and standard multipart writer to create a request body, 
	// then parse it to get a real FileHeader.

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	
	// manually create part with header
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="file"; filename="test.png"`)
	h.Set("Content-Type", "image/png")
	part, err := writer.CreatePart(h)
	assert.NoError(t, err)
	_, err = part.Write(buf.Bytes())
	assert.NoError(t, err)
	writer.Close()

	// Create a dummy HTTP request to parse the body
	r := multipart.NewReader(body, writer.Boundary())
	form, err := r.ReadForm(1024 * 1024)
	assert.NoError(t, err)
	
	return form.File["file"][0]
}

func TestLocalStorage_Integration_UploadImage(t *testing.T) {
	// Setup Temp Dir
	tempDir, err := os.MkdirTemp("", "file-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir) // Cleanup

	baseURL := "http://localhost:8080/files"
	storage := NewLocalStorage(tempDir, baseURL)

	// Case 1: Success Upload
	fileHeader := createTestImage(t)
	url, err := storage.UploadImage(fileHeader)
	
	assert.NoError(t, err)
	assert.Contains(t, url, baseURL)
	assert.Contains(t, url, ".png") // Or jpeg if converted
	
	// Verify file exists on disk
	filename := filepath.Base(url)
	_, err = os.Stat(filepath.Join(tempDir, filename))
	assert.NoError(t, err)

	// Case 2: Validation Failure (Invalid File Type)
	// Create text file
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.txt")
	assert.NoError(t, err)
	part.Write([]byte("not an image"))
	writer.Close()
	
	r := multipart.NewReader(body, writer.Boundary())
	_, err = r.ReadForm(1024)
	assert.NoError(t, err)
	// txtFile := form.File["file"][0] // Unused, we construct proper one below
	// Manually set content type header if needed, but ReadForm usually detects or sets from part header
	// The part header we created: `CreateFormFile` sets Content-Disposition. Content-Type is optional.
	// `LocalStorage` checks `file.Header.Get("Content-Type")`.
	// NewWriter doesn't automatically set Content-Type for the part unless we use CreatePart.
	// Let's create proper headers.
	
	body2 := &bytes.Buffer{}
	writer2 := multipart.NewWriter(body2)
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="file"; filename="test.txt"`)
	h.Set("Content-Type", "text/plain")
	part2, _ := writer2.CreatePart(h)
	part2.Write([]byte("not an image"))
	writer2.Close()
	
	r2 := multipart.NewReader(body2, writer2.Boundary())
	form2, _ := r2.ReadForm(1024)
	txtFile2 := form2.File["file"][0]

	_, err = storage.UploadImage(txtFile2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "file is not an image")
}
