package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"github.com/CyborgMaster/go-mandelbrot/mandelbrot"
)

// func main() {
// 	app := app.New()
// 	window := app.NewWindow("Hello")

// 	label := widget.NewLabel("Hello Fyne!")
// 	window.SetContent(container.NewVBox(
// 		label,
// 		widget.NewButton("Hi!", func() {
// 			label.SetText("Welcome :)")
// 		}),
// 	))

// 	window.ShowAndRun()
// }

// func main() {
// 	a := app.New()
// 	w := a.NewWindow("Hello")

// 	output := canvas.NewText(time.Now().Format(time.TimeOnly), color.NRGBA{G: 0xff, A: 0xff})
// 	output.TextStyle.Monospace = true
// 	output.TextSize = 32
// 	w.SetContent(output)

// 	go func() {
// 		ticker := time.NewTicker(time.Second)
// 		for range ticker.C {
// 			fyne.Do(func() {
// 				output.Text = time.Now().Format(time.TimeOnly)
// 				output.Refresh()
// 			})
// 		}
// 	}()
// 	w.ShowAndRun()
// }

// func main() {
// 	myApp := app.New()
// 	myWindow := myApp.NewWindow("Container")
// 	green := color.NRGBA{R: 0, G: 180, B: 0, A: 255}

// 	text1 := canvas.NewText("Hello", green)
// 	text2 := canvas.NewText("There", green)
// 	text2.Move(fyne.NewPos(20, 20))
// 	// content := container.NewWithoutLayout(text1, text2)
// 	content := container.New(layout.NewVBoxLayout(), text1, text2)

// 	clock := canvas.NewText(time.Now().Format(time.TimeOnly), color.NRGBA{G: 0xff, A: 0xff})
// 	clock.TextStyle.Monospace = true
// 	clock.TextSize = 32
// 	content.Add(clock)

// 	go func() {
// 		ticker := time.NewTicker(time.Second)
// 		for range ticker.C {
// 			fyne.Do(func() {
// 				clock.Text = time.Now().Format(time.TimeOnly)
// 				clock.Refresh()
// 			})
// 		}
// 	}()

// 	myWindow.SetContent(content)
// 	myWindow.ShowAndRun()
// }

func main() {
	myApp := app.New()
	w := myApp.NewWindow("Raster")

	raster := canvas.NewRasterWithPixels(mandelbrot.DrawPixel)
	// raster := canvas.NewRasterFromImage()
	w.SetContent(raster)
	w.Resize(fyne.NewSize(120, 100))
	w.ShowAndRun()
}
