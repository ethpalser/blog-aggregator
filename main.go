package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/ethpalser/blog-aggregator/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load(".env")
	port := os.Getenv("PORT")
	dbURL := os.Getenv("CONN")

	db, dbErr := sql.Open("postgres", dbURL)
	if dbErr != nil {
		log.Fatalf("Failed to connect to database: %s", dbURL)
	}
	defer db.Close()

	apiCfg := apiConfig{
		DB: database.New(db),
	}

	mux := http.NewServeMux()
	// Test handlers
	mux.HandleFunc("GET /v1/readiness", handlerReady)
	mux.HandleFunc("GET /v1/err", handlerError)
	// Users
	mux.HandleFunc("POST /v1/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("GET /v1/users", apiCfg.middlewareAuth(apiCfg.handlerGetUserByApiKey))
	// Feeds
	mux.HandleFunc("POST /v1/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	mux.HandleFunc("GET /v1/feeds", apiCfg.handlerGetAllFeeds)
	mux.HandleFunc("GET /v1/feeds/next-to-fetch", apiCfg.handlerGetNextToFetchFeeds)
	// Follows (User_Feeds)
	mux.HandleFunc("POST /v1/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
	mux.HandleFunc("GET /v1/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerGetUserFeedFollows))
	mux.HandleFunc("DELETE /v1/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollow))

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
