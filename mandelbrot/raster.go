package mandelbrot

import (
	"context"
	"fmt"
	"image"
	"image/draw"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

type MandelbrotImage struct {
	*canvas.Raster

	img           draw.Image
	cancelDrawing context.CancelFunc
}

func NewMandelbrotImage() *MandelbrotImage {
	r := &MandelbrotImage{}
	r.Raster = canvas.NewRaster(func(w, h int) image.Image {
		if r.img == nil ||
			r.img.Bounds().Size().X != w ||
			r.img.Bounds().Size().Y != h {

			r.redraw(w, h)
		}

		return r.img
	})

	return r
}

func (r *MandelbrotImage) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(r.Raster)
}

func (r *MandelbrotImage) Tapped(event *fyne.PointEvent) {
	fmt.Println("Tapped at", event.Position)
}

func (r *MandelbrotImage) redraw(w, h int) {
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
		fmt.Println("redrawing", w, h)
		for range DrawImage(ctx, img, runtime.GOMAXPROCS(-1)) {
			fyne.Do(r.Refresh)
		}
		r.cancelDrawing = nil
	}(ctx, img)
}
