-- +goose Up  
CREATE TABLE IF NOT EXISTS images (
image_id bigserial PRIMARY KEY,
user_id bigint  NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,

)
CREATE INDEX ON images(user_id);

-- +goose Down
DROP TABLE images;
