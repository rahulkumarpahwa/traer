package storage

import (
	"database/sql"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type UserStorage struct {
	DB *sql.DB
}

type UserStorageInterface interface {

}

func CreateUser(db *sql.DB, username, email, password string) error {
	// Validate inputs
	if username == "" || email == "" || password == "" {
		return fmt.Errorf("username, email, and password cannot be empty")
	}

	// Hash the password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	// Prepare the SQL statement
	stmt, err := db.Prepare(`
		INSERT INTO users (username, email, password_hash)
		VALUES (?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close() // Always close prepared statements

	// Execute the prepared statement
	_, err = stmt.Exec(username, email, string(hashedPassword))
	if err != nil {
		return fmt.Errorf("failed to execute statement: %v", err)
	}

	return nil
}