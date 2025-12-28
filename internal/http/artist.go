package http

import (
	"concert/internal/concert"
	"concert/internal/utils"
	"encoding/json"
	"log"
	"net/http"
	"text/template"
)

func (h *Handler) ListAllArtists(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	artists, err := h.Service.ListAllArtists()
	if err != nil {
		log.Printf("Error listing artists: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := struct {
		ArtistInfo
		User *concert.User
	}{
		ArtistInfo: ArtistInfo{Artists: artists},
		User:       h.getCurrentUser(r),
	}

	tmpl := template.Must(template.ParseFiles(
		utils.GetTemplatePath("base.html"),
		utils.GetTemplatePath("allartists.html"),
	))
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error executing artists template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) NewArtist(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	data := struct {
		User *concert.User
	}{
		User: h.getCurrentUser(r),
	}

	tmpl := template.Must(template.ParseFiles(
		utils.GetTemplatePath("base.html"),
		utils.GetTemplatePath("newartist.html"),
	))
	if err := tmpl.Execute(w, data); err != nil {
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

	Name := r.FormValue("Name")
	genre := r.FormValue("genre")
	ImageURL := r.FormValue("photo_url")
	albumURL := r.FormValue("album_url")

	if Name == "" || genre == "" {
		http.Error(w, "Name and genre are required", http.StatusBadRequest)
		return
	}

	if ImageURL == "" {
		if file, handler, err := r.FormFile("photo_file"); err == nil {
			defer file.Close()
			savedURL, err := h.saveUploadedFile(file, handler, "artists")
			if err != nil {
				log.Printf("Error saving photo file: %v", err)
			} else {
				ImageURL = savedURL
				log.Printf("Photo saved to: %s", ImageURL)
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
		Name:     Name,
		Genre:    genre,
		ImageURL: ImageURL,
		AlbumURL: albumURL,
	}

	log.Printf("Saving artist to database: Name=%s, ImageURL=%s, AlbumURL=%s", Name, ImageURL, albumURL)
	savedArtist, err := h.Service.SetArtist(artist)
	if err != nil {
		log.Printf("Error creating artist: %v", err)
		http.Error(w, "Failed to create artist: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Artist saved successfully - ID: %d, ImageURL: %s, AlbumURL: %s", savedArtist.ID, savedArtist.ImageURL, savedArtist.AlbumURL)
	http.Redirect(w, r, "/shows/new", http.StatusSeeOther)
}
