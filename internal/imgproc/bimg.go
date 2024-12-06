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

//
// func readImage(data []byte) ([]byte, error) {
// 	return bimg.Read(path)
// }

func (b *BimgProccessor) Resize(data []byte, width, height int) ([]byte, error) {
	return bimg.NewImage(data).Resize(width, height)
}

func (b *BimgProccessor) Rotate(data []byte, angle int) ([]byte, error) {
	return bimg.NewImage(data).Rotate(bimg.Angle(angle))
}

func (b *BimgProccessor) Crop(data []byte, width, height int) ([]byte, error) {
	return bimg.NewImage(data).Crop(width, height, bimg.GravityCentre)
}

func (b *BimgProccessor) Zoom(data []byte, factor int) ([]byte, error) {
	return bimg.NewImage(data).Zoom(factor)
}

func (b *BimgProccessor) Flip(data []byte) ([]byte, error) {
	return bimg.NewImage(data).Flip()
}

func (b *BimgProccessor) Convert(data []byte, imageType string) ([]byte, error) {
	switch imageType {
	case "png":
		return bimg.NewImage(data).Convert(bimg.PNG)
	case "jpeg":
		return bimg.NewImage(data).Convert(bimg.JPEG)
	case "webp":
		return bimg.NewImage(data).Convert(bimg.WEBP)
	case "svg":
		return bimg.NewImage(data).Convert(bimg.SVG)
	default:
		return nil, fmt.Errorf("%s is not a supported file format", imageType)
	}
}
