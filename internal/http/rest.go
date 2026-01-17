package http

import (
	"concert/internal/concert"
	"concert/internal/utils"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
	"go.temporal.io/sdk/client"
	"gorm.io/gorm"
)

func NewRouter(service concert.ConcertService, db *gorm.DB) (*Handler, error) {
	tmpClient, err := client.Dial(
		client.Options{HostPort: "127.0.0.1:7233"})
	if err != nil {
		return nil, err
	}
	return &Handler{
		Service:        service,
		Db:             db,
		TemporalClient: tmpClient,
	}, nil
}

func (h *Handler) Close() {
	if h.TemporalClient != nil {
		h.TemporalClient.Close()
	}
}

func (h *Handler) ChiSetRoutes() {
	h.Route = chi.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:8080"}, // Vite dev + Go production
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})
	h.Route.Use(c.Handler)

	h.Route.Group(func(r chi.Router) {
		r.Use(RateLimit)
		r.Post("/api/public/register", h.RegisterAPI)
		r.Post("/api/public/login", h.LoginAPI)

		r.Post("/api/public/forget", h.ForgetPassword)
		r.Get("/api/public/shows", h.ListAllShow)
		r.Get("/api/public/shows/{id}", h.GetShowPublic)
		r.Get("/api/public/artists/{id}", h.GetArtistPublic)
		r.Get("/api/public/artists", h.ListAllArtists)

	})

	h.Route.Group(func(r chi.Router) {
		r.Use(NeedsAuth(h.Db))

		r.Post("/api/bookings", h.CreateBooking)
		r.Get("/api/bookings", h.GetMyBookings)
		r.Delete("/api/bookings/{id}", h.CancelBooking)
	})

	h.Route.Group(func(r chi.Router) {
		r.Use(NeedsAuth(h.Db))
		r.Use(NeedsRole(h.Db, "admin"))
		h.AdminRoutes(r)
	})

	staticDir := utils.GetStaticDir()
	h.Route.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir(staticDir))))

	assetFs := http.FileServer(http.Dir(staticDir))
	h.Route.Get("/assets/*", http.HandlerFunc(assetFs.ServeHTTP))

	h.Route.Get("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		indexPath := filepath.Join(staticDir, "index.html")
		data, err := os.ReadFile(indexPath)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(data)
	}))
}
