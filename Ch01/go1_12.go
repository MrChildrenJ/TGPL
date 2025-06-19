package main

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
)

var palette = []color.Color{color.RGBA{245, 40, 145, 1}, color.RGBA{39, 245, 163, 1}}

const (
	whiteIndex = 0 // first color in palette
	blackIndex = 1 // next color in palette
)

type LissajousConfig struct {
	Cycles  float64 // number of complete x oscillator revolutions
	Res     float64 // angular resolution
	Size    int     // image canvas covers [-size..+size]
	Nframes int     // number of animation frames
	Delay   int     // delay between frames in 10ms units
}

func defaultConfig() *LissajousConfig {		// return address to a (LissajousConfig) struct
	return &LissajousConfig{
		Cycles:  5,
		Res:     0.001,
		Size:    100,
		Nframes: 64,
		Delay:   8,
	}
}

func main() {
	handler := func(w http.ResponseWriter, r* http.Request) {
		config := parseURLParams(r)
		lissajous(w, config)
	}

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8003", nil))
}

func parseURLParams(r* http.Request) *LissajousConfig {
	config := defaultConfig()

	query := r.URL.Query()

	if cycleStr := query.Get("cycles"); cycleStr != "" {
		cycles, err := strconv.ParseFloat(cycleStr, 64)
		if err != nil && cycles > 0 {
			config.Cycles = cycles
		}
	}
	return config
}

func lissajous(out io.Writer, config *LissajousConfig) {
	freq := rand.Float64() * 3.0 // relative frequency of y oscillator
	anim := gif.GIF{LoopCount: config.Nframes}
	phase := 0.0 // phase difference

	for i := 0; i < config.Nframes; i++ {
		rect := image.Rect(0, 0, 2*config.Size+1, 2*config.Size+1)
		img := image.NewPaletted(rect, palette)

		for t := 0.0; t < config.Cycles*2*math.Pi; t += config.Res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(
				config.Size+int(x*float64(config.Size)+0.5), config.Size+int(y*float64(config.Size)+0.5),
				blackIndex)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, config.Delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim) // NOTE: ignoring encoding errors
}

