package main

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"runtime"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"github.com/CyborgMaster/go-mandelbrot/mandelbrot"
	"github.com/lucasb-eyer/go-colorful"
)

type MandelbrotImage struct {
	*canvas.Raster

	image          *image.Paletted
	palette        color.Palette
	canvasSize     image.Point
	cancelDrawing  context.CancelFunc
	selectedBounds mandelbrot.Bounds
	drawnBounds    mandelbrot.Bounds

	dragging  bool
	dragStart fyne.Position
	dragEnd   fyne.Position

	colorAnimation *fyne.Animation
}

func NewMandelbrotImage() *MandelbrotImage {
	r := &MandelbrotImage{
		selectedBounds: mandelbrot.Bounds{
			Left:   -2.5,
			Right:  1.5,
			Top:    1.5,
			Bottom: -1.5,
		},
		palette: generatePallet(),
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

func generatePallet() color.Palette {
	palette := make(color.Palette, 256)
	palette[0] = color.Black
	for i := range 255 {
		palette[i+1] = colorful.Hsv(float64(360.0/255)*float64(i), 1, 0.5)

	}
	return palette
}

func rotatePalette(palette color.Palette) color.Palette {
	newPalette := make(color.Palette, 1, len(palette))
	newPalette[0] = palette[0] // black stays in the same spot
	newPalette = append(newPalette, palette[2:]...)
	newPalette = append(newPalette, palette[1])
	return newPalette
}

func (r *MandelbrotImage) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(r.Raster)
}

func (r *MandelbrotImage) Tapped(event *fyne.PointEvent) {
	if r.colorAnimation != nil {
		r.colorAnimation.Stop()
		r.colorAnimation = nil
		return
	}

	r.colorAnimation = fyne.NewAnimation(1*time.Hour, func(float32) {
		r.palette = rotatePalette(r.palette)
		r.image.Palette = r.palette
		r.Refresh()
	})
	r.colorAnimation.Start()
}

func (r *MandelbrotImage) DoubleTapped(event *fyne.PointEvent) {
	fmt.Println("DoubleTapped at", event.Position)
	r.selectedBounds = r.drawnBounds.ZoomToPoint(event.Position, r.Size(), 4)
	r.redraw()
}

func (r *MandelbrotImage) TappedSecondary(event *fyne.PointEvent) {
	fmt.Println("TappedSecondary at", event.Position)
	r.selectedBounds = r.drawnBounds.ZoomToPoint(event.Position, r.Size(), 1.0/4)
	r.redraw()
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
	r.selectedBounds = r.drawnBounds.ZoomToBox(r.dragStart, r.dragEnd, r.Size())
	r.redraw()
}

func (r *MandelbrotImage) redraw() {
	if r.cancelDrawing != nil {
		r.cancelDrawing()
		r.cancelDrawing = nil
	}
	ctx, cancel := context.WithCancel(context.Background())
	r.cancelDrawing = cancel
	size := image.Rect(0, 0, r.canvasSize.X, r.canvasSize.Y)
	r.image = image.NewPaletted(size, r.palette)
	r.drawnBounds = r.selectedBounds.MatchCanvasAspectRatio(r.canvasSize)
	go func(ctx context.Context, img *image.Paletted, bounds mandelbrot.Bounds) {
		linesDone := mandelbrot.DrawImage(ctx, img, bounds, runtime.GOMAXPROCS(-1))
		for range linesDone {
			fyne.Do(r.Refresh)
		}
		r.cancelDrawing = nil
	}(ctx, r.image, r.drawnBounds)
}
