package image

import (
	"bytes"
	"fmt"
	"github.com/ajstarks/svgo"
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	"golang.org/x/net/html"
	"image/color"
	"sync"
)

var Colors = make(map[string]color.RGBA)

func init() {
	Colors["purple"] = color.RGBA{R: 178, G: 102, B: 255, A: 255}
	Colors["green"] = color.RGBA{R: 102, G: 255, B: 102, A: 255}
	Colors["orange"] = color.RGBA{R: 255, G: 153, B: 51, A: 255}
	Colors["red"] = color.RGBA{R: 255, G: 102, B: 102, A: 255}
	Colors["yellow"] = color.RGBA{R: 255, G: 255, B: 102, A: 255}
	Colors["blue"] = color.RGBA{R: 51, G: 153, B: 255, A: 255}
}

var mu sync.RWMutex

func GetRandColor() string {
	mu.RLock()
	defer mu.RUnlock()
	for k := range Colors {
		return k
	}
	return "green"
}

func getCol(col string) color.RGBA {
	mu.RLock()
	defer mu.RUnlock()
	rgbaCol, ok := Colors[col]
	if !ok {
		return color.RGBA{G: 255, A: 255}
	}
	return rgbaCol
}

func GenImage(col, label string) ([]byte, error) {
	width := config.UserDefaultAvatarSize
	height := config.UserDefaultAvatarSize
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
func getBackgoundStyle(img []byte) string {
	buffer := new(bytes.Buffer)
	buffer.Write(img)
	doc, _ := html.Parse(buffer)
	var f func(*html.Node)
	var resStyle string
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "rect" {
			for _, s := range n.Attr {
				if s.Key == "style" {
					resStyle = s.Val
					return
				}
			}
		}

		for at := n.FirstChild; at != nil; at = at.NextSibling {
			f(at)
		}
	}

	f(doc)
	return resStyle
}
func UpdateImgText(img []byte, label string) ([]byte, error) {
	style := getBackgoundStyle(img)
	width := config.UserDefaultAvatarSize
	height := config.UserDefaultAvatarSize
	buffer := new(bytes.Buffer)
	canvas := svg.New(buffer)
	canvas.Start(width, height, `xlink:href="data:image/png;base64"`)
	canvas.Rect(0, 0, height, width, style)
	canvas.Text(width/2, height/2, label, "dominant-baseline:middle;text-anchor:middle;font-size:25px;fill:black;font-family:Arial")
	canvas.End()
	return buffer.Bytes(), nil
}
