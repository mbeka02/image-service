package server

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mbeka02/image-service/internal/database"
	"github.com/mbeka02/image-service/internal/imgproc"
	"github.com/mbeka02/image-service/internal/imgstore"
	"github.com/mbeka02/image-service/internal/models"
)

type ImageHandler struct {
	Store          *database.Store
	FileStorage    imgstore.Storage
	ImageProcessor imgproc.ImageProcessor
}

func (ih *ImageHandler) processImage(w http.ResponseWriter, r *http.Request,
	fn func(string) ([]byte, error),
	fileName string,
	successMsg string,
) {
	path, err := ih.FileStorage.DownloadTemp(r.Context(), fileName)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
	defer os.RemoveAll(path)
	fileData, err := fn(path)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
	// response := APIResponse{
	// 	Message: successMsg,
	// 	Data:    fileData,
	// 	Status:  http.StatusOK,
	// }
	respondWithImage(w, fileData)
	// respondWithJSON(w, http.StatusOK, response)
}

func (ih *ImageHandler) handleImageUpload(w http.ResponseWriter, r *http.Request) {
	// get the file
	_, fileHeader, err := r.FormFile("image")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Errorf("bad request:%v", err))
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
	idParam := chi.URLParam(r, "imageId")
	imageId, err := strconv.Atoi(idParam)
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

func (ih *ImageHandler) handleImageResize(w http.ResponseWriter, r *http.Request) {
	request := models.ResizeImageRequest{}
	if err := parseAndValidateRequest(r, &request); err != nil {

		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	ih.processImage(w, r, func(path string) ([]byte, error) {
		return ih.ImageProcessor.Resize(path, request.Width, request.Height)
	}, request.FileName, "Image Resized")
}

func (ih *ImageHandler) handleImageRotation(w http.ResponseWriter, r *http.Request) {
	request := models.RotateImageRequest{}
	if err := parseAndValidateRequest(r, &request); err != nil {

		respondWithError(w, http.StatusBadRequest, err)
		return

	}
	ih.processImage(w, r, func(path string) ([]byte, error) {
		return ih.ImageProcessor.Rotate(path, request.Angle)
	}, request.FileName, "Image Rotated")
}

func (ih *ImageHandler) handleImageCropping(w http.ResponseWriter, r *http.Request) {
	request := models.CropImageRequest{}
	if err := parseAndValidateRequest(r, &request); err != nil {

		respondWithError(w, http.StatusBadRequest, err)
		return

	}
	ih.processImage(w, r, func(path string) ([]byte, error) {
		return ih.ImageProcessor.Resize(path, request.Height, request.Width)
	}, request.FileName, "Image Resized")
}

func (ih *ImageHandler) handleImageFlip(w http.ResponseWriter, r *http.Request) {
	request := models.FlipImageRequest{}
	if err := parseAndValidateRequest(r, &request); err != nil {

		respondWithError(w, http.StatusBadRequest, err)
		return

	}
	ih.processImage(w, r, func(path string) ([]byte, error) {
		return ih.ImageProcessor.Flip(path)
	}, request.FileName, "Image Flipped")
}

func (ih *ImageHandler) handleImageConversion(w http.ResponseWriter, r *http.Request) {
	request := models.ConvertImageRequest{}
	if err := parseAndValidateRequest(r, &request); err != nil {

		respondWithError(w, http.StatusBadRequest, err)
		return

	}
	ih.processImage(w, r, func(path string) ([]byte, error) {
		return ih.ImageProcessor.Convert(path, request.ImageType)
	}, request.FileName, "Image Converted")
}

func (ih *ImageHandler) handleImageZoom(w http.ResponseWriter, r *http.Request) {
	request := models.ZoomImageRequest{}
	if err := parseAndValidateRequest(r, &request); err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	ih.processImage(w, r, func(path string) ([]byte, error) {
		return ih.ImageProcessor.Zoom(path, request.Factor)
	}, request.FileName, fmt.Sprint("zoomed by a factor of:", request.Factor))
}
