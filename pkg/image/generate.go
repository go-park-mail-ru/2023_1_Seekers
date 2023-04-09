package image

import (
	"bytes"
	"github.com/pkg/errors"
	"golang.org/x/image/font"
	"golang.org/x/image/font/inconsolata"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/png"
	"sync"
)

var colors = make(map[string]color.RGBA)

func init() {
	colors["purple"] = color.RGBA{R: 178, G: 102, B: 255, A: 255}
	colors["green"] = color.RGBA{R: 102, G: 255, B: 102, A: 255}
	colors["orange"] = color.RGBA{R: 255, G: 153, B: 51, A: 255}
	colors["red	"] = color.RGBA{R: 255, G: 102, B: 102, A: 255}
	colors["yellow"] = color.RGBA{R: 255, G: 255, B: 102, A: 255}
	colors["blue"] = color.RGBA{R: 51, G: 153, B: 255, A: 255}
}

var mu sync.RWMutex

func GetRandColor() string {
	mu.RLock()
	defer mu.RUnlock()
	for k := range colors {
		return k
	}
	return "green"
}

func getCol(col string) color.RGBA {
	mu.RLock()
	defer mu.RUnlock()
	rgbaCol, ok := colors[col]
	if !ok {
		return color.RGBA{G: 255, A: 255}
	}
	return rgbaCol
}

func addLabel(img *image.RGBA, x, y int, label string) {
	col := color.RGBA{A: 255}
	point := fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: inconsolata.Regular8x16,
		Dot:  point,
	}
	d.DrawString(label)
}

func GenImage(col, label string) ([]byte, error) {
	size := 20
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	backCol := getCol(col)
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			img.Set(i, j, backCol)
		}
	}

	addLabel(img, 6, 15, label)

	buffer := new(bytes.Buffer)
	if err := png.Encode(buffer, img); err != nil {
		return nil, errors.Wrap(err, "failed encode to png")
	}
	return buffer.Bytes(), nil
}
