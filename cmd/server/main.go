package main

import (
	"concert/internal/concert"
	"concert/internal/database"
	httptransport "concert/internal/http_transport"
	"fmt"
	"log"
	"net/http"
)

func Run() error {
	log.Println("Starting the server")

	db, err := database.DbSetup()
	if err != nil {
		return fmt.Errorf("failed to setup database: %w", err)
	}

	if err := database.Migrate(db); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database connected and migrated successfully")

	concertService := concert.NewConcert(db)
	handler := httptransport.NewRouter(concertService)
	handler.ChiSetRoutes()

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", handler.Route); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil
}

func main() {
	if err := Run(); err != nil {
		log.Fatalf("Can't start the server: %v", err)
	}
}
