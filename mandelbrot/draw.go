package mandelbrot

import (
	"image/color"
	"math/cmplx"
)

const MAX_ITERATIONS = 1000

var left = -2.0
var right = 1.0
var top = 1.0
var bottom = -1.0

func DrawPixel(x, y, w, h int) color.Color {
	z := 0 + 0i
	real_part := left + float64(x)/float64(w)*(right-left)
	img_part := top - float64(y)/float64(h)*(top-bottom)
	c := complex(real_part, img_part)

	for i := 0; i < MAX_ITERATIONS; i++ {
		if cmplx.Abs(z) >= 2.0 {
			return color.White
		}
		z = cmplx.Pow(z, 2) + c
	}
	return color.Black
}
