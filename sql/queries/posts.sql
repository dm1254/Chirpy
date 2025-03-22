-- name: CreatePosts :one 
INSERT INTO posts(id,created_at,updated_at,body,user_id)
values(
	gen_random_uuid(),
	NOW(),
	NOW(),
	$1,
	$2
)
RETURNING *;

-- name: GetPosts :many
SELECT * FROM posts
ORDER BY created_at ASC; 

-- name: GetSinglePost :one
SELECT * FROM posts
WHERE id = $1;


-- name: DeletePost :exec
DELETE FROM posts 
WHERE id = $1;

-- name: GetPostsByAuthor :many
SELECT * FROM posts
WHERE user_id  = $1;
