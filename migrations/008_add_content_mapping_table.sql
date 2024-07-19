-- +goose Up
CREATE TABLE content_mapping (
    id SERIAL PRIMARY KEY,
    content_block_id INTEGER NOT NULL,
    content_id INTEGER NOT NULL,
    rating INTEGER,
    deleted BOOLEAN NOT NULL DEFAULT false,
    FOREIGN KEY (content_block_id) REFERENCES blocks(id),
    FOREIGN KEY (content_id) REFERENCES contents(id)
);

-- +goose Down
DROP TABLE content_mapping;