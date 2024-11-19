package imgproc

type ImageProcessor interface {
	Resize(path string, width, height int) ([]byte, error)
	Rotate(path string, angle int) ([]byte, error)
	Convert(path, imageType string) ([]byte, error)
}
