package main

import(
	 "fmt"
	 //"math"


	 "github.com/faiface/pixel"
	 "github.com/faiface/pixel/pixelgl"
	 "github.com/faiface/pixel/text"
	 "github.com/faiface/pixel/imdraw"

	 "golang.org/x/image/colornames"

)

func run() {
	WIDTH 		:= 500.0
	HEIGHT 		:= 250.0
	p_width 	:= WIDTH/100
	p_height	:= HEIGHT/100

	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	imd := imdraw.New(nil)


	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, WIDTH, HEIGHT),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}


	imd.Color = colornames.Navy
	
	imd.Push(pixel.V(p_width * 5 , p_height * 5 ) ,  pixel.V( p_width * 35  , p_height * 35) ) // vertices for rect1 (bottom left)
	// buttons[0] = *pixelgl.NewCanvas(pixel.R(p_width * 5 , p_height * 5  , p_width * 15  , p_height * 13 ))
	imd.Rectangle(0)

	
	imd.Push(pixel.V(p_width * 95 , p_height * 5 ) ,  pixel.V( p_width * 65  , p_height * 35) ) // bottom right
	imd.Rectangle(0)

	imd.Push(pixel.V(p_width * 5 , p_height * 95 ) ,  pixel.V( p_width * 35  , p_height * 65) ) // top left
	imd.Rectangle(0)

	imd.Push(pixel.V(p_width * 95 , p_height * 95) , pixel.V( p_width * 65 , p_height * 65) ) //(top right)
	imd.Rectangle(0)
	// imd.Color = colornames.Limegreen
	// imd.Color = colornames.Navy
	// imd.Color = colornames.Red

	for !win.Closed() {
		win.Clear(colornames.Aliceblue)
		imd.Draw(win)

		if win.JustPressed(pixelgl.MouseButtonLeft) {
			button_handler(win)
		}

		win.Update()
	}
}
func main() {
	pixelgl.Run(run)
}


func button_handler (win *pixelgl.Window) {
	pos := win.MousePosition()

	fmt.Printf("%f" , pos.X)
}