package models

type ResizeImageRequest struct {
	Width  int `json:"width" validate:"required"`
	Height int `json:"height" validate:"required"`
}

type RotateImageRequest struct {
	Angle int `json:"angle" validate:"required"`
}

type CropImageRequest struct {
	Width  int `json:"width" validate:"required"`
	Height int `json:"height" validate:"required"`
}

type ConvertImageRequest struct {
	ImageType string `json:"image_type" validate:"required"`
}
type ZoomImageRequest struct {
	Factor int `json:"factor" validate:"required"`
}
type TransformationsRequest struct {
	Resize  *ResizeImageRequest  `json:"resize,omitempty"`
	Crop    *CropImageRequest    `json:"crop,omitempty"`
	Rotate  *RotateImageRequest  `json:"rotate,omitempty"`
	Zoom    *ZoomImageRequest    `json:"zoom,omitempty"`
	Convert *ConvertImageRequest `json:"convert,omitempty"`
	Flip    *bool                `json:"flip,omitempty"`
}
