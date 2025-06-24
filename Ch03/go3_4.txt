package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
)

const (
	cells   = 100
	xyrange = 30.0 // axis ranges (-xyrange..+xyrange)
	angle 	= math.Pi / 6 // angle of x, y axes (=30°)
)

func main() {
	http.HandleFunc("/", handler) // each request calls handler
	http.HandleFunc("/image/svg+xml", surfaceHandler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

// handler echoes the Path component of the requested URL.
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}

func surfaceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")

	// Parse query parameters with defaults
	height := parseIntParam(r, "height", 320)
	width := parseIntParam(r, "width", 600)
	color := r.URL.Query().Get("color")
	function := r.URL.Query().Get("function")

	if color == "" {
		color = "gradient" // default to gradient coloring
	}
	cellCount := parseIntParam(r, "cells", cells)

	// Validate parameters
	if height <= 0 || height > 2000 {
		height = 320
	}
	if width <= 0 || width > 2000 {
		width = 600
	}
	if cellCount <= 0 || cellCount > 200 {
		cellCount = cells
	}

	surface(w, height, width, color, cellCount, function)
}

func surface(out io.Writer, height, width int, colorMode string, cells int, function string) {
	fmt.Fprintf(out, "<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d'>", width, height)
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay, az, okA := corner(height, width, i+1, j, function)
			bx, by, bz, okB := corner(height, width, i, j, function)
			cx, cy, cz, okC := corner(height, width, i, j+1, function)
			dx, dy, dz, okD := corner(height, width, i+1, j+1, function)

			if okA && okB && okC && okD {
				avgZ := (az + bz + cz + dz) / 4
				var polygonColor string

				switch colorMode {
				case "gradient":
					polygonColor = heightToColor(avgZ)
				case "red":
					polygonColor = "#ff0000"
				case "blue":
					polygonColor = "#0000ff"
				case "green":
					polygonColor = "#00ff00"
				default:
					if isValidHexColor(colorMode) {
						polygonColor = colorMode
					} else {
						polygonColor = heightToColor(avgZ)
					}
				}

				fmt.Fprintf(out, "<polygon points='%g,%g %g,%g %g,%g %g,%g' fill='%s'/>\n",
					ax, ay, bx, by, cx, cy, dx, dy, polygonColor)
			}
		}
	}
	fmt.Fprintln(out, "</svg>")
}

func corner(height, width, i, j int, function string) (float64, float64, float64, bool) {
	xyscale := float64(width) / 2 / xyrange // pixels per x or y unit
	zscale := float64(height) * 0.4         // pixels per z unit

	var sin30, cos30 = math.Sin(angle), math.Cos(angle) // sin(30°), cos(30°)

	// Find point (x,y) at corner of cell (i,j).
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)


	// Compute surface height z.
	z := f(x, y, function)

	// Checking z
	if math.IsInf(z, 1) || math.IsInf(z, -1) || math.IsNaN(z) {
		return 0, 0, 0, false
	}

	// Project (x,y,z) isometrically onto 2-D SVG canvas (sx,sy).
	sx := float64(width)/2 + (x-y)*cos30*xyscale
	sy := float64(height)/2 + (x+y)*sin30*xyscale - z*zscale

	// Checking sx, sy
	if math.IsInf(sx, 0) || math.IsInf(sy, 0) || math.IsNaN(sx) || math.IsNaN(sy) {
		return 0, 0, 0, false
	}

	return sx, sy, z, true
}

func f(x, y float64, function string) float64 {
	r := math.Hypot(x, y) // distance from (0,0)
	
	switch function {
	case "Y0":
		return math.Y0(r)
	case "Y1":
		return math.Y1(r)
	case "J0":
		return math.J0(r)
	case "J1":
		return math.J1(r)
	default:
		return math.Sin(r) / r
	}
}

func heightToColor(z float64) string {
	minZ := -0.3
	maxZ := 1.0

	normalized := (z - minZ) / (maxZ - minZ)

	if normalized < 0 {
		normalized = 0
	} else if normalized > 1 {
		normalized = 1
	}

	red := int(normalized * 255)
	blue := int((1 - normalized) * 255)
	green := 0

	return fmt.Sprintf("#%02x%02x%02x", red, green, blue) // SPrintf return string i/o output to stdout
	// (255, 0, 0) -> "#ff0000"
}

func parseIntParam(r *http.Request, param string, defaultValue int) int {
	valueStr := r.URL.Query().Get(param)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func isValidHexColor(color string) bool {
	if len(color) != 7 || color[0] != '#' {
		return false
	}

	for i := 1; i < 7; i++ {
		c := color[i]
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}

/*
http://localhost:8000/image/svg+xml
http://localhost:8000/

Size Variations:
http://localhost:8000/image/svg+xml?width=800&height=600
http://localhost:8000/image/svg+xml?width=1200&height=800
http://localhost:8000/image/svg+xml?width=400&height=300

Color Options:
http://localhost:8000/image/svg+xml?color=red
http://localhost:8000/image/svg+xml?color=blue
http://localhost:8000/image/svg+xml?color=green
http://localhost:8000/image/svg+xml?color=gradient
http://localhost:8000/image/svg+xml?color=#ff6600

Grid Resolution:
http://localhost:8000/image/svg+xml?cells=50
http://localhost:8000/image/svg+xml?cells=150
http://localhost:8000/image/svg+xml?cells=25

Combined Parameters:
http://localhost:8000/image/svg+xml?width=1000&height=700&color=blue&cells=80
http://localhost:8000/image/svg+xml?width=600&height=400&color=#00ff88&cells=120
http://localhost:8000/image/svg+xml?width=800&height=600&color=gradient&cells=60

Bessel Functions:
http://localhost:8000/image/svg+xml?function=Y0
http://localhost:8000/image/svg+xml?function=Y1
http://localhost:8000/image/svg+xml?function=J0
http://localhost:8000/image/svg+xml?function=J1

Combined with other parameters:
http://localhost:8000/image/svg+xml?function=Y0&color=blue&width=800
http://localhost:8000/image/svg+xml?function=J1&color=gradient&cells=80
http://localhost:8000/image/svg+xml?function=Y1&color=#ff6600&width=1000&height=700
  */