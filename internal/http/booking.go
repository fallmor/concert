package http

import (
	"concert/internal/models"
	"concert/internal/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"

	"github.com/go-chi/chi/v5"
)

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
		utils.GetTemplatePath("base.html"),
		utils.GetTemplatePath("fan.html"),
	))

	data := struct {
		PageData
		User *models.User
	}{
		PageData: PageData{
			Title: fmt.Sprintf("Show Fan named %s", name),
			Fans:  fanDetails,
		},
		User: h.getCurrentUser(r),
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error executing fan template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) NewFan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	shows, err := h.Service.ListAllShow()
	if err != nil {
		log.Printf("Error listing shows: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := struct {
		ShowInfo
		User *models.User
	}{
		ShowInfo: ShowInfo{Shows: shows},
		User:     h.getCurrentUser(r),
	}

	tmpl := template.Must(template.ParseFiles(
		utils.GetTemplatePath("base.html"),
		utils.GetTemplatePath("participate.html"),
	))
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error rendering fan form: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

type CreateBookingRequest struct {
	ShowID      uint `json:"showId"`
	TicketCount int  `json:"ticketCount"`
}

type BookingResponse struct {
	ID          uint    `json:"ID"`
	ShowID      uint    `json:"showId"`
	ShowTitle   string  `json:"showTitle"`
	TicketCount int     `json:"ticketCount"`
	TotalPrice  float64 `json:"totalPrice"`
	Status      string  `json:"status"`
	BookingDate string  `json:"bookingDate"`
}

func (h *Handler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, err := GetUserFromCookie(h.Db, r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ShowID == 0 || req.TicketCount <= 0 {
		http.Error(w, "Invalid show ID or ticket count", http.StatusBadRequest)
		return
	}

	show, err := h.Service.GetShowByID(req.ShowID)
	if err != nil {
		http.Error(w, "Show not found", http.StatusNotFound)
		return
	}

	if show.AvailableSeats < req.TicketCount {
		http.Error(w, "Not enough seats available", http.StatusBadRequest)
		return
	}

	totalPrice := show.Price * float64(req.TicketCount)

	booking := models.Booking{
		UserID:      user.ID,
		ShowID:      req.ShowID,
		TicketCount: req.TicketCount,
		TotalPrice:  totalPrice,
		Status:      "confirmed",
	}

	tx := h.Db.Begin()

	if err := tx.Create(&booking).Error; err != nil {
		tx.Rollback()
		log.Printf("Error creating booking: %v", err)
		http.Error(w, "Failed to create booking", http.StatusInternalServerError)
		return
	}

	if err := tx.Model(&show).Update("available_seats", show.AvailableSeats-req.TicketCount).Error; err != nil {
		tx.Rollback()
		log.Printf("Error updating seats: %v", err)
		http.Error(w, "Failed to update seats", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("Error committing transaction: %v", err)
		http.Error(w, "Failed to complete booking", http.StatusInternalServerError)
		return
	}

	var fullBooking models.Booking
	if err := h.Db.Preload("Show.Artist").First(&fullBooking, booking.ID).Error; err != nil {
		log.Printf("Error reloading booking: %v", err)
		response := BookingResponse{
			ID:          booking.ID,
			ShowID:      show.ID,
			ShowTitle:   show.Title,
			TicketCount: req.TicketCount,
			TotalPrice:  totalPrice,
			Status:      booking.Status,
			BookingDate: booking.CreatedAt.Format("2006-01-02"),
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(fullBooking)
}

func (h *Handler) GetMyBookings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, err := GetUserFromCookie(h.Db, r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	bookings, err := h.Service.GetMyBookings(*user)

	if err != nil {
		log.Printf("Error fetching bookings: %v", err)
		http.Error(w, "Failed to fetch bookings", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(bookings)
}

func (h *Handler) CancelBooking(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, err := GetUserFromCookie(h.Db, r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idString := chi.URLParam(r, "id")
	bookingID, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Invalid booking ID", http.StatusBadRequest)
		return
	}

	booking, err := h.Service.GetBookingById(uint(bookingID))
	if err != nil {
		http.Error(w, "Booking not found", http.StatusNotFound)
		return
	}

	if booking.UserID != user.ID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if booking.Status == "cancelled" {
		http.Error(w, "Booking already cancelled", http.StatusBadRequest)
		return
	}

	tx := h.Db.Begin()

	if err := tx.Model(&booking).Update("status", "cancelled").Error; err != nil {
		tx.Rollback()
		http.Error(w, "Failed to cancel booking", http.StatusInternalServerError)
		return
	}

	if err := tx.Model(&booking.Show).Update("available_seats", booking.Show.AvailableSeats+booking.TicketCount).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Failed to update seats", http.StatusInternalServerError)
		return
	}

	tx.Commit()

	json.NewEncoder(w).Encode(map[string]string{"message": "Booking cancelled successfully"})
}
