 
package models

import (
	"time"
)

type File struct {
	ID              int       `json:"id"`
	UserID          int       `json:"user_id"`
	Filename        string    `json:"filename"`
	OriginalFilename string   `json:"original_filename"`
	FilePath        string    `json:"-"` // Hidden from API responses
	FileSize        int64     `json:"file_size"`
	FileType        string    `json:"file_type"`
	IsPublic        bool      `json:"is_public"`
	ShareToken      string    `json:"share_token,omitempty"`
	ShareExpiry     time.Time `json:"share_expiry,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	URL             string    `json:"url,omitempty"` // Computed field
}