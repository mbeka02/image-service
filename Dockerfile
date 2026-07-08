# ---- build stage ----
FROM golang:1.23.1-alpine AS builder

# build-base + vips-dev: bimg needs a C compiler and libvips headers to cgo-compile
RUN apk add --no-cache build-base pkgconfig vips-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# CGO must stay enabled — bimg won't build without it
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o main cmd/server/main.go

# ---- run stage ----
FROM alpine:3.21

# vips (runtime .so, not -dev) + ca-certificates for outbound HTTPS (GCS, SMTP, Cloud SQL)
RUN apk add --no-cache vips ca-certificates \
    && addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app
COPY --from=builder /app/main .
USER appuser

EXPOSE 3000
ENTRYPOINT ["/main"]
