package storage

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/rahulkumarpahwa/traer/internal/types"
	"golang.org/x/crypto/bcrypt"
)

type UserStorage struct {
	DB *sql.DB
}

type UserStorageInterface interface {
	CreateUser(email, password string) error
	UpdateUserByID(id int, password string) error
	DeleteUserByID(id int) error
	GetUserByID(id int) (*types.User, error)
	GetUserByEmail(email string) (*types.User, error)
}

func (u *UserStorage) CreateUser(email, password string) error {
	// Validate inputs
	if email == "" || password == "" {
		return fmt.Errorf("email, and password cannot be empty")
	}

	// Hash the password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	// Prepare the SQL statement
	stmt, err := u.DB.Prepare(`
		INSERT INTO users (email, password_hash)
		VALUES (?, ?)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close() // Always close prepared statements

	// Execute the prepared statement
	_, err = stmt.Exec(email, string(hashedPassword))
	if err != nil {
		return fmt.Errorf("failed to execute statement: %v", err)
	}

	return nil
}

func (u *UserStorage) UpdateUserByID(id int, password string) error {
	if id <= 0 {
		return fmt.Errorf("invalid user ID")
	}

	// Build query parts dynamically
	setClauses := []string{}
	args := []interface{}{}


	// If password is provided, hash it and add to the update
	if password != "" {
		hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %v", err)
		}
		setClauses = append(setClauses, "password_hash = ?")
		args = append(args, string(hashedBytes))
	}

	// If nothing to update, return error
	if len(setClauses) == 0 {
		return fmt.Errorf("no fields to update")
	}

	// Always update the timestamp
	setClauses = append(setClauses, "updated_at = CURRENT_TIMESTAMP")

	// Build final query
	query := fmt.Sprintf("UPDATE users SET %s WHERE id = ?",
		strings.Join(setClauses, ", "))
	args = append(args, id)

	// Prepare and execute
	stmt, err := u.DB.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare update: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(args...)
	if err != nil {
		return fmt.Errorf("failed to execute update: %v", err)
	}

	return nil
}

// DeleteUser removes a user by ID
func (u *UserStorage) DeleteUserByID(id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid user ID")
	}

	stmt, err := u.DB.Prepare(`DELETE FROM users WHERE id = ?`)
	if err != nil {
		return fmt.Errorf("failed to prepare delete: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("failed to execute delete: %v", err)
	}

	return nil
}

// GetUserByID retrieves a user by ID
func (u *UserStorage) GetUserByID(id int) (*types.User, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid user ID")
	}

	stmt, err := u.DB.Prepare(`SELECT id, email, password_hash, created_at, updated_at FROM users WHERE id = ?`)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare select: %v", err)
	}
	defer stmt.Close()

	var user types.User
	err = stmt.QueryRow(id).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil // No user found
	} else if err != nil {
		return nil, fmt.Errorf("failed to query user: %v", err)
	}

	return &user, nil
}

// GetUserByEmail retrieves a user by email
func (u *UserStorage) GetUserByEmail(email string) (*types.User, error) {
	if email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}

	stmt, err := u.DB.Prepare(`SELECT id, email, password_hash, created_at, updated_at FROM users WHERE email = ?`)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare select: %v", err)
	}
	defer stmt.Close()

	var user types.User
	err = stmt.QueryRow(email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to query user: %v", err)
	}

	return &user, nil
}