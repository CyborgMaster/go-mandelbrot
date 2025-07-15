package mandelbrot

import (
	"image"

	"fyne.io/fyne/v2"
)

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
