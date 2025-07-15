package main

import (
	"context"
	"fmt"
	"image"
	"image/draw"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"github.com/CyborgMaster/go-mandelbrot/mandelbrot"
)

type MandelbrotImage struct {
	*canvas.Raster

	image            draw.Image
	canvasSize       image.Point
	cancelDrawing    context.CancelFunc
	mandelbrotBounds mandelbrot.Bounds

	dragging  bool
	dragStart fyne.Position
	dragEnd   fyne.Position
}

func NewMandelbrotImage() *MandelbrotImage {
	r := &MandelbrotImage{
		mandelbrotBounds: mandelbrot.Bounds{
			Left:   -2.5,
			Right:  1.5,
			Top:    1.5,
			Bottom: -1.5,
		},
	}

	r.Raster = canvas.NewRaster(func(w, h int) image.Image {
		size := image.Point{X: w, Y: h}
		if r.image == nil || r.canvasSize != size {
			fmt.Println("resizing", w, h)
			r.canvasSize = size
			r.redraw()
		}

		return r.image
	})

	return r
}

func (r *MandelbrotImage) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(r.Raster)
}

func (r *MandelbrotImage) Tapped(event *fyne.PointEvent) {
	fmt.Println("Tapped at", event.Position)
}

func (r *MandelbrotImage) Dragged(event *fyne.DragEvent) {
	if !r.dragging {
		r.dragging = true
		r.dragStart = event.Position.Subtract(event.Dragged)
	}
	r.dragEnd = event.Position
}

func (r *MandelbrotImage) DragEnd() {
	fmt.Println("Dragged", r.dragStart, "to", r.dragEnd)
	r.dragging = false
	r.mandelbrotBounds = r.mandelbrotBounds.ZoomToBox(r.dragStart, r.dragEnd, r.Size())
	r.redraw()
}

func (r *MandelbrotImage) redraw() {
	if r.cancelDrawing != nil {
		r.cancelDrawing()
		r.cancelDrawing = nil
	}
	ctx, cancel := context.WithCancel(context.Background())
	r.cancelDrawing = cancel
	r.image = image.NewRGBA(image.Rect(0, 0, r.canvasSize.X, r.canvasSize.Y))
	go func(ctx context.Context, img draw.Image) {
		linesDone := mandelbrot.DrawImage(
			ctx,
			img,
			r.mandelbrotBounds,
			runtime.GOMAXPROCS(-1),
		)
		for range linesDone {
			fyne.Do(r.Refresh)
		}
		r.cancelDrawing = nil
	}(ctx, r.image)
}
