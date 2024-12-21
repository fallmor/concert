package httptransport

import (
	"concert/internal/concert"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"text/template"

	"github.com/go-chi/chi/v5"
)

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

func NewRouter(service *concert.Service) *Handler {
	return &Handler{
		Service: service,
	}
}

func (h *Handler) ChiSetRoutes() {
	h.Route = chi.NewRouter()
	fileServer := http.FileServer(http.Dir("../../static"))
	h.Route.Handle("/static/css/*", http.StripPrefix("/static/", fileServer))
	h.Route.Get("/fan/{name}", h.GetFan)
	h.Route.Get("/show/{artistname}", h.GetShow)
	h.Route.Get("/fans/new", h.NewFan)
	h.Route.Get("/shows", h.ListAllShow)
	h.Route.Get("/list", h.ListAllFan)
	h.Route.Post("/submit", h.ParticipateShow)
	h.Route.Post("/show", h.SetShow)
	h.Route.Post("/artist", h.SetArtist)
}

func (h *Handler) GetFan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	//w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	tmpl := template.Must(template.ParseFiles(filepath.Join("..", "..", "internal", "templates", "fan.html")))

	name := chi.URLParam(r, "name")
	log.Printf("Name: %s", name)
	Fandetails, err := h.Service.GetFan(name)
	if err != nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	data := PageData{
		Title: fmt.Sprintf("Show Fan named %s", name), // Use the artist name as the title
		Fans:  Fandetails,
	}

	if err := tmpl.Execute(w, data); err != nil {
		panic(err)
	}
}

func (h *Handler) GetShow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	artistnamename := chi.URLParam(r, "artistname")
	Showdetails, err := h.Service.GetShow(artistnamename)

	if err != nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	data := ShowInfo{
		Shows: Showdetails,
	}
	tmpl := template.Must(template.ParseFiles(filepath.Join("..", "..", "internal", "templates", "show_by_artist.html")))

	if err := tmpl.Execute(w, data); err != nil {
		panic(err)
	}
	// if err := json.NewEncoder(w).Encode(Showdetails); err != nil {
	// 	panic(err)
	// }
}

func (h *Handler) ListAllShow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	//Allshow, err := h.Service.ListAllShow()
	// if err != nil {
	// 	http.Error(w, http.StatusText(404), 404)
	// 	return
	// }
	shows, _ := h.Service.ListAllShow()
	data := ShowInfo{
		Shows: shows,
	}
	tmpl := template.Must(template.ParseFiles(filepath.Join("..", "..", "internal", "templates", "allshows.html")))

	if err := tmpl.Execute(w, data); err != nil {
		panic(err)
	}

	// if err := json.NewEncoder(w).Encode(Allshow); err != nil {
	// 	panic(err)
	// }
}

func (h *Handler) ListAllFan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	Allfan, err := h.Service.ListAllFan()
	data := FanInfo{
		Fans: Allfan,
	}
	if err != nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	tmpl := template.Must(template.ParseFiles(filepath.Join("..", "..", "internal", "templates", "allfans.html")))

	if err := tmpl.Execute(w, data); err != nil {
		panic(err)
	}
	// if err := json.NewEncoder(w).Encode(Allfan); err != nil {
	// 	panic(err)
	// }
}

func (h *Handler) SetShow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	var show concert.Show
	if err := json.NewDecoder(r.Body).Decode(&show); err != nil {
		fmt.Fprintln(w, "failed to decode the body")
	}
	myShow, err := h.Service.SetShow(show)
	if err != nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	if err := json.NewEncoder(w).Encode(myShow); err != nil {
		panic(err)
	}
}

func (h *Handler) SetArtist(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	var artist concert.Artist
	if err := json.NewDecoder(r.Body).Decode(&artist); err != nil {
		fmt.Fprintln(w, "failed to decode the body")
	}
	myArtist, err := h.Service.SetArtist(artist)
	if err != nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	if err := json.NewEncoder(w).Encode(myArtist); err != nil {
		panic(err)
	}
}

func (h *Handler) ParticipateShow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

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
		http.Error(w, "can't get the show from the show id", http.StatusBadRequest)
		return
	}

	Fan := concert.Fan{
		Nom:    name,
		ShowID: uint(showID),
		Show:   show,
		Price:  price,
	}

	// if err := json.NewDecoder(r.Body).Decode(&Fan); err != nil {
	// 	fmt.Fprintln(w, "failed to decode the body")
	// }

	_, err = h.Service.ParticipateShow(Fan)
	if err != nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	// if err := json.NewEncoder(w).Encode(Showdetails); err != nil {
	// 	panic(err)
	// }
}

func (h *Handler) NewFan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	shows, _ := h.Service.ListAllShow()
	data := ShowInfo{
		Shows: shows,
	}

	tmpl := template.Must(template.ParseFiles(filepath.Join("..", "..", "internal", "templates", "participate.html")))
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error rendering fan form: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/shows", http.StatusSeeOther)
}
