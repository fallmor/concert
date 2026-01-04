package http

// import (
// 	"concert/internal/models"
// 	"concert/internal/utils"
// 	"log"
// 	"net/http"
// 	"strconv"
// 	"text/template"
// 	"time"

// 	"github.com/go-chi/chi/v5"
// )

// func (h *Handler) ModeratorRoutes(r chi.Router) {
// 	r.Get("/moderator", h.GetModerator)
// 	r.Get("/shows/{id}/edit", h.EditShow)
// 	r.Get("/artists/{id}/edit", h.EditArtist)
// 	r.Get("/fans/{id}/edit", h.EditFan)
// 	r.Post("/shows/{id}/update", h.UpdateShow)
// 	r.Post("/artists/{id}/update", h.UpdateArtist)
// 	r.Post("/fans/{id}/update", h.UpdateFan)
// }

// func (h *Handler) GetModerator(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

// 	var stats struct {
// 		TotalShows   int64
// 		TotalArtists int64
// 		TotalFans    int64
// 		User         *models.User
// 	}

// 	h.Db.Model(&models.Show{}).Count(&stats.TotalShows)
// 	h.Db.Model(&models.Artist{}).Count(&stats.TotalArtists)
// 	h.Db.Model(&models.Booking{}).Count(&stats.TotalFans)
// 	stats.User = h.getCurrentUser(r)

// 	tmpl := template.Must(template.ParseFiles(
// 		utils.GetTemplatePath("base.html"),
// 		utils.GetTemplatePath("moderator.html"),
// 	))
// 	if err := tmpl.Execute(w, stats); err != nil {
// 		log.Printf("Error executing moderator template: %v", err)
// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 		return
// 	}
// }

// func (h *Handler) EditShow(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

// 	showIDStr := chi.URLParam(r, "id")
// 	showID, err := strconv.ParseUint(showIDStr, 10, 64)
// 	if err != nil {
// 		http.Error(w, "Invalid show ID", http.StatusBadRequest)
// 		return
// 	}

// 	show, err := h.Service.GetShowByID(uint(showID))
// 	if err != nil {
// 		log.Printf("Error getting show %d: %v", showID, err)
// 		http.Error(w, "Show not found", http.StatusNotFound)
// 		return
// 	}

// 	artists, err := h.Service.ListAllArtists()
// 	if err != nil {
// 		log.Printf("Error listing artists: %v", err)
// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 		return
// 	}

// 	data := struct {
// 		Show    models.Show
// 		Artists []models.Artist
// 		User    *models.User
// 	}{
// 		Show:    show,
// 		Artists: artists,
// 		User:    h.getCurrentUser(r),
// 	}

// 	tmpl := template.Must(template.ParseFiles(
// 		utils.GetTemplatePath("base.html"),
// 		utils.GetTemplatePath("editshow.html"),
// 	))
// 	if err := tmpl.Execute(w, data); err != nil {
// 		log.Printf("Error executing edit show template: %v", err)
// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 		return
// 	}
// }

// func (h *Handler) UpdateShow(w http.ResponseWriter, r *http.Request) {
// 	showIDStr := chi.URLParam(r, "id")
// 	showID, err := strconv.ParseUint(showIDStr, 10, 64)
// 	if err != nil {
// 		http.Error(w, "Invalid show ID", http.StatusBadRequest)
// 		return
// 	}

// 	if err := r.ParseForm(); err != nil {
// 		http.Error(w, "Unable to parse form", http.StatusBadRequest)
// 		return
// 	}

// 	show, err := h.Service.GetShowByID(uint(showID))
// 	if err != nil {
// 		log.Printf("Error getting show %d: %v", showID, err)
// 		http.Error(w, "Show not found", http.StatusNotFound)
// 		return
// 	}

// 	artistID, err := strconv.ParseUint(r.FormValue("artist_id"), 10, 64)
// 	if err != nil {
// 		http.Error(w, "Invalid Artist ID", http.StatusBadRequest)
// 		return
// 	}

// 	dateStr := r.FormValue("date")
// 	date, err := time.Parse("2006-01-02T15:04", dateStr)
// 	if err != nil {
// 		http.Error(w, "Invalid date format", http.StatusBadRequest)
// 		return
// 	}

// 	show.ArtistID = uint(artistID)
// 	show.Date = date
// 	show.Venue = r.FormValue("Venue")

// 	_, err = h.Service.SetShow(show)
// 	if err != nil {
// 		log.Printf("Error updating show: %v", err)
// 		http.Error(w, "Failed to update show", http.StatusInternalServerError)
// 		return
// 	}

// 	http.Redirect(w, r, "/shows", http.StatusSeeOther)
// }

// func (h *Handler) EditArtist(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

// 	artistIDStr := chi.URLParam(r, "id")
// 	artistID, err := strconv.ParseUint(artistIDStr, 10, 64)
// 	if err != nil {
// 		http.Error(w, "Invalid artist ID", http.StatusBadRequest)
// 		return
// 	}

// 	var artist models.Artist
// 	if err := h.Db.First(&artist, artistID).Error; err != nil {
// 		log.Printf("Error getting artist %d: %v", artistID, err)
// 		http.Error(w, "Artist not found", http.StatusNotFound)
// 		return
// 	}

// 	data := struct {
// 		Artist models.Artist
// 		User   *models.User
// 	}{
// 		Artist: artist,
// 		User:   h.getCurrentUser(r),
// 	}

