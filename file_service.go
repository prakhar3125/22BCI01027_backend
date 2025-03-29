package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"sync"
)

// FileService handles file operations
type FileService struct {
	fileRepo *FileRepository
	mutex    sync.Mutex
}

func NewFileService(fileRepo *FileRepository) *FileService {
	return &FileService{
		fileRepo: fileRepo,
		mutex:    sync.Mutex{},
	}
}

// UploadFile uploads a file to local storage and saves metadata to database
func (s *FileService) UploadFile(userID int, fileHeader *multipart.FileHeader) (*File, error) {
	// Open the uploaded file
	src, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	// Generate a unique filename
	uniqueFilename, err := generateUniqueFilename(fileHeader.Filename)
	if err != nil {
		return nil, err
	}

	// Create uploads directory if it doesn't exist
	uploadsDir := "./uploads"
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		return nil, err
	}

	// Create file path
	filePath := filepath.Join(uploadsDir, uniqueFilename)

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	// Copy file to destination
	if _, err = io.Copy(dst, src); err != nil {
		return nil, err
	}

	// Create file metadata
	file := &File{
		UserID:          userID,
		Filename:        uniqueFilename,
		OriginalFilename: fileHeader.Filename,
		FilePath:        filePath,
		FileSize:        fileHeader.Size,
		MimeType:        fileHeader.Header.Get("Content-Type"),
		IsPublic:        false,
	}

	// Save file metadata to database
	fileID, err := s.fileRepo.Create(file)
	if err != nil {
		// Delete the file if metadata saving fails
		os.Remove(filePath)
		return nil, err
	}

	file.ID = fileID
	return file, nil
}

// UploadFileAsync uploads multiple files concurrently
func (s *FileService) UploadFilesAsync(userID int, fileHeaders []*multipart.FileHeader) ([]int, error) {
	var wg sync.WaitGroup
	fileIDs := make([]int, len(fileHeaders))
	errorsChan := make(chan error, len(fileHeaders))

	// Process each file in its own goroutine
	for i, fileHeader := range fileHeaders {
		wg.Add(1)
		go func(i int, fileHeader *multipart.FileHeader) {
			defer wg.Done()

			file, err := s.UploadFile(userID, fileHeader)
			if err != nil {
				errorsChan <- err
				return
			}
			fileIDs[i] = file.ID
		}(i, fileHeader)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	close(errorsChan)

	// Check for errors
	if len(errorsChan) > 0 {
		return nil, <-errorsChan
	}

	return fileIDs, nil
}

// GetUserFiles retrieves all files for a user
func (s *FileService) GetUserFiles(userID int) ([]*File, error) {
	return s.fileRepo.GetByUserID(userID)
}

// SearchFiles searches for files by name
func (s *FileService) SearchFiles(userID int, name string) ([]*File, error) {
	return s.fileRepo.SearchByName(userID, name)
}

// GetFile retrieves a file by ID
func (s *FileService) GetFile(fileID, userID int) (*File, error) {
	file, err := s.fileRepo.GetByID(fileID)
	if err != nil {
		return nil, err
	}

	// Check if the file belongs to the user or is public
	if file.UserID != userID && !file.IsPublic {
		return nil, errors.New("file not found or you don't have permission to access it")
	}

	return file, nil
}

// ShareFile makes a file publicly accessible
func (s *FileService) ShareFile(fileID, userID int) (string, error) {
	// Check if the file exists and belongs to the user
	file, err := s.fileRepo.GetByID(fileID)
	if err != nil {
		return "", err
	}

	if file.UserID != userID {
		return "", errors.New("file not found or you don't have permission to share it")
	}

	// Make the file public
	if !file.IsPublic {
		if err := s.fileRepo.UpdatePublicStatus(fileID, userID, true); err != nil {
			return "", err
		}
	}

	// Return the public URL
	return fmt.Sprintf("/files/%d", fileID), nil
}

// DeleteFile deletes a file and its metadata
func (s *FileService) DeleteFile(fileID, userID int) error {
	// Check if the file exists and belongs to the user
	file, err := s.fileRepo.GetByID(fileID)
	if err != nil {
		return err
	}

	if file.UserID != userID {
		return errors.New("file not found or you don't have permission to delete it")
	}

	// Lock to prevent concurrent access to the file
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Delete the file from local storage
	if err := os.Remove(file.FilePath); err != nil && !os.IsNotExist(err) {
		return err
	}

	// Delete the file metadata from database
	return s.fileRepo.Delete(fileID, userID)
}

// generateUniqueFilename generates a unique filename
func generateUniqueFilename(originalFilename string) (string, error) {
	// Generate random bytes
	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	// Get file extension
	ext := filepath.Ext(originalFilename)

	// Create unique filename
	return hex.EncodeToString(randomBytes) + ext, nil
}