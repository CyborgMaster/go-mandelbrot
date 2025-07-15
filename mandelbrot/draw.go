package mandelbrot

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math/cmplx"

	"github.com/lucasb-eyer/go-colorful"
)

const MAX_ITERATIONS = 1000

func Pixel(point image.Point, size image.Point, bounds Bounds) color.Color {
	z := 0 + 0i
	c := complex(bounds.PixelOffset(point, size))

	for i := range MAX_ITERATIONS {
		if cmplx.Abs(z) >= 2.0 {
			return colorful.Hsv(float64(i%360), 1, 0.5)
		}
		z = cmplx.Pow(z, 2) + c
	}
	return color.Black
}

func DrawImage(
	ctx context.Context,
	img draw.Image,
	bounds Bounds,
	gorouintes int,
) (lineDone <-chan struct{}) {
	size := img.Bounds().Size()
	pixels := make(chan image.Point)
	lines := make(chan struct{})

	fmt.Println("starting render routines:", gorouintes)

	for range gorouintes {
		go func() {
			for pixel := range pixels {
				img.Set(pixel.X, pixel.Y, Pixel(pixel, size, bounds))
				if ctx.Err() != nil {
					return
				}
			}
		}()
	}

	go func() {
		for y := range size.Y {
			for x := range size.X {
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
