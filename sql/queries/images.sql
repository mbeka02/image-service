-- name: CreateImage :one
INSERT INTO images(user_id , file_name , file_size , storage_url) VALUES ($1,$2,$3,$4)RETURNING file_name , file_size , storage_url;

-- name: GetUserImages :many
SELECT * FROM images WHERE user_id=$1 LIMIT $2 OFFSET $3;
-- name: GetImage :one
SELECT * FROM images WHERE image_id=$1;
-- name: DeleteUserImage :exec
DELETE FROM images WHERE image_id=$1 AND user_id=$2; 
