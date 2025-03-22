-- +goose Up
CREATE TABLE users(
	id UUID PRIMARY KEY,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL, 
	email TEXT UNIQUE NOT NULL
);

-- +goose Down 
DROP TABLE users;

-- +goose Up
ALTER TABLE users
ADD COLUMN is_chirpy_red BOOLEAN
DEFAULT false;
