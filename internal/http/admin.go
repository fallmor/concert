package http

import (
	"concert/internal/concert"
	"concert/internal/utils"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) getCurrentUser(r *http.Request) *concert.User {
	if !Authenticated(r) {
		return nil
	}
	user, err := GetUserFromCookie(h.Db, r)
	if err != nil {
		return nil
	}
	return user
}

func (h *Handler) GetAdmin(w http.ResponseWriter, r *http.Request) {
	var stats struct {
		TotalUsers   int64
		TotalShows   int64
		TotalArtists int64
		TotalFans    int64
		User         *concert.User
	}

	h.Db.Model(&concert.User{}).Count(&stats.TotalUsers)
	h.Db.Model(&concert.Show{}).Count(&stats.TotalShows)
	h.Db.Model(&concert.Artist{}).Count(&stats.TotalArtists)
	h.Db.Model(&concert.Fan{}).Count(&stats.TotalFans)
	stats.User = h.getCurrentUser(r)

	utils.RenderTemplate(w, "admin.html", stats)
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.Service.GetAllUsers()
	if err != nil {
		log.Printf("Error getting all users: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := struct {
		UserInfo
		User *concert.User
	}{
		UserInfo: UserInfo{Users: users},
		User:     h.getCurrentUser(r),
	}

	utils.RenderTemplate(w, "admin/users.html", data)
}

func (h *Handler) AdminRoutes(r chi.Router) {
	r.Get("/admin", h.GetAdmin)
	r.Get("/admin/users", h.GetUsers)
	r.Get("/shows/new", h.NewShow)
	r.Post("/show", h.SetShow)
	r.Get("/artist/new", h.NewArtist)
	r.Post("/artist", h.SetArtist)

	// Admin can only delete
	r.Delete("/shows/{id}", h.DeleteShow)
	r.Delete("/artists/{id}", h.DeleteArtist)
	r.Delete("/fans/{id}", h.DeleteFan)
}

func (h *Handler) DeleteShow(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := utils.ParseID(idStr)
	if err != nil {
		http.Error(w, "Invalid show ID", http.StatusBadRequest)
		return
	}

	if err := h.Db.Delete(&concert.Show{}, id).Error; err != nil {
		log.Printf("Error deleting show %d: %v", id, err)
		http.Error(w, "Failed to delete show", http.StatusInternalServerError)
		return
	}

	utils.WriteJSONSuccess(w)
}

func (h *Handler) DeleteArtist(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := utils.ParseID(idStr)
	if err != nil {
		http.Error(w, "Invalid artist ID", http.StatusBadRequest)
		return
	}

	if err := h.Db.Delete(&concert.Artist{}, id).Error; err != nil {
		log.Printf("Error deleting artist %d: %v", id, err)
		http.Error(w, "Failed to delete artist", http.StatusInternalServerError)
		return
	}

	utils.WriteJSONSuccess(w)
}

func (h *Handler) DeleteFan(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := utils.ParseID(idStr)
	if err != nil {
		http.Error(w, "Invalid fan ID", http.StatusBadRequest)
		return
	}

	if err := h.Db.Delete(&concert.Fan{}, id).Error; err != nil {
		log.Printf("Error deleting fan %d: %v", id, err)
		http.Error(w, "Failed to delete fan", http.StatusInternalServerError)
		return
	}

	utils.WriteJSONSuccess(w)
}
