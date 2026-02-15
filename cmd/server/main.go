package main

import (
	"concert/internal/concert"
	"concert/internal/database"
	httpTransport "concert/internal/http"
	"concert/internal/utils"
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
	handler, err := httpTransport.NewRouter(concertService, db)
	if err != nil {
		return fmt.Errorf("could not initialize the handler: %w", err)
	}
	defer handler.Close()
	handler.ChiSetRoutes()

	port := utils.GetEnvOrDefault("PORT", "8080")
	addr := ":" + port
	log.Printf("Server starting on %s", addr)
	if err := http.ListenAndServe(addr, handler.Route); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil
}

func main() {
	if err := Run(); err != nil {
		log.Fatalf("Can't start the server: %v", err)
	}
}
