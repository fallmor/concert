package http

import (
	"concert/internal/concert"
	"concert/internal/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/go-chi/chi/v5"
	"go.temporal.io/sdk/client"
	"gorm.io/gorm"
)

func NewRouter(service concert.ConcertService, db *gorm.DB) (*Handler, error) {
	tmpClient, err := client.Dial(
		client.Options{HostPort: "127.0.0.1:7233"})
	if err != nil {
		return nil, err
	}
	return &Handler{
		Service:        service,
		Db:             db,
		TemporalClient: tmpClient,
	}, nil
}

func (h *Handler) Close() {
	if h.TemporalClient != nil {
		h.TemporalClient.Close()
	}
}

func (h *Handler) ChiSetRoutes() {
	h.Route = chi.NewRouter()

	staticDir := utils.GetStaticDir()
	fileServer := http.FileServer(http.Dir(staticDir))
	h.Route.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	h.Route.Group(func(r chi.Router) {
		r.Use(RateLimit)
		r.Get("/", h.Home)
		r.Get("/login", h.GetLogin)
		r.Post("/login", h.Login)
		r.Get("/register", h.GetRegister)
		r.Post("/register", h.Register)
		r.Get("/forget-password", h.GetForgetPassword)
		r.Post("/forget-password", h.ForgetPassword)
		r.NotFound(h.NotFoundHandler)
	})

	h.Route.Group(func(r chi.Router) {
		r.Use(NeedsAuth(h.Db))

		r.Get("/fan/{name}", h.GetFan)
		r.Get("/show/{artistname}", h.GetShow)
		r.Get("/fans/new", h.NewFan)
		r.Get("/shows", h.ListAllShow)
		r.Get("/artists", h.ListAllArtists)
		r.Get("/list", h.ListAllFan)
		r.Post("/submit", h.ParticipateShow)
		r.Post("/upload/image", h.UploadImage)
		r.Post("/logout", h.Logout)
		r.Get("/profile", h.GetProfile)
	})

	h.Route.Group(func(r chi.Router) {
		r.Use(NeedsAuth(h.Db))
		r.Use(NeedsRole(h.Db, "moderator"))
		h.ModeratorRoutes(r)
	})

	// Admin routes (
	h.Route.Group(func(r chi.Router) {
		r.Use(NeedsAuth(h.Db))
		r.Use(NeedsRole(h.Db, "admin"))
		h.AdminRoutes(r)
		// Also allow moderators' edit routes for admins
		r.Get("/shows/{id}/edit", h.EditShow)
		r.Get("/artists/{id}/edit", h.EditArtist)
		r.Get("/fans/{id}/edit", h.EditFan)
		r.Post("/shows/{id}/update", h.UpdateShow)
		r.Post("/artists/{id}/update", h.UpdateArtist)
		r.Post("/fans/{id}/update", h.UpdateFan)
	})
}

func (h *Handler) saveUploadedFile(file multipart.File, handler *multipart.FileHeader, folder string) (string, error) {
	defer file.Close()

	projectRoot := utils.GetProjectRoot("go.mod")
	uploadDir := filepath.Join(projectRoot, "static", "images", folder)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	ext := filepath.Ext(handler.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return "", fmt.Errorf("invalid file extension: %s", ext)
	}
	// for security reasons we need to sanitize the filename
	baseName := strings.TrimSuffix(handler.Filename, ext)
	baseName = strings.ReplaceAll(baseName, " ", "_")
	baseName = strings.ReplaceAll(baseName, "/", "_")
	baseName = strings.ReplaceAll(baseName, "\\", "_")
	baseName = strings.ReplaceAll(baseName, "..", "_")
	filename := fmt.Sprintf("%d_%s%s", time.Now().Unix(), baseName, ext)
	filePath := filepath.Join(uploadDir, filename)

	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	imageURL := fmt.Sprintf("/static/images/%s/%s", folder, filename)
	log.Printf("Image file saved: %s -> %s (stored in database)", handler.Filename, imageURL)
	return imageURL, nil
}

func (h *Handler) UploadImage(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "File too large or invalid", http.StatusBadRequest)
		return
	}

	folder := r.FormValue("folder")
	if folder != "artists" && folder != "albums" {
		folder = "artists"
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "No file uploaded", http.StatusBadRequest)
		return
	}
	defer file.Close()

	url, err := h.saveUploadedFile(file, handler, folder)
	if err != nil {
		log.Printf("Error uploading image: %v", err)
		http.Error(w, "Failed to upload image", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"url": url})
}

func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	user, err := GetUserFromCookie(h.Db, r)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if err := h.Db.Preload("Fans.Show.Artist").First(user, user.ID).Error; err != nil {
		log.Printf("Error loading user profile: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	//  this struct satisfies both templates:
	// base.html and profile.html
	data := struct {
		User      *concert.User
		ID        uint
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt gorm.DeletedAt
		Email     string
		Username  string
		FirstName string
		LastName  string
		Role      string
		Fans      []concert.Fan
	}{
		User:      user,
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt,
		Email:     user.Email,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		Fans:      user.Fans,
	}

	tmpl := template.Must(template.ParseFiles(
		utils.GetTemplatePath("base.html"),
		utils.GetTemplatePath("profile.html"),
	))
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error executing profile template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
