package image

import (
	"bytes"
	"fmt"
	"github.com/ajstarks/svgo"
	"image/color"
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

func GenImage(col, label string) ([]byte, error) {
	width := 46
	height := 46
	buffer := new(bytes.Buffer)
	canvas := svg.New(buffer)
	canvas.Start(width, height, `xlink:href="data:image/png;base64"`)
	rgbaBackCol := getCol(col)
	backCol := canvas.RGBA(int(rgbaBackCol.R), int(rgbaBackCol.G), int(rgbaBackCol.B), float64(rgbaBackCol.A))
	canvas.Rect(0, 0, height, width, fmt.Sprintf("fill:%s", backCol))
	canvas.Text(width/2, height/2, label, "dominant-baseline:middle;text-anchor:middle;font-size:25px;fill:black;font-family:Arial")
	canvas.End()
	return buffer.Bytes(), nil

}
