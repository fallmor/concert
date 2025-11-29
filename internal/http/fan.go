package http

import (
	"concert/internal/concert"
	"concert/internal/utils"
	"fmt"
	"log"
	"net/http"
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
		User *concert.User
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
		User *concert.User
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