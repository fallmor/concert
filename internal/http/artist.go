package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) ListAllArtists(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	artists, err := h.Service.ListAllArtists()
	if err != nil {
		log.Printf("Error listing artists: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(artists)
}

func (h *Handler) GetArtistPublic(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idString := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	response, err := h.Service.GetArtistByID(uint(id))
	if err != nil {
		http.Error(w, "Artist not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(response)
}
