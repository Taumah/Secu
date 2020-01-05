package main

import(
	 "fmt"
	 //"math"


	 "github.com/faiface/pixel"
	 "github.com/faiface/pixel/pixelgl"

	 "golang.org/x/image/colornames"
	 "github.com/faiface/pixel/imdraw"

)
func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 400, 120),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	imd := imdraw.New(nil)

	imd.Color = colornames.Navy
	imd.Push(pixel.V(20,20), pixel.V( 30 , 30) )
	myrect := imd.Rectangle(0)
	
	// imd.Color = colornames.Limegreen
	// imd.Color = colornames.Navy
	// imd.Color = colornames.Red

	for !win.Closed() {
		win.Clear(colornames.Aliceblue)
		imd.Draw(win)

		
		win.Update()
	}
}
func main() {
	fmt.Println("Hello, world.")

	pixelgl.Run(run)
}
