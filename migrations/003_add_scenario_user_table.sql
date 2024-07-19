-- +goose Up
CREATE TABLE scenario_user (
    id SERIAL PRIMARY KEY,
    scenario_id INTEGER  NOT NULL,
    user_id INTEGER  NOT NULL,
    deleted BOOLEAN NOT NULL DEFAULT false,
    FOREIGN KEY (scenario_id) REFERENCES scenarios(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);


-- +goose Down
DROP TABLE scenario_user;