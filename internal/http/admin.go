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

type CreateShowRequest struct {
	Title       string  `json:"title"`
	Date        string  `json:"date"` // "2026-02-15"
	Time        string  `json:"time"` // "20:00"
	ArtistID    uint    `json:"artistId"`
	Venue       string  `json:"venue"`
	Price       float64 `json:"price"`
	TotalSeats  int     `json:"totalSeats"`
	Description string  `json:"description"`
	ImageURL    string  `json:"imageUrl"`
}

type CreateArtistRequest struct {
	Name     string `json:"name"`
	Genre    string `json:"genre"`
	Bio      string `json:"bio"`
	ImageURL string `json:"imageUrl"`
}

type AdminBookingResponse struct {
	ID          uint    `json:"id"`
	UserID      uint    `json:"userId"`
	Username    string  `json:"username"`
	UserEmail   string  `json:"userEmail"`
	ShowID      uint    `json:"showId"`
	ShowTitle   string  `json:"showTitle"`
	ArtistName  string  `json:"artistName"`
	TicketCount int     `json:"ticketCount"`
	TotalPrice  float64 `json:"totalPrice"`
	Status      string  `json:"status"`
	CreatedAt   string  `json:"createdAt"`
}

type Stats struct {
	TotalShows    int64   `json:"totalShows"`
	TotalArtists  int64   `json:"totalArtists"`
	TotalBookings int64   `json:"totalBookings"`
	TotalRevenue  float64 `json:"totalRevenue"`
	TotalUsers    int64   `json:"totalUsers"`
}

func (h *Handler) ListShows(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	shows, err := h.Service.ListAllShow()
	if err != nil {
		http.Error(w, "Could not list shows", http.StatusInternalServerError)
	}

	// Check available seats
	for i := range shows {
		confirmBook := h.Service.CountConfirmSeats(uint(i))
		shows[i].AvailableSeats = shows[i].TotalSeats - int(confirmBook)
	}

	json.NewEncoder(w).Encode(shows)
}
func (h *Handler) getCurrentUser(r *http.Request) *models.User {
	if !Authenticated(r) {
		return nil
	}
	user, err := GetUserFromCookie(h.Db, r)
	if err != nil {
		return nil
	}
	return user
}

// func (h *Handler) GetAdmin(w http.ResponseWriter, r *http.Request) {
// 	var stats struct {
// 		TotalUsers   int64
// 		TotalShows   int64
// 		TotalArtists int64
// 		TotalFans    int64
// 		User         *models.User
// 	}

// 	h.Db.Model(&models.User{}).Count(&stats.TotalUsers)
// 	h.Db.Model(&models.Show{}).Count(&stats.TotalShows)
// 	h.Db.Model(&models.Artist{}).Count(&stats.TotalArtists)
// 	h.Db.Model(&models.Booking{}).Count(&stats.TotalFans)
// 	stats.User = h.getCurrentUser(r)

// 	utils.RenderTemplate(w, "admin.html", stats)
// }

func (h *Handler) CreateShow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req CreateShowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var artist models.Artist
	if err := h.Db.First(&artist, req.ArtistID).Error; err != nil {
		http.Error(w, "Artist not found", http.StatusNotFound)
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		http.Error(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	if req.Time == "" {
		req.Time = "20:00"
	}
	if req.TotalSeats == 0 {
		req.TotalSeats = 100
	}
	if req.Price == 0 {
		req.Price = 50.0
	}

	show := models.Show{
		Title:          req.Title,
		Date:           date,
		Time:           req.Time,
		ArtistID:       req.ArtistID,
		Venue:          req.Venue,
		Price:          req.Price,
		TotalSeats:     req.TotalSeats,
		AvailableSeats: req.TotalSeats,
		Description:    req.Description,
		ImageURL:       req.ImageURL,
	}

	show, err = h.Service.SetShow(show)
	if err != nil {
		log.Printf("Error creating show: %v", err)
		http.Error(w, "Failed to create show", http.StatusInternalServerError)
		return
	}

	log.Printf("Show created: %s by %s on %s", show.Title, artist.Name, show.Date.Format("2006-01-02"))

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(show)
}

func (h *Handler) UpdateShow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idString := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var show models.Show
	show, err = h.Service.GetShowByID(uint(id))
	if err != nil {
		http.Error(w, "Show not found", http.StatusNotFound)
		return
	}

	var req CreateShowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Title != "" {
		show.Title = req.Title
	}
	if req.Date != "" {
		date, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			http.Error(w, "Invalid date format", http.StatusBadRequest)
			return
		}
		show.Date = date
	}
	if req.Time != "" {
		show.Time = req.Time
	}
	if req.ArtistID != 0 {
		h.Service.GetArtistByID(req.ArtistID)
		_, err := h.Service.GetArtistByID(req.ArtistID)
		if err != nil {
			http.Error(w, "Artist not found", http.StatusNotFound)
			return
		}
		show.ArtistID = req.ArtistID
	}
	if req.Venue != "" {
		show.Venue = req.Venue
	}
	if req.Price > 0 {
		show.Price = req.Price
	}
	if req.TotalSeats > 0 {
		show.TotalSeats = req.TotalSeats
	}
	if req.Description != "" {
		show.Description = req.Description
	}
	if req.ImageURL != "" {
		show.ImageURL = req.ImageURL
	}

	show, err = h.Service.SetShow(show)

	if err != nil {
		log.Printf("Error updating show: %v", err)
		http.Error(w, "Failed to update show", http.StatusInternalServerError)
		return
	}

	log.Printf("Show updated: %s (ID: %d)", show.Title, show.ID)

	json.NewEncoder(w).Encode(show)
}

