-- +goose Up
ALTER TABLE users
DROP COLUMN IF EXISTS is_chirpy_red;


ALTER TABLE users
ADD COLUMN is_chirpy_red BOOLEAN NOT NULL
DEFAULT false;

-- +goose Down
ALTER TABLE users
DROP COLUMN IF EXISTS is_chirpy_red;
