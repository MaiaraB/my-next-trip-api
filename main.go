package main

import (
	"log"
	"net/http"
	"os"

	"github.com/MaiaraB/travel-plan/handler"
)

func main() {
	log.Printf("Server started")

	router := handler.NewRouter()
	port := getPort()
	log.Fatal(http.ListenAndServe(port, router))
}

func getPort() string {
	p := os.Getenv("PORT")
	if p != "" {
		return ":" + p
	}
	return ":8080"
}
