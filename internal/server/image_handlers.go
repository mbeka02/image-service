package server

import (
	"errors"
	"fmt"
	"image"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	_ "image/jpeg"
	_ "image/png"

	"github.com/go-chi/chi/v5"
	"github.com/mbeka02/image-service/internal/database"
	"github.com/mbeka02/image-service/internal/imgproc"
	"github.com/mbeka02/image-service/internal/imgstore"
	"github.com/mbeka02/image-service/internal/models"
)

const maxFileSize = 1024 * 1024 * 10

type ImageHandler struct {
	Store          *database.Store
	FileStorage    imgstore.Storage
	ImageProcessor imgproc.ImageProcessor
}

func getImageId(r *http.Request) (int, error) {
	idParam := chi.URLParam(r, "imageId")
	imageId, err := strconv.Atoi(idParam)
	if err != nil {
		return 0, errors.New("invalid url param")
	}
	return imageId, nil
}

func extractMetadata(file multipart.File) (string, int, int, error) {
	buff := make([]byte, 512)
	if _, err := file.Read(buff); err != nil {
		return "", 0, 0, err
	}
	contentType := http.DetectContentType(buff)
	file.Seek(0, 0)

	config, _, err := image.DecodeConfig(file)
	if err != nil {
		return contentType, 0, 0, err
	}

	file.Seek(0, 0)
	fmt.Printf("contentType:%v,width:%v,height:%v", contentType, config.Width, config.Height)
	return contentType, config.Width, config.Height, nil
}

func (ih *ImageHandler) handleImageUpload(w http.ResponseWriter, r *http.Request) {
	// get the file
	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Errorf("bad request:%v", err))
		return
	}
	allowedFileTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
	}
	contentType, width, height, err := extractMetadata(file)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Errorf("uanble to extract metadata:%v", err))
		fmt.Println(err)
		return
	}
	if _, ok := allowedFileTypes[contentType]; !ok {
		respondWithError(w, http.StatusBadRequest, fmt.Errorf("invalid file format:%v", contentType))
		return
	}

	// upload the file to GC storage
	uploadResponse, err := ih.FileStorage.Upload(r.Context(), fileHeader)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Errorf("internal server error : %v", err))
		return
	}
	payload, err := getAuthPayload(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
	// save to DB
	createdImage, err := ih.Store.CreateImage(r.Context(), database.CreateImageParams{
		UserID:     payload.UserID,
		FileName:   uploadResponse.FileName,
		StorageUrl: uploadResponse.StorageUrl,
		FileSize:   uploadResponse.Size,
		Metadata:   "",
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	response := APIResponse{
		Status:  http.StatusOK,
		Data:    createdImage,
		Message: "uploaded",
	}
	respondWithJSON(w, http.StatusOK, response)
}

func (ih *ImageHandler) handleGetImages(w http.ResponseWriter, r *http.Request) {
	payload, err := getAuthPayload(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	// Get limit from query parameter, default to 10 if not provided
	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10 // default limit
	}

	// Get offset from query parameter, default to 0 if not provided
	offsetStr := r.URL.Query().Get("offset")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0 // default offset
	}
	data, err := ih.Store.GetUserImages(r.Context(), database.GetUserImagesParams{
		UserID: payload.UserID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
	response := APIResponse{
		Status:  http.StatusOK,
		Data:    data,
		Message: "images",
	}
	respondWithJSON(w, http.StatusOK, response)
}

func (ih *ImageHandler) handleDeleteImage(w http.ResponseWriter, r *http.Request) {
	imageId, err := getImageId(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, errors.New("invalid url param"))
		return
	}
	payload, err := getAuthPayload(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
	image, err := ih.Store.GetImage(r.Context(), int64(imageId))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, errors.New("unable to get image"))
		return
	}
	if image.UserID != payload.UserID {
		respondWithError(w, http.StatusUnauthorized, errors.New("unauthorized!"))
		return
	}
	if err = ih.FileStorage.Delete(r.Context(), image.FileName); err != nil {
		respondWithError(w, http.StatusInternalServerError, errors.New("unable to delete the image"))
		return
	}

	ih.Store.DeleteUserImage(r.Context(), database.DeleteUserImageParams{
		UserID:  payload.UserID,
		ImageID: int64(imageId),
	})
	response := APIResponse{
		Status:  http.StatusOK,
		Message: "deleted the image sucessfully",
		Data:    nil,
	}
	respondWithJSON(w, http.StatusOK, response)
}

func (ih *ImageHandler) handleGetImage(w http.ResponseWriter, r *http.Request) {
	imageId, err := getImageId(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, errors.New("invalid url param"))
		return
	}
	payload, err := getAuthPayload(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
	image, err := ih.Store.GetImage(r.Context(), int64(imageId))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, errors.New("unable to get image"))
		return
	}
	if image.UserID != payload.UserID {
		respondWithError(w, http.StatusUnauthorized, errors.New("unauthorized!"))
		return
	}
	respondWithJSON(w, http.StatusOK, APIResponse{
		Message: "image:",
		Data:    image,
		Status:  http.StatusOK,
	})
}

