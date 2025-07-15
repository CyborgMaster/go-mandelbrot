package main

import (
	"fmt"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	fmt.Println("GOMAXPROCS:", runtime.GOMAXPROCS(-1))

	myApp := app.New()
	w := myApp.NewWindow("Raster")
	w.Resize(fyne.NewSize(800, 600))
	w.SetContent(NewMandelbrotImage())
	w.SetPadded(false)
	w.ShowAndRun()
}

// TODO: Zooming, by double click
// TODO: Progressive rendering (low res first)
// TODO: color cycling, use a better color palette based on HCL
// TODO: zoom out
// TODO: when dragging, show a rectangle that shows the area that will be zoomed in on
