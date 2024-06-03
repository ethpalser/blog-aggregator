package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ethpalser/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, u database.User) {
	type CreateFollowReq struct {
		FeedId uuid.UUID `json:"feed_id"`
	}
	req := CreateFollowReq{}
	decoder := json.NewDecoder(r.Body)
	jsonErr := decoder.Decode(&req)
	if jsonErr != nil {
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	ctx := context.Background()
	_, dbErr := cfg.DB.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		UserID: u.ID.UUID,
		FeedID: req.FeedId,
	})
	if dbErr != nil {
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondWithJSON(w, http.StatusCreated, "")
}

func (cfg *apiConfig) handlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request, u database.User) {
	type DeleteFollowReq struct {
		FeedId uuid.UUID `json:"feed_id"`
	}
	req := DeleteFollowReq{}
	decoder := json.NewDecoder(r.Body)
	jsonErr := decoder.Decode(&req)
	if jsonErr != nil {
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	ctx := context.Background()
	dat, dbErr := cfg.DB.GetFeedFollowById(ctx, req.FeedId)
	if dbErr != nil {
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if dat.UserID != u.ID.UUID {
		respondWithError(w, http.StatusForbidden, "you do not have permission to this resource")
		return
	}

	_, delErr := cfg.DB.DeleteFeedFollow(ctx, database.DeleteFeedFollowParams{
		UserID: u.ID.UUID,
		FeedID: req.FeedId,
	})
	if delErr != nil {
		respondWithError(w, http.StatusInternalServerError, "an error occurred deleting resource")
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

func (cfg *apiConfig) handlerGetUserFeedFollows(w http.ResponseWriter, r *http.Request, u database.User) {
	ctx := context.Background()
	dat, dbErr := cfg.DB.GetUserFeedFollows(ctx, u.ID.UUID)
	if dbErr != nil {
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	userFeeds := []FeedView{}
	for _, feed := range dat {
		userFeeds = append(userFeeds, FeedView{
			Id:        feed.FeedID,
			CreatedAt: feed.CreatedAt.Time,
			UpdatedAt: feed.UpdatedAt.Time,
			Name:      feed.Name.String,
			Url:       feed.Url.String,
		})
	}

	respondWithJSON(w, http.StatusOK, userFeeds)
}
