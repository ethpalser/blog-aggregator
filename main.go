package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load(".env")
	port := os.Getenv("PORT")

	mux := http.NewServeMux()
	// Test handlers
	mux.HandleFunc("GET /v1/readiness", handlerReady)
	mux.HandleFunc("GET /v1/err", handlerError)

	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Printf("Server closed: %s", err)
	}
}

func handlerReady(w http.ResponseWriter, r *http.Request) {
	type okResponse struct {
		Status string `json:"status"`
	}
	respondWithJSON(w, http.StatusOK, okResponse{Status: "ok"})
}

func handlerError(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
}
