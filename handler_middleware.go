package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/ethpalser/blog-aggregator/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get("Authorization")
		if authToken == "" {
			respondWithError(w, http.StatusUnauthorized, "Invalid auth token")
			return
		}

		apikey := strings.TrimPrefix(authToken, "ApiKey ")

		ctx := context.Background()
		user, dbErr := cfg.DB.GetUserByApiKey(ctx, apikey)
		if dbErr != nil {
			log.Printf("Error fetching user by api key: %s", dbErr.Error())
			respondWithError(w, http.StatusInternalServerError, "Internal server error")
			return
		}

		handler(w, r, user)
	}
}
