package http

import (
	"concert/internal/models"

	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) ListAllShow(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	shows, err := h.Service.ListAllShow()
	if err != nil {
		log.Printf("Error listing shows: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}


	json.NewEncoder(w).Encode(shows)

}

func (h *Handler) GetShowPublic(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idString := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		log.Println("Invalid ID:", err)
		http.Error(w, "Invalid show ID", http.StatusBadRequest)
		return
	}

	show, err := h.Service.GetShowByID(uint(id))
	if err != nil {
		http.Error(w, "Could not get the show by ID", http.StatusBadRequest)
		return
	}
	// Calculate available seats

	json.NewEncoder(w).Encode(show)
}	

func (h *Handler) SetShow(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")

	if contentType == "application/json" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		var show models.Show
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

	Venue := r.FormValue("Venue")
	if Venue == "" {
		http.Error(w, "Venue is required", http.StatusBadRequest)
		return
	}

	show := models.Show{
		Date:     date,
		ArtistID: uint(artistID),
		Venue:    Venue,
	}

	_, err = h.Service.SetShow(show)
	if err != nil {
		log.Printf("Error creating show: %v", err)
		http.Error(w, "Failed to create show: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/shows", http.StatusSeeOther)
}


func (h *Handler) ParticipateShow(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	user, err := GetUserFromCookie(h.Db, r)
	if err != nil {
		log.Printf("Error getting user from cookie: %v", err)
		http.Error(w, "You must be logged in to register for a show", http.StatusUnauthorized)
		return
	}

	//name := r.FormValue("Name")
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

	fan := models.Booking{
		User:       *user,
		ShowID:     uint(showID),
		Show:       show,
		TotalPrice: float64(price),
		UserID:     user.ID,
	}

	_, err = h.Service.ParticipateShow(fan)
	if err != nil {
		log.Printf("Error registering fan: %v", err)
		http.Error(w, "Failed to register fan: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/shows", http.StatusSeeOther)
}