func (ih *ImageHandler) handleImageTransformations(w http.ResponseWriter, r *http.Request) {
	imageId, err := getImageId(r)
	image, err := ih.Store.GetImage(r.Context(), int64(imageId))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, errors.New("unable to get image"))
		return
	}
	payload, err := getAuthPayload(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
	if image.UserID != payload.UserID {
		respondWithError(w, http.StatusUnauthorized, errors.New("unauthorized!"))
		return
	}
	request := models.TransformationsRequest{}
	err = parseAndValidateRequest(r, &request)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
	}

	path, err := ih.FileStorage.DownloadTemp(r.Context(), image.FileName)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, APIError{
			Message: "unable to perform the transformations",
			Status:  http.StatusInternalServerError,
			Detail:  err.Error(),
		})
		return
	}
	fileData, err := ih.applyTransformations(path, &request)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, APIError{
			Message: "unable to perform the transformations",
			Status:  http.StatusInternalServerError,
			Detail:  err.Error(),
		})
		return
	}
	respondWithImage(w, fileData)
}

func (ih *ImageHandler) applyTransformations(imagePath string, request *models.TransformationsRequest) ([]byte, error) {
	var err error
	var currentImageData []byte = nil
	// Read the initial image
	currentImageData, err = os.ReadFile(imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read initial image: %v", err)
	}

	// Apply transformations in a specific order
	transformationFuncs := []func() ([]byte, error){
		func() ([]byte, error) {
			if request.Resize != nil {
				return ih.ImageProcessor.Resize(currentImageData, request.Resize.Width, request.Resize.Height)
			}
			return currentImageData, nil
		},
		func() ([]byte, error) {
			if request.Rotate != nil {
				return ih.ImageProcessor.Rotate(currentImageData, request.Rotate.Angle)
			}

			return currentImageData, nil
		},
		func() ([]byte, error) {
			if request.Crop != nil {
				return ih.ImageProcessor.Crop(currentImageData, request.Crop.Width, request.Crop.Height)
			}

			return currentImageData, nil
		},
		func() ([]byte, error) {
			if request.Flip != nil {
				return ih.ImageProcessor.Flip(currentImageData)
			}
			return currentImageData, nil
		},
		func() ([]byte, error) {
			if request.Convert != nil {
				return ih.ImageProcessor.Convert(currentImageData, request.Convert.ImageType)
			}
			return currentImageData, nil
		},
		func() ([]byte, error) {
			if request.Zoom != nil {
				return ih.ImageProcessor.Zoom(currentImageData, request.Zoom.Factor)
			}
			return currentImageData, nil
		},
	}

	defer os.Remove(imagePath)
	// Apply transformations sequentially
	for _, transformFunc := range transformationFuncs {

		// Apply transformation
		currentImageData, err = transformFunc()
		if err != nil {
			return nil, fmt.Errorf("transformation failed: %v", err)
		}
	}

	return currentImageData, nil
}
