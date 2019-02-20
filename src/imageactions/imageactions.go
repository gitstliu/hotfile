package imageactions

import (
	"io"

	goimage "github.com/hunterhug/go_image"
)

func ImageToSpecificSize(reader io.Reader, width int, height int) ([]byte, error) {
	return goimage.ThumbnailS2B(reader, width, height)
}
