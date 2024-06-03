-- name: CreateFeedFollow :one
INSERT INTO feed_follows (user_id, feed_id)
VALUES ($1, $2)
RETURNING *;

-- name: DeleteFeedFollow :one
DELETE FROM feed_follows
WHERE user_id = $1 AND feed_id = $2
RETURNING *;

-- name: GetUserFeedFollows :many
SELECT * FROM feed_follows as ff
LEFT JOIN feeds as fd
ON ff.feed_id = fd.id
WHERE ff.user_id = $1
ORDER BY fd.updated_at;

-- name: GetFeedFollowById :one
SELECT * FROM feed_follows
WHERE feed_id = $1;