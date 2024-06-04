-- name: CreateFeed :one
INSERT INTO feeds(id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetAllFeeds :many
SELECT * FROM feeds;

-- name: GetNextToFetchFeeds :many
SELECT * FROM feeds
WHERE (last_fetched_at IS NULL) OR (last_fetched_at <= $1)
ORDER BY last_fetched_at ASC
LIMIT $2;

-- name: UpdateFeedFetchedAt :one
UPDATE feeds
SET last_fetched_at = now(), updated_at = now()
WHERE id = $1
RETURNING *;