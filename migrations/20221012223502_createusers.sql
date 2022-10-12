-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id SERIAL,
    username varchar(255) NOT NULL UNIQUE,
    password varchar(255) NOT NULL,
    is_admin BOOLEAN DEFAULT FALSE,
    created_at timestamp DEFAULT NULL,
    updated_at timestamp DEFAULT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
