package imgstore

import (
	"context"
	"mime/multipart"
)

type Storage interface {
	Upload(ctx context.Context, file *multipart.FileHeader) (string, error)
	Get(ctx context.Context, fileName string) ([]byte, error)
	Delete(ctx context.Context, fileName string) error
}
