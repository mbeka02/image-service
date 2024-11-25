-- +goose Up  
CREATE TABLE IF NOT EXISTS images (
image_id bigserial PRIMARY KEY,
user_id bigint  NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
file_name varchar NOT NULL,
file_size bigint NOT NULL,
storage_url varchar NOT NULL,
metadata jsonb,
created_at timestamptz NOT NULL DEFAULT (now()),
updated_at timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z'

);
CREATE INDEX ON images(user_id);

-- +goose Down
DROP TABLE images;
