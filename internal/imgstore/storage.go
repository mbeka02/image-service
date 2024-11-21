package imgstore

import "mime/multipart"

type Storage interface {
	Upload(file *multipart.FileHeader) (string, error)
}