// 	tmpl := template.Must(template.ParseFiles(
// 		utils.GetTemplatePath("base.html"),
// 		utils.GetTemplatePath("editartist.html"),
// 	))
// 	if err := tmpl.Execute(w, data); err != nil {
// 		log.Printf("Error executing edit artist template: %v", err)
// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 		return
// 	}
// }

// func (h *Handler) UpdateArtist(w http.ResponseWriter, r *http.Request) {
// 	artistIDStr := chi.URLParam(r, "id")
// 	artistID, err := strconv.ParseUint(artistIDStr, 10, 64)
// 	if err != nil {
// 		http.Error(w, "Invalid artist ID", http.StatusBadRequest)
// 		return
// 	}

// 	if err := r.ParseMultipartForm(10 << 20); err != nil {
// 		http.Error(w, "Unable to parse form", http.StatusBadRequest)
// 		return
// 	}

// 	var artist models.Artist
// 	if err := h.Db.First(&artist, artistID).Error; err != nil {
// 		log.Printf("Error getting artist %d: %v", artistID, err)
// 		http.Error(w, "Artist not found", http.StatusNotFound)
// 		return
// 	}

// 	artist.Name = r.FormValue("Name")
// 	artist.Genre = r.FormValue("genre")

// 	if ImageURL := r.FormValue("photo_url"); ImageURL != "" {
// 		artist.ImageURL = ImageURL
// 	} else if file, handler, err := r.FormFile("photo_file"); err == nil {
// 		defer file.Close()
// 		savedURL, err := h.saveUploadedFile(file, handler, "artists")
// 		if err == nil {
// 			artist.ImageURL = savedURL
// 		}
// 	}

// 	if albumURL := r.FormValue("album_url"); albumURL != "" {
// 		artist.AlbumURL = albumURL
// 	} else if file, handler, err := r.FormFile("album_file"); err == nil {
// 		defer file.Close()
// 		savedURL, err := h.saveUploadedFile(file, handler, "albums")
// 		if err == nil {
// 			artist.AlbumURL = savedURL
// 		}
// 	}

// 	_, err = h.Service.SetArtist(artist)
// 	if err != nil {
// 		log.Printf("Error updating artist: %v", err)
// 		http.Error(w, "Failed to update artist", http.StatusInternalServerError)
// 		return
// 	}

// 	http.Redirect(w, r, "/artists", http.StatusSeeOther)
// }

// func (h *Handler) EditFan(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

// 	fanIDStr := chi.URLParam(r, "id")
// 	fanID, err := strconv.ParseUint(fanIDStr, 10, 64)
// 	if err != nil {
// 		http.Error(w, "Invalid fan ID", http.StatusBadRequest)
// 		return
// 	}

// 	var fan models.Booking
// 	if err := h.Db.Preload("Show").Preload("Show.Artist").First(&fan, fanID).Error; err != nil {
// 		log.Printf("Error getting fan %d: %v", fanID, err)
// 		http.Error(w, "Fan not found", http.StatusNotFound)
// 		return
// 	}

// 	shows, err := h.Service.ListAllShow()
// 	if err != nil {
// 		log.Printf("Error listing shows: %v", err)
// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 		return
// 	}

// 	data := struct {
// 		Fan   models.Booking
// 		Shows []models.Show
// 		User  *models.User
// 	}{
// 		Fan:   fan,
// 		Shows: shows,
// 		User:  h.getCurrentUser(r),
// 	}

// 	tmpl := template.Must(template.ParseFiles(
// 		utils.GetTemplatePath("base.html"),
// 		utils.GetTemplatePath("editfan.html"),
// 	))
// 	if err := tmpl.Execute(w, data); err != nil {
// 		log.Printf("Error executing edit fan template: %v", err)
// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 		return
// 	}
// }

// func (h *Handler) UpdateFan(w http.ResponseWriter, r *http.Request) {
// 	fanIDStr := chi.URLParam(r, "id")
// 	fanID, err := strconv.ParseUint(fanIDStr, 10, 64)
// 	if err != nil {
// 		http.Error(w, "Invalid fan ID", http.StatusBadRequest)
// 		return
// 	}

// 	if err := r.ParseForm(); err != nil {
// 		http.Error(w, "Unable to parse form", http.StatusBadRequest)
// 		return
// 	}

// 	var fan models.Booking
// 	if err := h.Db.First(&fan, fanID).Error; err != nil {
// 		log.Printf("Error getting fan %d: %v", fanID, err)
// 		http.Error(w, "Fan not found", http.StatusNotFound)
// 		return
// 	}

// 	fan.User.FirstName = r.FormValue("Name")
// 	showID, err := strconv.ParseUint(r.FormValue("show_id"), 10, 64)
// 	if err == nil {
// 		fan.ShowID = uint(showID)
// 	}
// 	// price, err := strconv.Atoi(r.FormValue("price"))
// 	// if err == nil {
// 	// 	fan.Price = price
// 	// }

// 	if err := h.Db.Save(&fan).Error; err != nil {
// 		log.Printf("Error updating fan: %v", err)
// 		http.Error(w, "Failed to update fan", http.StatusInternalServerError)
// 		return
// 	}

// 	http.Redirect(w, r, "/list", http.StatusSeeOther)
// }
