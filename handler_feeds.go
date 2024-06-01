package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/ethpalser/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, u database.User) {
	type CreateFeedReq struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}
	req := CreateFeedReq{}
	decoder := json.NewDecoder(r.Body)
	dErr := decoder.Decode(&req)
	if dErr != nil {
		respondWithError(w, http.StatusInternalServerError, "invalid request structure")
		return
	}

	now := time.Now()
	ctx := context.Background()
	feed, dbErr := cfg.DB.CreateFeed(ctx, database.CreateFeedParams{
		ID:        uuid.NullUUID{UUID: uuid.New(), Valid: true},
		CreatedAt: sql.NullTime{Time: now, Valid: true},
		UpdatedAt: sql.NullTime{Time: now, Valid: true},
		Name:      sql.NullString{String: req.Name, Valid: true},
		Url:       sql.NullString{String: req.Url, Valid: true},
		UserID:    u.ID,
	})
	if dbErr != nil {
		log.Printf("Failed to create feed: %s", dbErr)
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondWithJSON(w, http.StatusAccepted, feed)
}
