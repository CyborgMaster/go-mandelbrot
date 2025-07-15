package main

import (
	"fmt"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/CyborgMaster/go-mandelbrot/mandelbrot"
)

func main() {
	fmt.Println("GOMAXPROCS:", runtime.GOMAXPROCS(-1))

	myApp := app.New()
	w := myApp.NewWindow("Raster")
	w.Resize(fyne.NewSize(800, 600))
	w.SetContent(mandelbrot.NewMandelbrotImage())
	w.SetPadded(false)
	w.ShowAndRun()
}

// TODO: Keep aspect ratio the same, expand the rendered image either vertically or horizontally to match the aspect ratio of the window
// TODO: don't flip things when we drag the rectangle the other way
// TODO: Zooming, by double click
// TODO: Progressive rendering (low res first)
// TODO: color cycling
// TODO: zoom out
// TODO: when dragging, show a rectangle that shows the area that will be zoomed in on
