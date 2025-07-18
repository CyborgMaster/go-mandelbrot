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

// TODO: display and allow adjusting the iteration limit
// TODO: when dragging, show a rectangle that shows the area that will be zoomed in on
// TODO: indicate render progress
// TODO: indicate zoom level using the exponent of the selected size
// TODO: have a history and you can go back to previous bounds
