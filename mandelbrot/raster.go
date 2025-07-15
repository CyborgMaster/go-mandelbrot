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

	img              draw.Image
	cancelDrawing    context.CancelFunc
	mandelbrotBounds Bounds

	dragging  bool
	dragStart fyne.Position
	dragEnd   fyne.Position
}

func NewMandelbrotImage() *MandelbrotImage {
	r := &MandelbrotImage{
		mandelbrotBounds: Bounds{
			Left:   -2.5,
			Right:  1.5,
			Top:    1.5,
			Bottom: -1.5,
		},
	}

	r.Raster = canvas.NewRaster(func(w, h int) image.Image {
		if r.img == nil ||
			r.img.Bounds().Size().X != w ||
			r.img.Bounds().Size().Y != h {

			fmt.Println("resizing", w, h)
			rect := image.Rect(0, 0, w, h)
			img := image.NewRGBA(rect)
			r.img = img
			r.redraw()
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

	go func(ctx context.Context, img draw.Image) {
		for range DrawImage(ctx, img, r.mandelbrotBounds, runtime.GOMAXPROCS(-1)) {
			fyne.Do(r.Refresh)
		}
		r.cancelDrawing = nil
	}(ctx, r.img)
}
