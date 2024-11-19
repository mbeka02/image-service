package imgproc

import "github.com/h2non/bimg"

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

func (b *BimgProccessor) Resize(path string, width, height int) ([]byte, error) {
	buff, err := bimg.Read(path)
	if err != nil {
		return nil, err
	}

	return bimg.NewImage(buff).Resize(width, height)
}
