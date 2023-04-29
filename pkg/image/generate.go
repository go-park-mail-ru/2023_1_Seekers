package image

import (
	"bytes"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/common"
	"github.com/golang/freetype/truetype"
	log "github.com/sirupsen/logrus"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/png"
	"os"
	"sync"
	"unicode"
)

var Colors = make(map[string]color.RGBA)
var TextColors = make(map[color.RGBA]color.RGBA)
var fontTtf *truetype.Font
var fontSize = 20.0

func Init(ttfPath string) {
	Colors["purple"] = color.RGBA{R: 178, G: 102, B: 255, A: 255}
	Colors["green"] = color.RGBA{R: 102, G: 255, B: 102, A: 255}
	Colors["orange"] = color.RGBA{R: 255, G: 153, B: 51, A: 255}
	Colors["red"] = color.RGBA{R: 255, G: 102, B: 102, A: 255}
	Colors["yellow"] = color.RGBA{R: 255, G: 255, B: 102, A: 255}
	Colors["blue"] = color.RGBA{R: 51, G: 153, B: 255, A: 255}

	TextColors[color.RGBA{R: 178, G: 102, B: 255, A: 255}] = color.RGBA{R: 255, G: 250, B: 240, A: 255} //purple
	TextColors[color.RGBA{R: 102, G: 255, B: 102, A: 255}] = color.RGBA{R: 75, G: 0, B: 30, A: 255}     // green
	TextColors[color.RGBA{R: 255, G: 153, B: 51, A: 255}] = color.RGBA{R: 0, G: 0, B: 0, A: 255}        // orange
	TextColors[color.RGBA{R: 255, G: 102, B: 102, A: 255}] = color.RGBA{R: 0, G: 0, B: 0, A: 255}       // red
	TextColors[color.RGBA{R: 255, G: 255, B: 102, A: 255}] = color.RGBA{R: 248, G: 165, B: 29, A: 255}  // yellow
	TextColors[color.RGBA{R: 51, G: 153, B: 255, A: 255}] = color.RGBA{R: 0, G: 0, B: 0, A: 255}        // blue

	fontBytes, err := os.ReadFile(ttfPath)
	if err != nil {
		log.Fatal("cant read font file")
	}

	fontTtf, err = truetype.Parse(fontBytes)
	if err != nil {
		log.Fatal("cant parse ttf")
	}
}

var muCol sync.RWMutex
var muTextCol sync.RWMutex

func GetRandColor() string {
	muCol.RLock()
	defer muCol.RUnlock()
	for k := range Colors {
		return k
	}
	return "green"
}

func getCol(col string) color.RGBA {
	muTextCol.RLock()
	defer muTextCol.RUnlock()
	rgbaCol, ok := Colors[col]
	if !ok {
		return color.RGBA{G: 255, A: 255}
	}
	return rgbaCol
}

func getLabelCol(col color.RGBA) color.RGBA {
	muCol.RLock()
	defer muCol.RUnlock()
	rgbaCol, ok := TextColors[col]
	if !ok {
		return color.RGBA{R: 0, B: 0, G: 0, A: 0}
	}
	return rgbaCol
}

func addLabel(img *image.RGBA, x, y int, face font.Face, label string, labelCol color.RGBA) error {
	point := fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(labelCol),
		Face: face,
		Dot:  point,
	}
	d.DrawString(label)
	return nil

}

func NewPNG(width, height int, col color.RGBA) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	// Set color for each pixel.
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, col)
		}
	}
	return img
}

func GenImage(col, label string, height, width int) ([]byte, error) {
	buffer := new(bytes.Buffer)
	rgbaBgCol := getCol(col)
	img := NewPNG(width, height, rgbaBgCol)

	var labelRune = common.GetFirstRune(label)

	face := truetype.NewFace(fontTtf, &truetype.Options{
		Size:    fontSize,
		DPI:     100.0,
		Hinting: font.HintingNone,
	})
	var runeWidth int
	w, ok := face.GlyphAdvance(labelRune)
	if !ok {
		runeWidth = (width - int(fontSize)) / 2
	} else {
		runeWidth = w.Floor()
	}

	if unicode.IsUpper(labelRune) {
		if err := addLabel(img, (width-runeWidth)/2, height-(height-int(fontSize))/2, face, string(labelRune), getLabelCol(rgbaBgCol)); err != nil {
			return nil, err
		}
	} else {
		if err := addLabel(img, (width-runeWidth)/2, height-(height-int(fontSize))/2-3, face, string(labelRune), getLabelCol(rgbaBgCol)); err != nil {
			return nil, err
		}
	}

	if err := png.Encode(buffer, img); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func UpdateImgText(img []byte, label string, height, width int) ([]byte, error) {
	pngImg, err := png.Decode(bytes.NewReader(img))
	if err != nil {
		return nil, err
	}

	r, g, b, a := pngImg.At(0, 0).RGBA()
	col := color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: uint8(a),
	}

	buffer := new(bytes.Buffer)
	resultPng := NewPNG(width, height, col)

	var labelRune = common.GetFirstRune(label)

	face := truetype.NewFace(fontTtf, &truetype.Options{
		Size:    fontSize,
		DPI:     100.0,
		Hinting: font.HintingNone,
	})
	var runeWidth int
	w, ok := face.GlyphAdvance(labelRune)
	if !ok {
		runeWidth = (width - int(fontSize)) / 2
	} else {
		runeWidth = w.Floor()
	}

	if err = addLabel(resultPng, (width-runeWidth)/2, height-(height-int(fontSize))/2, face, string(labelRune), getLabelCol(col)); err != nil {
		return nil, err
	}
	if unicode.IsUpper(labelRune) {
		if err = addLabel(resultPng, (width-runeWidth)/2, height-(height-int(fontSize))/2, face, string(labelRune), getLabelCol(col)); err != nil {
			return nil, err
		}
	} else {
		if err = addLabel(resultPng, (width-runeWidth)/2, height-(height-int(fontSize))/2-3, face, string(labelRune), getLabelCol(col)); err != nil {
			return nil, err
		}
	}

	if err = png.Encode(buffer, resultPng); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
