package main

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Password is not included in JSON responses
	CreatedAt time.Time `json:"created_at"`
}

// File represents a file stored in the system
type File struct {
	ID              int       `json:"id"`
	UserID          int       `json:"user_id"`
	Filename        string    `json:"filename"`         // System-generated unique filename
	OriginalFilename string    `json:"original_filename"` // Original file name
	FilePath        string    `json:"file_path"`
	FileSize        int64     `json:"file_size"`
	MimeType        string    `json:"mime_type"`
	IsPublic        bool      `json:"is_public"`
	CreatedAt       time.Time `json:"created_at"`
}

// repositories.go
