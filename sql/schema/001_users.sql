-- +goose Up 
CREATE TABLE IF NOT EXISTS users (
user_id  bigserial PRIMARY KEY,
full_name varchar NOT NULL,
password varchar NOT NULL,
email varchar UNIQUE NOT NULL,
created_at timestamptz NOT NULL DEFAULT (now()),
password_changed_at timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z'
);

-- +goose Down
DROP TABLE users;
