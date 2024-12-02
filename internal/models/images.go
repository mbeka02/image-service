package models

type ResizeImageRequest struct {
	Width    int    `json:"width" validate:"required"`
	Height   int    `json:"height" validate:"required"`
	FileName string `json:"file_name" validate:"required"`
}

type RotateImageRequest struct {
	Angle    int    `json:"angle" validate:"required"`
	FileName string `json:"file_name" validate:"required"`
}

type CropImageRequest struct {
	Width    int    `json:"width" validate:"required"`
	Height   int    `json:"height" validate:"required"`
	FileName string `json:"file_name" validate:"required"`
}

type FlipImageRequest struct {
	FileName string `json:"file_name" validate:"required"`
}

type ConvertImageRequest struct {
	ImageType string `json:"image_type" validate:"required"`
	FileName  string `json:"file_name" validate:"required"`
}
