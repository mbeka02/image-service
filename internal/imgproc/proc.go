package imgproc

type ImageProcessor interface {
	Resize(path string, width, height int) ([]byte, error)
}
