-- +goose Up
-- +goose StatementBegin
CREATE TABLE quotes (
    id SERIAL,
    quote varchar(255) NOT NULL,
    created_at timestamp DEFAULT NULL,
    updated_at timestamp DEFAULT NULL
);

INSERT INTO quotes VALUES
(0, 'amogus', '2022-11-08 01:01:01', '2022-11-08 01:01:01'),
(1, 'Mada kono sekai wa', '2022-11-08 01:01:01', '2022-11-08 01:01:01');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE quotes;
-- +goose StatementEnd
