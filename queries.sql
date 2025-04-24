-- name: GetUser :one
SELECT * FROM users WHERE username = $1;

-- name: CreateUser :one
INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING *;

-- name: GetUserItems :many
SELECT * FROM items WHERE user_id = $1;

-- name: CreateItem :one
INSERT INTO items (user_id, value) VALUES ($1, $2) RETURNING *;

-- name: GetItemsContaining :many
SELECT * FROM items WHERE value LIKE '%' || $1 || '%' and user_id = $2;

-- name: DeleteItem :exec
DELETE FROM items WHERE id = $1;

-- name: GetItem :one
SELECT * FROM items WHERE id = $1;