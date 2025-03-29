package main


import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// InitDB initializes the database connection
func InitDB() (*sql.DB, error) {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	if dbHost == "" {
		dbHost = "localhost"
	}
	if dbPort == "" {
		dbPort = "3306"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", 
		dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	// Create tables if they don't exist
	if err = createTables(db); err != nil {
		return nil, err
	}

	return db, nil
}

// createTables creates the necessary tables if they don't exist
func createTables(db *sql.DB) error {
	// Create users table
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		email VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`)
	if err != nil {
		return err
	}

	// Create files table
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS files (
		id INT AUTO_INCREMENT PRIMARY KEY,
		user_id INT NOT NULL,
		filename VARCHAR(255) NOT NULL,
		original_filename VARCHAR(255) NOT NULL,
		file_path VARCHAR(255) NOT NULL,
		file_size BIGINT NOT NULL,
		mime_type VARCHAR(100),
		is_public BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);`)
	if err != nil {
		return err
	}

	return nil
}

