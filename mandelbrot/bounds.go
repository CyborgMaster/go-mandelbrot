package mandelbrot

import (
	"image"

	"fyne.io/fyne/v2"
)

type Point struct {
	X float64
	Y float64
}

func (p Point) Complex() complex128 {
	return complex(p.X, p.Y)
}

// These match the X/Y plane of the mandelbrot, so up is positive Y, right is
// positive X.
type Bounds struct {
	Left   float64
	Right  float64
	Top    float64
	Bottom float64
}

func NewBoundsFromPoints(p1, p2 Point) Bounds {
	return Bounds{
		Left:   min(p1.X, p2.X),
		Right:  max(p1.X, p2.X),
		Top:    max(p1.Y, p2.Y),
		Bottom: min(p1.Y, p2.Y),
	}
}

func (b Bounds) Width() float64 {
	return b.Right - b.Left
}

func (b Bounds) Height() float64 {
	return b.Top - b.Bottom
}

func (b Bounds) Center() Point {
	return Point{
		X: (b.Left + b.Right) / 2,
		Y: (b.Top + b.Bottom) / 2,
	}
}

func (b Bounds) PixelOffset(point image.Point, size image.Point) Point {
	return Point{
		X: b.Left + float64(point.X)/float64(size.X)*(b.Width()),
		Y: b.Top - float64(point.Y)/float64(size.Y)*(b.Height()),
	}
}

func (b Bounds) PositionOffset(pos fyne.Position, size fyne.Size) Point {
	return Point{
		X: b.Left + float64(pos.X)/float64(size.Width)*(b.Width()),
		Y: b.Top - float64(pos.Y)/float64(size.Height)*(b.Height()),
	}
}

func (b Bounds) ZoomToBox(
	topLeft fyne.Position,
	botRight fyne.Position,
	size fyne.Size,
) Bounds {
	return NewBoundsFromPoints(
		b.PositionOffset(topLeft, size),
		b.PositionOffset(botRight, size),
	)
}

func (b Bounds) MatchCanvasAspectRatio(size image.Point) Bounds {
	// Calculate the aspect ratio of the current bounds
	currentAspect := b.Width() / b.Height()
	// Calculate the aspect ratio of the image size
	canvasAspect := float64(size.X) / float64(size.Y)

	newBounds := b

	if currentAspect > canvasAspect {
		// Current bounds are wider than the canvas, adjust height
		newHeight := b.Width() / canvasAspect
		centerY := b.Center().Y
		newBounds.Top = centerY + newHeight/2
		newBounds.Bottom = centerY - newHeight/2
	} else {
		// Current bounds are taller than the canvas, adjust width
		newWidth := b.Height() * canvasAspect
		centerX := b.Center().X
		newBounds.Left = centerX - newWidth/2
		newBounds.Right = centerX + newWidth/2
	}

	return newBounds
}
