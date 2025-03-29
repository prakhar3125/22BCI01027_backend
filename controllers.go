package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AuthController handles authentication-related requests
type AuthController struct {
	authService *AuthService
}

func NewAuthController(authService *AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// Register handles user registration
func (c *AuthController) Register(ctx *gin.Context) {
	var request struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := c.authService.Register(request.Email, request.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": id, "message": "User registered successfully"})
}

// Login handles user login
func (c *AuthController) Login(ctx *gin.Context) {
	var request struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := c.authService.Login(request.Email, request.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

// FileController handles file-related requests
type FileController struct {
	fileService *FileService
}

func NewFileController(fileService *FileService) *FileController {
	return &FileController{fileService: fileService}
}

// UploadFile handles file upload
func (c *FileController) UploadFile(ctx *gin.Context) {
	// Get user ID from context
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get file
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	// Upload file
	uploadedFile, err := c.fileService.UploadFile(userID.(int), file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return file metadata
	ctx.JSON(http.StatusCreated, gin.H{
		"id": uploadedFile.ID,
		"filename": uploadedFile.OriginalFilename,
		"size": uploadedFile.FileSize,
		"url": "/files/" + strconv.Itoa(uploadedFile.ID),
	})
}

// GetUserFiles handles retrieval of user files
func (c *FileController) GetUserFiles(ctx *gin.Context) {
	// Get user ID from context
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Check if search query is provided
	searchQuery := ctx.Query("search")
	var files []*File
	var err error

	if searchQuery != "" {
		// Search files by name
		files, err = c.fileService.SearchFiles(userID.(int), searchQuery)
	} else {
		// Get all user files
		files, err = c.fileService.GetUserFiles(userID.(int))
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return files
	ctx.JSON(http.StatusOK, gin.H{"files": files})
}

// GetFile handles file retrieval
func (c *FileController) GetFile(ctx *gin.Context) {
	// Get file ID from URL
	fileID, err := strconv.Atoi(ctx.Param("file_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	// Get user ID from context
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get file
	file, err := c.fileService.GetFile(fileID, userID.(int))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Return file
	ctx.File(file.FilePath)
}

// ShareFile handles file sharing
func (c *FileController) ShareFile(ctx *gin.Context) {
	// Get file ID from URL
	fileID, err := strconv.Atoi(ctx.Param("file_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	// Get user ID from context
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Share file
	url, err := c.fileService.ShareFile(fileID, userID.(int))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Return file URL
	ctx.JSON(http.StatusOK, gin.H{"url": url})
}



// controllers.go (continued)

// DeleteFile handles file deletion
func (c *FileController) DeleteFile(ctx *gin.Context) {
	// Get file ID from URL
	fileID, err := strconv.Atoi(ctx.Param("file_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	// Get user ID from context
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Delete file
	err = c.fileService.DeleteFile(fileID, userID.(int))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Return success message
	ctx.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}