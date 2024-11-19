package imgproc

type ImageProcessor interface {
	Resize(path string, width, height int) ([]byte, error)
	Rotate(path string, angle int) ([]byte, error)
	Crop(path string, width, height int) ([]byte, error)
	Zoom(path string, factor int) ([]byte, error)
	Flip(path string) ([]byte, error)
	Convert(path, imageType string) ([]byte, error)
}
