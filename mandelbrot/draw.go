package mandelbrot

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math/cmplx"

	"fyne.io/fyne/v2"
	"github.com/lucasb-eyer/go-colorful"
)

const MAX_ITERATIONS = 1000

type Bounds struct {
	Left   float64
	Right  float64
	Top    float64
	Bottom float64
}

func (b Bounds) Width() float64 {
	return b.Right - b.Left
}
func (b Bounds) Height() float64 {
	return b.Top - b.Bottom
}
func (b Bounds) PixelOffset(point image.Point, size image.Point) (float64, float64) {
	x := b.Left + float64(point.X)/float64(size.X)*(b.Width())
	y := b.Top - float64(point.Y)/float64(size.Y)*(b.Height())
	return x, y
}
func (b Bounds) PositionOffset(pos fyne.Position, size fyne.Size) (float64, float64) {
	x := b.Left + float64(pos.X)/float64(size.Width)*(b.Width())
	y := b.Top - float64(pos.Y)/float64(size.Height)*(b.Height())
	return x, y
}

func (b Bounds) ZoomToBox(
	topLeft fyne.Position,
	botRight fyne.Position,
	size fyne.Size,
) Bounds {
	newBounds := Bounds{}
	newBounds.Left, newBounds.Top = b.PositionOffset(topLeft, size)
	newBounds.Right, newBounds.Bottom = b.PositionOffset(botRight, size)
	return newBounds
}

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
