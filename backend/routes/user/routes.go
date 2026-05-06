package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/rahulkumarpahwa/traer/storage"
	"github.com/rahulkumarpahwa/traer/types"
	"github.com/rahulkumarpahwa/traer/utils"
	"golang.org/x/crypto/bcrypt"
)

// UserHandler handles HTTP requests for user operations
type UserHandler struct {
	UserStorage storage.UserStorageInterface
}

// CreateUserRequest represents the JSON body for creating a user
type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UpdateUserRequest represents the JSON body for updating a user
type UpdateUserRequest struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// LoginRequest represents the JSON body for login
type LoginRequest struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password"`
}

// HandleLogin handles POST /login
func (u *UserHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ErrorResponse(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid JSON body"))
		return
	}

	// Validate input: must have password and either username or email
	if req.Password == "" || (req.Username == "" && req.Email == "") {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("must provide either username+password or email+password"))
		return
	}

	// Fetch user from storage
	var userObj interface{}
	var err error

	if req.Username != "" {
		userObj, err = u.UserStorage.GetUserByUsername(req.Username)
	} else {
		userObj, err = u.UserStorage.GetUserByEmail(req.Email)
	}

	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	if userObj == nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, fmt.Errorf("invalid credentials"))
		return
	}

	// Type assert to *types.User
	user, ok := userObj.(*types.User) // Replace with actual type, e.g., *types.User
	if !ok {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("internal type assertion error"))
		return
	}

	// Compare password hash
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, fmt.Errorf("invalid credentials"))
		return
	}


	// Token Generation


	utils.JSONResponse(w, http.StatusOK, map[string]string{
		"message": "login successful",
		"user_id": fmt.Sprintf("%d", user.ID),
	})
}

// HandleCreateUser handles POST /users
func (hs *UserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ErrorResponse(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid JSON body"))
		return
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("username, email, and password are required"))
		return
	}

	if err := hs.UserStorage.CreateUser(req.Username, req.Email, req.Password); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	utils.JSONResponse(w, http.StatusCreated, map[string]string{"message": "User created successfully"})
}

// HandleUpdateUser handles PUT /users?id=123
func (hs *UserHandler) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.ErrorResponse(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid user ID"))
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid JSON body"))
		return
	}

	if req.Username == "" && req.Password == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("at least one field (username or password) must be provided"))
		return
	}

	if err := hs.UserStorage.UpdateUserByID(id, req.Username, req.Password); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	utils.JSONResponse(w, http.StatusOK, map[string]string{"message": "User updated successfully"})
}

// HandleDeleteUser handles DELETE /users?id=123
func (hs *UserHandler) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.ErrorResponse(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid user ID"))
		return
	}

	if err := hs.UserStorage.DeleteUserByID(id); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	utils.JSONResponse(w, http.StatusOK, map[string]string{"message": "User deleted successfully"})
}

// HandleGetUserByID handles GET /users/id?id=123
func (hs *UserHandler) HandleGetUserByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.ErrorResponse(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid user ID"))
		return
	}

	user, err := hs.UserStorage.GetUserByID(id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	if user == nil {
		utils.ErrorResponse(w, http.StatusNotFound, fmt.Errorf("user not found"))
		return
	}

	utils.JSONResponse(w, http.StatusOK, user)
}

// HandleGetUserByEmail handles GET /users/email?email=...
func (hs *UserHandler) HandleGetUserByEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.ErrorResponse(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}

	email := r.URL.Query().Get("email")
	if email == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("email is required"))
		return
	}

	user, err := hs.UserStorage.GetUserByEmail(email)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	if user == nil {
		utils.ErrorResponse(w, http.StatusNotFound, fmt.Errorf("user not found"))
		return
	}

	utils.JSONResponse(w, http.StatusOK, user)
}

// HandleGetUserByUsername handles GET /users/username?username=...
func (hs *UserHandler) HandleGetUserByUsername(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.ErrorResponse(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("username is required"))
		return
	}

	user, err := hs.UserStorage.GetUserByUsername(username)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	if user == nil {
		utils.ErrorResponse(w, http.StatusNotFound, fmt.Errorf("user not found"))
		return
	}

	utils.JSONResponse(w, http.StatusOK, user)
}
