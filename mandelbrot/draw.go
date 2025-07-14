package mandelbrot

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math/cmplx"
)

const MAX_ITERATIONS = 1000

var left = -2.0
var right = 1.0
var top = 1.0
var bottom = -1.0

func Pixel(x, y, w, h int) color.Color {
	z := 0 + 0i
	real_part := left + float64(x)/float64(w)*(right-left)
	img_part := top - float64(y)/float64(h)*(top-bottom)
	c := complex(real_part, img_part)

	for range MAX_ITERATIONS {
		if cmplx.Abs(z) >= 2.0 {
			return color.White
		}
		z = cmplx.Pow(z, 2) + c
	}
	return color.Black
}

func DrawImage(
	ctx context.Context,
	img draw.Image,
	gorouintes int,
) (lineDone <-chan struct{}) {
	h := img.Bounds().Size().Y
	w := img.Bounds().Size().X

	pixels := make(chan image.Point)
	lines := make(chan struct{})

	fmt.Println("starting render routines:", gorouintes)

	for range gorouintes {
		go func() {
			for pixel := range pixels {
				img.Set(pixel.X, pixel.Y, Pixel(pixel.X, pixel.Y, w, h))
				if ctx.Err() != nil {
					return
				}
			}
		}()
	}

	go func() {
		for y := range h {
			for x := range w {
				pixels <- image.Point{X: x, Y: y}
				if ctx.Err() != nil {
					fmt.Println("redrawing raster cancelled")
					return
				}
			}
			lines <- struct{}{}
		}
		close(pixels)
		close(lines)
	}()

	return lines
}
