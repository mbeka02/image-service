-- name: GetUsers :many
SELECT user_id , full_name , email FROM users LIMIT $1 OFFSET $2;
-- name: GetUserByID :one
SELECT user_id , full_name , email FROM users WHERE user_id=$1;
