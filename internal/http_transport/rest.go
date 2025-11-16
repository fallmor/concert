package httptransport

import (
	"concert/internal/concert"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/go-chi/chi/v5"
)

func getProjectRoot() string {
	wd, err := os.Getwd()
	if err != nil {
		wd = "."
	}
	
	currentDir := wd
	for {
		goModPath := filepath.Join(currentDir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return currentDir
		}
		
		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			return wd
		}
		currentDir = parentDir
	}
}

func getTemplatePath(filename string) string {
	projectRoot := getProjectRoot()
	return filepath.Join(projectRoot, "internal", "templates", filename)
}

func getStaticDir() string {
	projectRoot := getProjectRoot()
	return filepath.Join(projectRoot, "static")
}

type Handler struct {
	Route   chi.Router
	Service *concert.Service
}

type PageData struct {
	Title string
	Fans  []concert.Fan
}

type ShowInfo struct {
	Shows []concert.Show
}

type FanInfo struct {
	Fans []concert.Fan
}

type ArtistInfo struct {
	Artists []concert.Artist
}

func NewRouter(service *concert.Service) *Handler {
	return &Handler{
		Service: service,
	}
}

func (h *Handler) ChiSetRoutes() {
	h.Route = chi.NewRouter()
	staticDir := getStaticDir()
	fileServer := http.FileServer(http.Dir(staticDir))
	h.Route.Handle("/static/*", http.StripPrefix("/static/", fileServer))
	h.Route.Get("/fan/{name}", h.GetFan)
	h.Route.Get("/show/{artistname}", h.GetShow)
	h.Route.Get("/fans/new", h.NewFan)
	h.Route.Get("/shows", h.ListAllShow)
	h.Route.Get("/shows/new", h.NewShow)
	h.Route.Get("/artists", h.ListAllArtists)
	h.Route.Get("/artist/new", h.NewArtist)
	h.Route.Get("/list", h.ListAllFan)
	h.Route.Post("/submit", h.ParticipateShow)
	h.Route.Post("/show", h.SetShow)
	h.Route.Post("/artist", h.SetArtist)
	h.Route.Post("/upload/image", h.UploadImage)
}

