package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	fmt.Println(r)
	email := r.FormValue("email")
	username := r.FormValue("username")
	password := r.FormValue("password")
	firstName := r.FormValue("FirstName")
	lastName := r.FormValue("LastName")
	role := r.FormValue("role")

	if email == "" || username == "" || password == "" {
		http.Error(w, "Email, username, and password are required", http.StatusBadRequest)
		return
	}

	exists, err := UserExists(h.Db, email, username)
	if err != nil {
		log.Printf("Error checking existing user: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "Email or username already exists", http.StatusConflict)
		return
	}

	user, err := InsertUser(h.Db, email, username, password, firstName, lastName, role)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	log.Printf("User registered successfully: email=%s, username=%s", email, username)

	SetCookie(w, user.ID)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

type RegisterRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type UserResponse struct {
	ID        uint   `json:"id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Role      string `json:"role"`
}

func (h *Handler) RegisterAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode JSON: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("RegisterAPI received: %+v", req)

	if req.Email == "" || req.Username == "" || req.Password == "" {
		http.Error(w, "Email, username, and password are required", http.StatusBadRequest)
		return
	}

	exists, err := UserExists(h.Db, req.Email, req.Username)
	if err != nil {
		log.Printf("Error checking user: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if exists {
		http.Error(w, "Email or username already exists", http.StatusConflict)
		return
	}

	var role string

	switch {
	case strings.HasPrefix(req.Username, "admin-"):
		role = "admin"
	case strings.HasPrefix(req.Username, "mod-"):
	default:
		role = "user"
	}

	user, err := InsertUser(
		h.Db,
		req.Email,
		req.Username,
		req.Password,
		req.FirstName,
		req.LastName,
		role,
	)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	SetCookie(w, user.ID)

	response := UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) LoginAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode JSON: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	user, err := AuthenticateUser(h.Db, req.Email, req.Password)
	if err != nil {
		log.Printf("Authentication failed: %v", err)
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	SetCookie(w, user.ID)
	log.Printf("User logged in: %s", user.Username)

	response := UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *Handler) ForgetPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/forget-password", http.StatusSeeOther)
		return
	}
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode JSON: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	email := req.Email

	user, err := FindUserByEmailOrUsername(h.Db, email)
	if err != nil {
		log.Printf("Error finding user: %v", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	log.Printf("User found: email=%s, username=%s", user.Email, user.Username)

	if err := h.SendResetPasswordEmail(email); err != nil {
		log.Printf("Error sending reset password email: %v", err)
		http.Error(w, "Failed to send reset password email", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Password reset email sent successfully",
	})
}

// execute a temporal workflow to send a reset password email
func (h *Handler) SendResetPasswordEmail(email string) error {
	ctx := context.Background()
	// generate a password
	// not secure

	generatePass := uuid.NewString()[:8]
	if err := UpdatePassword(h.Db, email, generatePass); err != nil {
		return err
	}
	workflowID := "email-workflow-" + time.Now().Format("20060102-150405")
	workflowOptions := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: "email-task-queue",
	}

	user := TemporalUserInput{
		Email:    email,
		Password: generatePass,
	}
	_, err := h.TemporalClient.ExecuteWorkflow(ctx, workflowOptions, "SendMailWorkflow", user)
	if err != nil {
		log.Printf("Error executing workflow: %v", err)
		return err
	}

	return nil
}