func (h *Handler) DeleteShow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idString := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	_, err = h.Service.GetShowByID(uint(id))
	if err != nil {
		http.Error(w, "Show not found", http.StatusNotFound)
		return
	}

	confirmBook := h.Service.CountConfirmSeats(uint(id))

	if confirmBook > 0 {
		http.Error(w, "Cannot delete show with confirmed bookings", http.StatusBadRequest)
		return
	}

	h.Service.DeleteShow(uint(id))
	if err != nil {
		log.Printf("Error deleting show: %v", err)
		http.Error(w, "Failed to delete show", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Show deleted successfully"})
}

func (h *Handler) ListArtists(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	artists, err := h.Service.ListAllArtists()
	if err != nil {
		log.Printf("Error listing artists: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	log.Printf("Returning %d artists", len(artists)) // ← Add debug
	for _, a := range artists {
		log.Printf("  - Artist: ID=%d, Name=%s", a.ID, a.Name) // ← Add debug
	}
	json.NewEncoder(w).Encode(artists)
}

func (h *Handler) CreateArtist(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req CreateArtistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	artist := models.Artist{
		Name:     req.Name,
		Genre:    req.Genre,
		Bio:      req.Bio,
		ImageURL: req.ImageURL,
	}

	_, err := h.Service.SetArtist(artist)

	if err != nil {
		log.Printf("Error creating artist: %v", err)
		http.Error(w, "Failed to create artist", http.StatusInternalServerError)
		return
	}

	log.Printf("Artist created: %s", artist.Name)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(artist)
}

func (h *Handler) UpdateArtist(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idString := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var artist models.Artist
	_, err = h.Service.GetArtistByID(uint(id))
	if err != nil {
		http.Error(w, "Artist not found", http.StatusNotFound)
		return
	}

	var req CreateArtistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name != "" {
		artist.Name = req.Name
	}
	if req.Genre != "" {
		artist.Genre = req.Genre
	}
	if req.Bio != "" {
		artist.Bio = req.Bio
	}
	if req.ImageURL != "" {
		artist.ImageURL = req.ImageURL
	}

	_, err = h.Service.SetArtist(artist)
	if err != nil {
		log.Printf("Error updating artist: %v", err)
		http.Error(w, "Failed to update artist", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(artist)
}

func (h *Handler) DeleteArtist(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idString := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	artist, err := h.Service.GetArtistByID(uint(id))
	if err != nil {
		http.Error(w, "Artist not found", http.StatusNotFound)
		return
	}

	var showCount int64
	h.Db.Model(&models.Show{}).Where("artist_id = ?", id).Count(&showCount)

	if showCount > 0 {
		http.Error(w, "Cannot delete artist with existing shows", http.StatusBadRequest)
		return
	}

	h.Service.DeleteArtist(uint(artist.ID))
	if err != nil {
		log.Printf("Error deleting artist: %v", err)
		http.Error(w, "Failed to delete artist", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Artist deleted successfully"})
}

func (h *Handler) ListBookings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	bookings, err := h.Service.GetAllBookings()
	if err != nil {
		log.Printf("Error fetching bookings: %v", err)
		http.Error(w, "Failed to fetch bookings", http.StatusInternalServerError)
		return
	}

	responses := make([]AdminBookingResponse, len(bookings))
	for i, booking := range bookings {
		responses[i] = AdminBookingResponse{
			ID:          booking.ID,
			UserID:      booking.UserID,
			Username:    booking.User.Username,
			UserEmail:   booking.User.Email,
			ShowID:      booking.ShowID,
			ShowTitle:   booking.Show.Title,
			ArtistName:  booking.Show.Artist.Name,
			TicketCount: booking.TicketCount,
			TotalPrice:  booking.TotalPrice,
			Status:      booking.Status,
			CreatedAt:   booking.CreatedAt.Format("2006-01-02 15:04"),
		}
	}

	json.NewEncoder(w).Encode(responses)
}

// Get admin stats
func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var stats Stats

	h.Db.Model(&models.Show{}).Count(&stats.TotalShows)
	h.Db.Model(&models.Artist{}).Count(&stats.TotalArtists)
	h.Db.Model(&models.Booking{}).Count(&stats.TotalBookings)
	h.Db.Model(&models.User{}).Count(&stats.TotalUsers)

	h.Db.Model(&models.Booking{}).
		Where("status = ?", "confirmed").
		Select("COALESCE(SUM(total_price), 0)").
		Scan(&stats.TotalRevenue)

	json.NewEncoder(w).Encode(stats)
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	// users, err := h.Service.GetAllUsers()
	// if err != nil {
	// 	log.Printf("Error getting all users: %v", err)
	// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// 	return
	// }

	// data := struct {
	// 	UserInfo
	// 	User *concert.User
	// }{
	// 	UserInfo: UserInfo{Users: users},
	// 	User:     h.getCurrentUser(r),
	// }

	// utils.RenderTemplate(w, "admin/users.html", data)
}

func (h *Handler) AdminRoutes(r chi.Router) {
	// r.Get("/admin", h.GetAdmin)
	// //	r.Get("/admin/users", h.GetUsers)
	// r.Get("/shows/new", h.NewShow)
	// r.Post("/show", h.SetShow)
	// r.Get("/artist/new", h.NewArtist)
	// r.Post("/artist", h.SetArtist)

	// // Admin can only delete
	// r.Delete("/shows/{id}", h.DeleteShow)
	// r.Delete("/artists/{id}", h.DeleteArtist)
	// r.Delete("/fans/{id}", h.DeleteFan)

	r.Get("/api/admin/shows", h.ListShows)
	r.Post("/api/admin/shows", h.CreateShow)
	r.Put("/api/admin/shows/{id}", h.UpdateShow)
	r.Delete("/api/admin/shows/{id}", h.DeleteShow)

	// Artists
	r.Get("/api/admin/artists", h.ListArtists)
	r.Post("/api/admin/artists", h.CreateArtist)
	r.Put("/api/admin/artists/{id}", h.UpdateArtist)
	r.Delete("/api/admin/artists/{id}", h.DeleteArtist)

	// Bookings
	r.Get("/api/admin/bookings", h.ListBookings)

	// Stats
	r.Get("/api/admin/stats", h.GetStats)
}

// func (h *Handler) DeleteShow(w http.ResponseWriter, r *http.Request) {
// 	idStr := chi.URLParam(r, "id")
// 	id, err := utils.ParseID(idStr)
// 	if err != nil {
// 		http.Error(w, "Invalid show ID", http.StatusBadRequest)
// 		return
// 	}

// 	if err := h.Db.Delete(&models.Show{}, id).Error; err != nil {
// 		log.Printf("Error deleting show %d: %v", id, err)
// 		http.Error(w, "Failed to delete show", http.StatusInternalServerError)
// 		return
// 	}

// 	utils.WriteJSONSuccess(w)
// }

// func (h *Handler) DeleteArtist(w http.ResponseWriter, r *http.Request) {
// 	idStr := chi.URLParam(r, "id")
// 	id, err := utils.ParseID(idStr)
// 	if err != nil {
// 		http.Error(w, "Invalid artist ID", http.StatusBadRequest)
// 		return
// 	}

// 	if err := h.Db.Delete(&models.Artist{}, id).Error; err != nil {
// 		log.Printf("Error deleting artist %d: %v", id, err)
// 		http.Error(w, "Failed to delete artist", http.StatusInternalServerError)
// 		return
// 	}

// 	utils.WriteJSONSuccess(w)
// }

// func (h *Handler) DeleteFan(w http.ResponseWriter, r *http.Request) {
// 	idStr := chi.URLParam(r, "id")
// 	id, err := utils.ParseID(idStr)
// 	if err != nil {
// 		http.Error(w, "Invalid fan ID", http.StatusBadRequest)
// 		return
// 	}

// 	if err := h.Db.Delete(&models.Booking{}, id).Error; err != nil {
// 		log.Printf("Error deleting fan %d: %v", id, err)
// 		http.Error(w, "Failed to delete fan", http.StatusInternalServerError)
// 		return
// 	}

// 	utils.WriteJSONSuccess(w)
// }
