package file

import (
	"mime/multipart"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStorageProvider
type MockStorageProvider struct {
	mock.Mock
}

func (m *MockStorageProvider) UploadImage(file *multipart.FileHeader) (string, error) {
	args := m.Called(file)
	return args.String(0), args.Error(1)
}

func (m *MockStorageProvider) DeleteFile(filename string) error {
	args := m.Called(filename)
	return args.Error(0)
}

func (m *MockStorageProvider) GetFileURL(filename string) string {
	args := m.Called(filename)
	return args.String(0)
}

func TestFileService_UploadImage(t *testing.T) {
	mockStorage := new(MockStorageProvider)
	service := NewService(mockStorage)

	fileHeader := &multipart.FileHeader{Filename: "test.jpg"}
	
	mockStorage.On("UploadImage", fileHeader).Return("uploads/test.jpg", nil)

	url, err := service.UploadImage(fileHeader)
	
	assert.NoError(t, err)
	assert.Equal(t, "uploads/test.jpg", url)
	mockStorage.AssertExpectations(t)
}

func TestFileService_DeleteFile(t *testing.T) {
	mockStorage := new(MockStorageProvider)
	service := NewService(mockStorage)

	mockStorage.On("DeleteFile", "uploads/test.jpg").Return(nil)

	err := service.DeleteFile("uploads/test.jpg")
	
	assert.NoError(t, err)
	mockStorage.AssertExpectations(t)
}
