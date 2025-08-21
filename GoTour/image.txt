package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

type Image struct {
	W int
	H int
}

func (i Image) ColorModel() color.Model {
	return color.RGBAModel
}

func (i Image) Bounds() image.Rectangle {
	return image.Rect(0, 0, i.W, i.H)
}

func (i Image) At(x, y int) color.Color {
	v := uint8(x ^ y + x | y)
	return color.RGBA{v, v, 255, 255}
}

func main() {
	m := Image{W: 200, H: 300}

	f, err := os.Create("output.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	png.Encode(f, m)
	fmt.Println("圖像已保存為 output.png")
}
