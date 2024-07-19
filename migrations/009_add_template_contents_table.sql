-- +goose Up
CREATE TABLE template_contents (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    template_content JSONB NOT NULL,
    content_type_id INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted BOOLEAN NOT NULL DEFAULT false,
    FOREIGN KEY (content_type_id) REFERENCES content_type_dic(id)
);


-- +goose Down
DROP TABLE template_contents;
