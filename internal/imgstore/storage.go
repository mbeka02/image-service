package imgstore

import (
	"context"
	"io"
	"mime/multipart"
)

type Storage interface {
	Upload(ctx context.Context, FileHeader *multipart.FileHeader) (*UploadResponse, error)
	Get(ctx context.Context, fileName string) (io.Reader, error)
	Delete(ctx context.Context, fileName string) error
	DownloadTemp(ctx context.Context, fileName string) (string, error)
}

type UploadResponse struct {
	FileName   string
	StorageUrl string
	Size       int64
}
