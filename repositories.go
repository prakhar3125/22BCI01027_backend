package main

import (
	"database/sql"
	"errors"
)

// UserRepository handles database operations for users
type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(email, hashedPassword string) (int, error) {
	query := "INSERT INTO users (email, password) VALUES (?, ?)"
	result, err := r.db.Exec(query, email, hashedPassword)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (r *UserRepository) GetByEmail(email string) (*User, error) {
	query := "SELECT id, email, password, created_at FROM users WHERE email = ?"
	row := r.db.QueryRow(query, email)

	var user User
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByID(id int) (*User, error) {
	query := "SELECT id, email, password, created_at FROM users WHERE id = ?"
	row := r.db.QueryRow(query, id)

	var user User
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// FileRepository handles database operations for files
type FileRepository struct {
	db *sql.DB
}

func NewFileRepository(db *sql.DB) *FileRepository {
	return &FileRepository{db: db}
}

func (r *FileRepository) Create(file *File) (int, error) {
	query := `
		INSERT INTO files (user_id, filename, original_filename, file_path, file_size, mime_type, is_public)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	result, err := r.db.Exec(
		query, 
		file.UserID, 
		file.Filename,
		file.OriginalFilename,
		file.FilePath,
		file.FileSize,
		file.MimeType,
		file.IsPublic,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (r *FileRepository) GetByID(id int) (*File, error) {
	query := `
		SELECT id, user_id, filename, original_filename, file_path, file_size, mime_type, is_public, created_at
		FROM files
		WHERE id = ?
	`
	row := r.db.QueryRow(query, id)

	var file File
	err := row.Scan(
		&file.ID,
		&file.UserID,
		&file.Filename,
		&file.OriginalFilename,
		&file.FilePath,
		&file.FileSize,
		&file.MimeType,
		&file.IsPublic,
		&file.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("file not found")
		}
		return nil, err
	}

	return &file, nil
}

func (r *FileRepository) GetByUserID(userID int) ([]*File, error) {
	query := `
		SELECT id, user_id, filename, original_filename, file_path, file_size, mime_type, is_public, created_at
		FROM files
		WHERE user_id = ?
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []*File
	for rows.Next() {
		var file File
		err := rows.Scan(
			&file.ID,
			&file.UserID,
			&file.Filename,
			&file.OriginalFilename,
			&file.FilePath,
			&file.FileSize,
			&file.MimeType,
			&file.IsPublic,
			&file.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		files = append(files, &file)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return files, nil
}

func (r *FileRepository) SearchByName(userID int, name string) ([]*File, error) {
	query := `
		SELECT id, user_id, filename, original_filename, file_path, file_size, mime_type, is_public, created_at
		FROM files
		WHERE user_id = ? AND original_filename LIKE ?
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, userID, "%"+name+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []*File
	for rows.Next() {
		var file File
		err := rows.Scan(
			&file.ID,
			&file.UserID,
			&file.Filename,
			&file.OriginalFilename,
			&file.FilePath,
			&file.FileSize,
			&file.MimeType,
			&file.IsPublic,
			&file.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		files = append(files, &file)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return files, nil
}

func (r *FileRepository) Delete(id int, userID int) error {
	query := "DELETE FROM files WHERE id = ? AND user_id = ?"
	result, err := r.db.Exec(query, id, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("file not found or you don't have permission to delete it")
	}

	return nil
}

func (r *FileRepository) UpdatePublicStatus(id int, userID int, isPublic bool) error {
	query := "UPDATE files SET is_public = ? WHERE id = ? AND user_id = ?"
	result, err := r.db.Exec(query, isPublic, id, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("file not found or you don't have permission to update it")
	}

	return nil
}
