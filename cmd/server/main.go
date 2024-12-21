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
	log.Println("starting the server")

	db, err := database.DbSetup()
	if err != nil {
		log.Println("Could not setup the Db")
		return err
	}
	if err := database.Migrate(db); err != nil {
		log.Fatal("Could not migrate the database")
	}
	myconcert := concert.NewConcert(db)
	handler := httptransport.NewRouter(myconcert)
	handler.ChiSetRoutes()
	fmt.Println("connected to the database")

	if err := http.ListenAndServe(":8080", handler.Route); err != nil {
		fmt.Println("Failed to set up server")
		return err
	}
	return nil

}

func main() {
	if err := Run(); err != nil {
		fmt.Println("Can't start the server")
	}
}
