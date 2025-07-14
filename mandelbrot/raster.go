package mandelbrot

import (
	"context"
	"fmt"
	"image"
	"image/draw"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

type Raster struct {
	r *canvas.Raster

	img           draw.Image
	cancelDrawing context.CancelFunc
}

func (r *Raster) redraw(w, h int) {
	if r.cancelDrawing != nil {
		r.cancelDrawing()
		r.cancelDrawing = nil
	}
	rect := image.Rect(0, 0, w, h)
	img := image.NewRGBA(rect)
	r.img = img
	ctx, cancel := context.WithCancel(context.Background())
	r.cancelDrawing = cancel

	go func(ctx context.Context, img draw.Image) {
		fmt.Println("redrawing raster", w, h)
		h := img.Bounds().Size().Y
		w := img.Bounds().Size().X
		for y := range h {
			for x := range w {
				img.Set(x, y, Pixel(x, y, w, h))
				if ctx.Err() != nil {
					fmt.Println("redrawing raster cancelled")
					return
				}
			}
		}
		r.cancelDrawing = nil
		fyne.Do(r.r.Refresh)
	}(ctx, img)
}

func NewRaster() *canvas.Raster {
	r := &Raster{}
	r.r = canvas.NewRaster(func(w, h int) image.Image {
		if r.img == nil ||
			r.img.Bounds().Size().X != w ||
			r.img.Bounds().Size().Y != h {

			r.redraw(w, h)
		}

		return r.img
	})

	return r.r
}
