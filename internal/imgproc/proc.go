package imgproc

type ImageProcessor interface {
	Resize(data []byte, width, height int) ([]byte, error)
	Rotate(data []byte, angle int) ([]byte, error)
	Crop(data []byte, width, height int) ([]byte, error)
	Zoom(data []byte, factor int) ([]byte, error)
	Flip(data []byte) ([]byte, error)
	Convert(data []byte, imageType string) ([]byte, error)
}
