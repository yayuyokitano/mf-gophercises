-- +migrate Up
CREATE TABLE urls (
    short_url TEXT PRIMARY KEY,
    target_url TEXT NOT NULL
);

-- +migrate Down
DROP TABLE urls;
