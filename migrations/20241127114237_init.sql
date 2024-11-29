-- +goose Up
-- +goose StatementBegin
CREATE TABLE songs (
    id UUID PRIMARY KEY,
    group_name TEXT NOT NULL,
    title TEXT NOT NULL,
    release_date DATE,
    text TEXT,
    link TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE songs;
-- +goose StatementEnd