-- name: CreateUsers :one
INSERT INTO users (id,created_at,updated_at,email,hashed_password)
values(
	gen_random_uuid(),
	NOW(),
	NOW(),
	$1,
	$2
)
RETURNING *; 

-- name: Reset :exec 
DELETE FROM users;

-- name: GetUserPass :one
SELECT * FROM users
WHERE email = $1;

-- name: UpdateUserEmailAndPass :one
UPDATE users 
SET email = $1, hashed_password = $2
RETURNING *;

-- name: UpgradeUserToRed :exec
UPDATE users 
SET is_chirpy_red = true
WHERE id = $1;
