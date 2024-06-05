-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetPost :one
SELECT * FROM posts
WHERE id = $1
ORDER BY published_at DESC
LIMIT $2;

-- name: GetPostsByUser :many
SELECT p.* FROM posts as p
RIGHT JOIN (
    SELECT feed_id
    FROM feed_follows
    WHERE user_id = $1
) as uf
ON p.feed_id = uf.feed_id
ORDER BY published_at DESC
LIMIT $2;