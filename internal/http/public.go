package http

import (
	"concert/internal/concert"
	"concert/internal/utils"
	"context"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
)

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	shows, err := h.Service.ListAllShow()
	if err != nil {
		log.Printf("Error listing shows for home page: %v", err)
		shows = []concert.Show{}
	}

	data := struct {
		ShowInfo
		User *concert.User
	}{
		ShowInfo: ShowInfo{Shows: shows},
		User:     h.getCurrentUser(r),
	}

	tmpl := template.Must(template.ParseFiles(
		utils.GetTemplatePath("base.html"),
		utils.GetTemplatePath("home.html"),
	))
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error executing home template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	username := r.FormValue("username")
	password := r.FormValue("password")
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	role := r.FormValue("role")

	// because this value are required
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

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	emailOrUsername := r.FormValue("email")
	password := r.FormValue("password")

	if emailOrUsername == "" || password == "" {
		http.Error(w, "Email or username and password are required", http.StatusBadRequest)
		return
	}

	user, err := AuthenticateUser(h.Db, emailOrUsername, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	log.Printf("User logged in successfully: email=%s, username=%s", user.Email, user.Username)

	SetCookie(w, user.ID)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	// Clear session cookie
	DeleteCookie(w)

	// Redirect to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) GetLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	tmpl := template.Must(template.ParseFiles(
		utils.GetTemplatePath("base.html"),
		utils.GetTemplatePath("login.html"),
	))
	if err := tmpl.Execute(w, nil); err != nil {
		log.Printf("Error executing login template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	tmpl := template.Must(template.ParseFiles(
		utils.GetTemplatePath("base.html"),
		utils.GetTemplatePath("register.html"),
	))
	if err := tmpl.Execute(w, nil); err != nil {
		log.Printf("Error executing register template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	tmpl := template.Must(template.ParseFiles(
		utils.GetTemplatePath("base.html"),
		utils.GetTemplatePath("error.html"),
	))
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "404 - Page Not Found", http.StatusNotFound)
		return
	}
}

func (h *Handler) ForgetPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/forget-password", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}
	email := r.FormValue("email")
	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}
	user, err := FindUserByEmailOrUsername(h.Db, email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)

		return
	}
	log.Printf("User found: email=%s, username=%s", user.Email, user.Username)

	// Send reset password email via Temporal workflow
	if err := h.SendResetPasswordEmail(user.Email); err != nil {
		log.Printf("Error sending reset password email: %v", err)
		http.Error(w, "Failed to send reset password email", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) GetForgetPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	tmpl := template.Must(template.ParseFiles(
		utils.GetTemplatePath("base.html"),
		utils.GetTemplatePath("forget-password.html"),
	))
	if err := tmpl.Execute(w, nil); err != nil {
		log.Printf("Error executing forget password template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
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
