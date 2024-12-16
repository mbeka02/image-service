# Image Processing Service

This is an API for a simple image processing service written in Golang - https://roadmap.sh/projects/image-processing-service

## Features

- User registration.
- JWT Authentication.
- Email verification (PENDING).
- Image uploading and downloading to and from Google Cloud Storage
- Image transformation using the [h2non/bimg](https://github.com/h2non/bimg) package
- Postgres for data storage

## Installation

1. Clone the repository
2. Configure the environment variables in the `.env` file;
3. Execute DB migrations using goose or any other migration tool.

## Usage

1. Run the project using the `air` command or `make watch` if make is installed.

## Dependencies

### Go packages

- [h2non/bimg](https://github.com/h2non/bimg)
- [chi](https://github.com/go-chi/chi)
- [jwt](https://github.com/golang-jwt/jwt)
- [uuid](https://github.com/google/uuid)
- [pq](https://github.com/lib/pq)
- [sqlc](https://github.com/sqlc-dev/sqlc)
- [goose](https://github.com/pressly/goose)

### External services

- [PostgreSQL](https://www.postgresql.org/)
- [Google Cloud Storage](https://cloud.google.com/?hl=en)

## NOTE:

- This project is primarily a learning exercise and is not intended for commercial use.
