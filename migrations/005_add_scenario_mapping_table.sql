-- +goose Up
CREATE TABLE scenario_mapping (
    id SERIAL PRIMARY KEY,
    scenario_id INTEGER NOT NULL,
    key VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted BOOLEAN NOT NULL DEFAULT false,
    FOREIGN KEY (scenario_id) REFERENCES scenarios(id),
    FOREIGN KEY (key) REFERENCES blocks(key) 
);

-- +goose Down
DROP TABLE blocks;