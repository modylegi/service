-- +goose Up
CREATE TABLE content_type_dic (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT
);


-- +goose Down
DROP TABLE content_type_dic;
