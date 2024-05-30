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

type UserView struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type CreateUserReq struct {
		Name string `json:"name"`
	}
	req := CreateUserReq{}
	decoder := json.NewDecoder(r.Body)
	dErr := decoder.Decode(&req)
	if dErr != nil {
		respondWithError(w, http.StatusInternalServerError, "Invalid request structure")
		return
	}

	ctx := context.Background()
	now := time.Now()
	user, dbErr := cfg.DB.CreateUser(ctx, database.CreateUserParams{
		ID:        uuid.NullUUID{UUID: uuid.New(), Valid: true},
		CreatedAt: sql.NullTime{Time: now, Valid: true},
		UpdatedAt: sql.NullTime{Time: now, Valid: true},
		Name:      sql.NullString{String: req.Name, Valid: true},
	})
	if dbErr != nil {
		log.Printf("Error creating user")
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	data := UserView{
		Id:        user.ID.UUID,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
		Name:      user.Name.String,
	}
	respondWithJSON(w, http.StatusOK, data)
}
