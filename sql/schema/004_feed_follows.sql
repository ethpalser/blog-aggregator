-- +goose Up
CREATE TABLE feed_follows(
    user_id UUID,
    feed_id UUID,
    PRIMARY KEY (user_id, feed_id)
);

-- +goose Down
DROP TABLE feed_follows;