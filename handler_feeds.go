package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ethpalser/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

type FeedView struct {
	Id            uuid.UUID  `json:"id"`
	CreatedAt     *time.Time `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
	Name          string     `json:"name"`
	Url           string     `json:"url"`
	LastFetchedAt *time.Time `json:"last_fetched_at"`
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

	respondWithJSON(w, http.StatusAccepted, DBFeedToView(&feed))
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

// Temp
func (cfg *apiConfig) handlerGetNextToFetchFeeds(w http.ResponseWriter, r *http.Request) {
	type NextFetchReq struct {
		Since time.Time `json:"since"`
		Limit int       `json:"limit"`
	}

	req := NextFetchReq{}
	decoder := json.NewDecoder(r.Body)
	jsonErr := decoder.Decode(&req)
	if jsonErr != nil {
		respondWithError(w, http.StatusInternalServerError, "internal server error: failed to parse json")
		return
	}

	res, err := cfg.getNextToFetchFeeds(req.Since, int32(req.Limit))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, res)
}

func (cfg *apiConfig) getNextToFetchFeeds(since time.Time, limit int32) ([]FeedView, error) {
	ctx := context.Background()
	dat, dbErr := cfg.DB.GetNextToFetchFeeds(ctx, database.GetNextToFetchFeedsParams{
		LastFetchedAt: sql.NullTime{Time: since, Valid: true},
		Limit:         limit,
	})
	if dbErr != nil {
		return []FeedView{}, errors.New("database error")
	}

	res := []FeedView{}
	for _, feed := range dat {
		res = append(res, DBFeedToView(&feed))
	}
	return res, nil
}

func (cfg *apiConfig) fetchFeedData(feed FeedView) error {
	// Perform https request at url
	resp, httpsErr := http.Get(feed.Url)
	if httpsErr != nil {
		return fmt.Errorf("error fetching data: %s", httpsErr)
	}
	// Parse XML
	type RSSFeedData struct {
		Title string `json:"title"`
	}
	dat := RSSFeedData{}
	decoder := xml.NewDecoder(resp.Body)
	xmlErr := decoder.Decode(&dat)
	if xmlErr != nil {
		return fmt.Errorf("error decoding xml: %s", xmlErr)
	}
	// Update feed
	ctx := context.Background()
	cfg.DB.UpdateFeedFetchedAt(ctx, uuid.NullUUID{UUID: feed.Id, Valid: true})
	// Process feed
	log.Println(dat.Title)
	return nil
}

func DBFeedToView(feed *database.Feed) FeedView {
	return FeedView{
		Id:            feed.ID.UUID,
		CreatedAt:     timeOrNil(&feed.CreatedAt),
		UpdatedAt:     timeOrNil(&feed.UpdatedAt),
		Name:          feed.Name.String,
		Url:           feed.Url.String,
		LastFetchedAt: timeOrNil(&feed.LastFetchedAt),
	}
}
