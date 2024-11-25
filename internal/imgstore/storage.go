package imgstore

import (
	"context"
	"mime/multipart"
)

type Storage interface {
	Upload(ctx context.Context, FileHeader *multipart.FileHeader) (*UploadResponse, error)
	Get(ctx context.Context, fileName string) ([]byte, error)
	Delete(ctx context.Context, fileName string) error
}

type UploadResponse struct {
	FileName   string
	StorageUrl string
	Size       int64
}
