-- name: CreateImage :one
INSERT INTO images(user_id , file_name , file_size , storage_url) VALUES ($1,$2,$3,$4)RETURNING file_name , file_size , storage_url;
