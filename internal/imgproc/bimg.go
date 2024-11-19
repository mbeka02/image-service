package imgproc

import (
	"fmt"

	"github.com/h2non/bimg"
)

type BimgProccessor struct {
	ProcessorOptions bimg.Options
}

func NewBimgProcessor(quality, compression int) ImageProcessor {
	options := bimg.Options{
		Quality:     quality,
		Compression: compression,
	}
	return &BimgProccessor{
		options,
	}
}

func readImage(path string) ([]byte, error) {
	return bimg.Read(path)
}

func (b *BimgProccessor) Resize(path string, width, height int) ([]byte, error) {
	buff, err := readImage(path)
	if err != nil {
		return nil, err
	}

	return bimg.NewImage(buff).Resize(width, height)
}

func (b *BimgProccessor) Rotate(path string, angle int) ([]byte, error) {
	buff, err := readImage(path)
	if err != nil {
		return nil, err
	}
	return bimg.NewImage(buff).Rotate(bimg.Angle(angle))
}

func (b *BimgProccessor) Crop(path string, width, height int) ([]byte, error) {
	buff, err := readImage(path)
	if err != nil {
		return nil, err
	}
	return bimg.NewImage(buff).Crop(width, height, bimg.GravityCentre)
}

func (b *BimgProccessor) Zoom(path string, factor int) ([]byte, error) {
	buff, err := readImage(path)
	if err != nil {
		return nil, err
	}
	return bimg.NewImage(buff).Zoom(factor)
}

func (b *BimgProccessor) Flip(path string) ([]byte, error) {
	buff, err := readImage(path)
	if err != nil {
		return nil, err
	}
	return bimg.NewImage(buff).Flip()
}

func (b *BimgProccessor) Convert(path, imageType string) ([]byte, error) {
	buff, err := readImage(path)
	if err != nil {
		return nil, err
	}
	switch imageType {
	case "png":
		return bimg.NewImage(buff).Convert(bimg.PNG)
	case "jpeg":
		return bimg.NewImage(buff).Convert(bimg.JPEG)
	case "webp":
		return bimg.NewImage(buff).Convert(bimg.WEBP)
	case "svg":
		return bimg.NewImage(buff).Convert(bimg.SVG)
	default:
		return nil, fmt.Errorf("%s is not a supported file format", imageType)
	}
}
