-- +goose Up
CREATE TABLE blocks (
    id SERIAL PRIMARY KEY,
    key VARCHAR(255) UNIQUE NOT NULL,  
    title VARCHAR(255) NOT NULL,
    background_image VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted BOOLEAN NOT NULL DEFAULT false
);


-- +goose Down
DROP TABLE blocks;