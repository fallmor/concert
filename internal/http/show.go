package http

import (
	"concert/internal/concert"

	"concert/internal/utils"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)


func (h *Handler) ListAllFan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	allFans, err := h.Service.ListAllFan()
	if err != nil {
		log.Printf("Error listing fans: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := struct {
		FanInfo
		User *concert.User
	}{
		FanInfo: FanInfo{Fans: allFans},
		User:    h.getCurrentUser(r),
	}

	tmpl := template.Must(template.ParseFiles(
		utils.GetTemplatePath("base.html"),
		utils.GetTemplatePath("allfans.html"),
	))

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error executing fans template: %v", err)
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

	data := struct {
		ShowInfo
		User *concert.User
	}{
		ShowInfo: ShowInfo{Shows: shows},
		User:     h.getCurrentUser(r),
	}

	tmpl := template.Must(template.ParseFiles(
		utils.GetTemplatePath("base.html"),
		utils.GetTemplatePath("allshows.html"),
	))

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error executing shows template: %v", err)
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

	data := struct {
		ArtistInfo
		User *concert.User
	}{
		ArtistInfo: ArtistInfo{Artists: artists},
		User:       h.getCurrentUser(r),
	}

	tmpl := template.Must(template.ParseFiles(
		utils.GetTemplatePath("base.html"),
		utils.GetTemplatePath("newshow.html"),
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

func (h *Handler) GetShow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	artistName := chi.URLParam(r, "artistname")
	showDetails, err := h.Service.GetShow(artistName)
	if err != nil {
		log.Printf("Error getting shows for artist %s: %v", artistName, err)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	data := struct {
		ShowInfo
		User *concert.User
	}{
		ShowInfo: ShowInfo{Shows: showDetails},
		User:     h.getCurrentUser(r),
	}

	tmpl := template.Must(template.ParseFiles(
		utils.GetTemplatePath("base.html"),
		utils.GetTemplatePath("show_by_artist.html"),
	))

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error executing show template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ParticipateShow(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Get the current user from the cookie
	user, err := GetUserFromCookie(h.Db, r)
	if err != nil {
		log.Printf("Error getting user from cookie: %v", err)
		http.Error(w, "You must be logged in to register for a show", http.StatusUnauthorized)
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
		UserID: user.ID,
	}

	_, err = h.Service.ParticipateShow(fan)
	if err != nil {
		log.Printf("Error registering fan: %v", err)
		http.Error(w, "Failed to register fan: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/shows", http.StatusSeeOther)
}