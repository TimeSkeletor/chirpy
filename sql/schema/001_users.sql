-- +goose Up
CREATE TABLE users (
    id uuid DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    email TEXT NOT NULL
);

-- +goose Down
DROP TABLE users;