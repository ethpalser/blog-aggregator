-- +goose Up
ALTER TABLE users ADD CONSTRAINT unq_users_id UNIQUE (id);

CREATE TABLE feeds(
    id UUID UNIQUE,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    name TEXT,
    url TEXT,
    user_id UUID,
    CONSTRAINT fk_feed_user_id
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

-- +goose Down
DROP CONSTRAINT unq_users_id;

DROP TABLE feeds;