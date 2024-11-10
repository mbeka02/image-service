-- name: GetUsers :many
SELECT user_id , full_name , email FROM users LIMIT $1 OFFSET $2;
-- name: GetUserByID :one
SELECT user_id , full_name , email FROM users WHERE user_id=$1;
-- name: CreateUser :one
INSERT INTO users(full_name , password , email) VALUES ($1,$2,$3) RETURNING *;
