package imgstore

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"time"

	"cloud.google.com/go/storage"
)

type GCStorage struct {
	client     *storage.Client
	bucketName string
	projectId  string
}

func NewGCStorage(projectId, bucketName string) (Storage, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to setup the storage client: %v", err)
	}
	return &GCStorage{
		client,
		bucketName,
		projectId,
	}, nil
}

func (g *GCStorage) Upload(ctx context.Context, fileHeader *multipart.FileHeader) (*UploadResponse, error) {
	// open the associated file
	srcFile, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("unable to open the file:%v", err)
	}

	defer srcFile.Close()
	// create a unique filename
	fileName := fmt.Sprintf("%s_%d", fileHeader.Filename, time.Now().UnixNano())

	// get the bucket handle
	bucket := g.client.Bucket(g.bucketName)
	objectHandle := bucket.Object(fileName)

	writer := objectHandle.NewWriter(ctx)
	writer.ContentType = fileHeader.Header.Get("Content-Type")

	// Copy the file to the Object
	written, err := io.Copy(writer, srcFile)
	if err != nil {
		return nil, fmt.Errorf("unable to copy to storage:%v", err)
	}
	defer writer.Close()
	// make the uploaded images public for Now
	/*	if err := objectHandle.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return "", fmt.Errorf("unable to make the file public:%v", err)
	}*/
	return &UploadResponse{
		FileName:   fileName,
		Size:       written,
		StorageUrl: fmt.Sprintf("https://storage.googleapis.com/%s/%s", g.bucketName, fileName),
	}, nil
}

func (g *GCStorage) Download(ctx context.Context, fileName string) (io.Reader, error) {
	object := g.client.Bucket(g.bucketName).Object(fileName)
	reader, err := object.NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to create a new reader:%v", err)
	}

	defer reader.Close()

	return reader, nil
}

func (g *GCStorage) Delete(ctx context.Context, fileName string) error {
	object := g.client.Bucket(g.bucketName).Object(fileName)

	if err := object.Delete(ctx); err != nil {
		return fmt.Errorf("unable to delete the file:%v", err)
	}
	return nil
}

func (g *GCStorage) DownloadTemp(ctx context.Context, fileName string) (string, error) {
	// Get the file
	fileData, err := g.Download(ctx, fileName)
	if err != nil {
		return "", fmt.Errorf("unable to get the file:%v", err)
	}

	// Create a temporary file
	tempFile, err := os.CreateTemp("", fileName)
	if err != nil {
		return "", fmt.Errorf("unable to save the file locally:%v", err)
	}

	defer tempFile.Close()
	if _, err = io.Copy(tempFile, fileData); err != nil {
		return "", fmt.Errorf("unable to copy the file contents")
	}

	return tempFile.Name(), nil
}

func (g *GCStorage) Close() error {
	return g.client.Close()
}
