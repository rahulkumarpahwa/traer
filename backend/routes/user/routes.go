package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/rahulkumarpahwa/traer/storage"
	"github.com/rahulkumarpahwa/traer/utils"
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
