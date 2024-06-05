-- +goose Up
CREATE TABLE posts(
    id UUID,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    title TEXT NOT NULL,
    url TEXT NOT NULL,
    description TEXT,
    published_at TIMESTAMP,
    feed_id UUID,
    CONSTRAINT fk_post_feed_id
    FOREIGN KEY (feed_id)
    REFERENCES feeds(id)
);

-- +goose Down
DROP TABLE posts;