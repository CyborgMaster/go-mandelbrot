package mandelbrot

import (
	"context"
	"fmt"
	"image"
	"math/cmplx"
)

const MAX_ITERATIONS = 1000

func PixelColor(point image.Point, size image.Point, bounds Bounds) uint8 {
	z := 0 + 0i
	c := bounds.PixelOffset(point, size).Complex()

	for i := range MAX_ITERATIONS {
		if cmplx.Abs(z) >= 2.0 {
			return uint8(i%255) + 1 // Use a simple color palette based on iteration count
		}
		z = cmplx.Pow(z, 2) + c
	}
	return 0
}

func DrawImage(
	ctx context.Context,
	img *image.Paletted,
	bounds Bounds,
	gorouintes int,
) (lineDone <-chan struct{}) {
	fmt.Println("starting render routines:", gorouintes)

	type Square struct {
		TopLeft image.Point
		Size    int
	}

	imageSize := img.Bounds().Size()

	minSide := min(imageSize.X, imageSize.Y)

	// We use a square size that is a power of 3, because then when rendering each
	// layer smaller, we can leave the middle of the 3x3 grid already drawn
	// because we sampled the center point when we did the larger square.
	startingSquareSize := 3
	for startingSquareSize <= minSide {
		startingSquareSize *= 3
	}
	startingSquareSize /= 27
	if startingSquareSize < 1 {
		startingSquareSize = 1
	}

	squares := make(chan Square)
	lines := make(chan struct{})

	for range gorouintes {
		go func() {
			for square := range squares {
				center := image.Point{
					X: square.TopLeft.X + square.Size/2,
					Y: square.TopLeft.Y + square.Size/2,
				}
				DrawSquare(img, square.TopLeft, square.Size,
					PixelColor(center, imageSize, bounds))
				if ctx.Err() != nil {
					return
				}
			}
		}()
	}

	go func() {
		for squareSize := startingSquareSize; squareSize >= 1; squareSize /= 3 {
			for y := 0; y < imageSize.Y; y += squareSize {
				for x := 0; x < imageSize.X; x += squareSize {
					row := y / squareSize
					col := x / squareSize
					if row%3 == 1 && col%3 == 1 && squareSize != startingSquareSize {
						// Skip the center of each 3x3 because we got the color right on the
						// larger square size.  We can't do this for the largest square
						// because there wasn't anything drawn before it.
						continue
					}

					squares <- Square{
						TopLeft: image.Point{X: x, Y: y},
						Size:    squareSize,
					}

					if ctx.Err() != nil {
						fmt.Println("redrawing raster cancelled")
						return
					}
				}
				lines <- struct{}{}
			}
		}
		close(squares)
		close(lines)
	}()

	return lines
}

func DrawSquare(
	img *image.Paletted,
	topLeft image.Point,
	size int,
	colorPalletIndex uint8,
) {
	for y := range size {
		for x := range size {
			img.SetColorIndex(topLeft.X+x, topLeft.Y+y, colorPalletIndex)
		}
	}
}