func (h *Handler) GetFan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	name := chi.URLParam(r, "name")
	fanDetails, err := h.Service.GetFan(name)
	if err != nil {
		log.Printf("Error getting fan %s: %v", name, err)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	tmpl := template.Must(template.ParseFiles(
		getTemplatePath("base.html"),
		getTemplatePath("fan.html"),
	))

	data := PageData{
		Title: fmt.Sprintf("Show Fan named %s", name),
		Fans:  fanDetails,
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error executing fan template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetShow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	artistName := chi.URLParam(r, "artistname")
	showDetails, err := h.Service.GetShow(artistName)
	if err != nil {
		log.Printf("Error getting shows for artist %s: %v", artistName, err)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	data := ShowInfo{
		Shows: showDetails,
	}

	tmpl := template.Must(template.ParseFiles(
		getTemplatePath("base.html"),
		getTemplatePath("show_by_artist.html"),
	))

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error executing show template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ListAllShow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	shows, err := h.Service.ListAllShow()
	if err != nil {
		log.Printf("Error listing shows: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := ShowInfo{
		Shows: shows,
	}

	tmpl := template.Must(template.ParseFiles(
		getTemplatePath("base.html"),
		getTemplatePath("allshows.html"),
	))

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error executing shows template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ListAllFan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	allFans, err := h.Service.ListAllFan()
	if err != nil {
		log.Printf("Error listing fans: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := FanInfo{
		Fans: allFans,
	}

	tmpl := template.Must(template.ParseFiles(
		getTemplatePath("base.html"),
		getTemplatePath("allfans.html"),
	))

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error executing fans template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ListAllArtists(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	artists, err := h.Service.ListAllArtists()
	if err != nil {
		log.Printf("Error listing artists: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := ArtistInfo{
		Artists: artists,
	}

	tmpl := template.Must(template.ParseFiles(
		getTemplatePath("base.html"),
		getTemplatePath("allartists.html"),
	))
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error executing artists template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) NewShow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	artists, err := h.Service.ListAllArtists()
	if err != nil {
		log.Printf("Error listing artists: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	
	data := ArtistInfo{
		Artists: artists,
	}

	tmpl := template.Must(template.ParseFiles(
		getTemplatePath("base.html"),
		getTemplatePath("newshow.html"),
	))
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error rendering show form: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) SetShow(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")

	if contentType == "application/json" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		var show concert.Show
		if err := json.NewDecoder(r.Body).Decode(&show); err != nil {
			http.Error(w, "failed to decode the body", http.StatusBadRequest)
			return
		}
		myShow, err := h.Service.SetShow(show)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(myShow); err != nil {
			log.Printf("Error encoding show response: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}
	
	artistID, err := strconv.ParseUint(r.FormValue("artist_id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid Artist ID", http.StatusBadRequest)
		return
	}
	
	dateStr := r.FormValue("date")
	if dateStr == "" {
		http.Error(w, "Date and time are required", http.StatusBadRequest)
		return
	}
	
	date, err := time.Parse("2006-01-02T15:04", dateStr)
	if err != nil {
		log.Printf("Error parsing date '%s': %v", dateStr, err)
		http.Error(w, "Invalid date or time format. Please select both date and time.", http.StatusBadRequest)
		return
	}
	
	place := r.FormValue("place")
	if place == "" {
		http.Error(w, "Place is required", http.StatusBadRequest)
		return
	}
	
	show := concert.Show{
		Date:     date,
		ArtistID: uint(artistID),
		Place:    place,
	}
	
	_, err = h.Service.SetShow(show)
	if err != nil {
		log.Printf("Error creating show: %v", err)
		http.Error(w, "Failed to create show: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/shows", http.StatusSeeOther)
}

func (h *Handler) NewArtist(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	tmpl := template.Must(template.ParseFiles(
		getTemplatePath("base.html"),
		getTemplatePath("newartist.html"),
	))
	if err := tmpl.Execute(w, nil); err != nil {
		log.Printf("Error rendering artist form: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) SetArtist(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")

	if contentType == "application/json" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		var artist concert.Artist
		if err := json.NewDecoder(r.Body).Decode(&artist); err != nil {
			http.Error(w, "failed to decode the body", http.StatusBadRequest)
			return
		}
		myArtist, err := h.Service.SetArtist(artist)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(myArtist); err != nil {
			log.Printf("Error encoding artist response: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}
	
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	nom := r.FormValue("nom")
	genre := r.FormValue("genre")
	photoURL := r.FormValue("photo_url")
	albumURL := r.FormValue("album_url")

	if nom == "" || genre == "" {
		http.Error(w, "Name and genre are required", http.StatusBadRequest)
		return
	}

	if photoURL == "" {
		if file, handler, err := r.FormFile("photo_file"); err == nil {
			defer file.Close()
			savedURL, err := h.saveUploadedFile(file, handler, "artists")
			if err != nil {
				log.Printf("Error saving photo file: %v", err)
			} else {
				photoURL = savedURL
				log.Printf("Photo saved to: %s", photoURL)
			}
		}
	}

	if albumURL == "" {
		if file, handler, err := r.FormFile("album_file"); err == nil {
			defer file.Close()
			savedURL, err := h.saveUploadedFile(file, handler, "albums")
			if err != nil {
				log.Printf("Error saving album cover file: %v", err)
			} else {
				albumURL = savedURL
				log.Printf("Album cover saved to: %s", albumURL)
			}
		}
	}

	artist := concert.Artist{
		Nom:      nom,
		Genre:    genre,
		PhotoURL: photoURL,
		AlbumURL: albumURL,
	}

	log.Printf("Saving artist to database: Name=%s, PhotoURL=%s, AlbumURL=%s", nom, photoURL, albumURL)
	savedArtist, err := h.Service.SetArtist(artist)
	if err != nil {
		log.Printf("Error creating artist: %v", err)
		http.Error(w, "Failed to create artist: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Artist saved successfully - ID: %d, PhotoURL: %s, AlbumURL: %s", savedArtist.ID, savedArtist.PhotoURL, savedArtist.AlbumURL)
	http.Redirect(w, r, "/shows/new", http.StatusSeeOther)
}

func (h *Handler) ParticipateShow(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}
	name := r.FormValue("nom")
	showID, err := strconv.ParseUint(r.FormValue("show_id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid ShowID", http.StatusBadRequest)
		return
	}
	price, err := strconv.Atoi(r.FormValue("price"))
	if err != nil {
		http.Error(w, "Invalid Price", http.StatusBadRequest)
		return
	}
	show, err := h.Service.GetShowByID(uint(showID))
	if err != nil {
		log.Printf("Error getting show by ID %d: %v", showID, err)
		http.Error(w, "Show not found", http.StatusBadRequest)
		return
	}

	fan := concert.Fan{
		Nom:    name,
		ShowID: uint(showID),
		Show:   show,
		Price:  price,
	}

	_, err = h.Service.ParticipateShow(fan)
	if err != nil {
		log.Printf("Error registering fan: %v", err)
		http.Error(w, "Failed to register fan: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/shows", http.StatusSeeOther)
}

func (h *Handler) NewFan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	shows, err := h.Service.ListAllShow()
	if err != nil {
		log.Printf("Error listing shows: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	
	data := ShowInfo{
		Shows: shows,
	}

	tmpl := template.Must(template.ParseFiles(
		getTemplatePath("base.html"),
		getTemplatePath("participate.html"),
	))
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error rendering fan form: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) saveUploadedFile(file multipart.File, handler *multipart.FileHeader, folder string) (string, error) {
	defer file.Close()

	projectRoot := getProjectRoot()
	uploadDir := filepath.Join(projectRoot, "static", "images", folder)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	ext := filepath.Ext(handler.Filename)
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
