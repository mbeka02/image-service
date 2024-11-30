package models

type ResizeImageRequest struct {
	Width    int    `json:"width" validate:"required"`
	Height   int    `json:"height" validate:"required"`
	FileName string `json:"file_name" validate:"required"`
}
