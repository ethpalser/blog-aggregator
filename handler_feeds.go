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

type FeedView struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
}

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

	respondWithJSON(w, http.StatusAccepted, FeedView{
		Id:        feed.ID.UUID,
		CreatedAt: feed.CreatedAt.Time,
		UpdatedAt: feed.UpdatedAt.Time,
		Name:      feed.Name.String,
		Url:       feed.Url.String,
	})
}

func (cfg *apiConfig) handlerGetAllFeeds(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	dbFeed, dbErr := cfg.DB.GetAllFeeds(ctx)
	if dbErr != nil {
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	resp := []FeedView{}
	for _, feed := range dbFeed {
		resp = append(resp, DBFeedToView(&feed))
	}

	respondWithJSON(w, http.StatusOK, resp)
}

func DBFeedToView(feed *database.Feed) FeedView {
	return FeedView{
		Id:        feed.ID.UUID,
		CreatedAt: feed.CreatedAt.Time,
		UpdatedAt: feed.UpdatedAt.Time,
		Name:      feed.Name.String,
		Url:       feed.Url.String,
	}
}
