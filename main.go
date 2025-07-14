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
	w.SetContent(mandelbrot.NewMandelbrotImage())
	w.Resize(fyne.NewSize(800, 600))
	w.ShowAndRun()
}

// TODO: Zooming
// TODO: Keep aspect ratio the same
// TODO: Colors
// TODO: Progressive rendering (low res first)
// TODO: color cycling
