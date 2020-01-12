package main

import (
	"fmt"
	"os"

	//"math"
	// "io"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/text"

	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

var LOCATION_file [100]byte
var LOCATION_matrix string

var matrix [][]uint8 = make([][]uint8, 4, 16)

var WIDTH float64 = 1100.0
var HEIGHT float64 = 700.0
var p_width float64 = WIDTH / 100
var p_height float64 = HEIGHT / 100

func run() {

	// basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	imd := imdraw.New(nil)

	cfg := pixelgl.WindowConfig{
		Title:  "Wesh Hennou",
		Bounds: pixel.R(0, 0, WIDTH, HEIGHT),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	imd.Color = colornames.Navy

	imd.Push(pixel.V(p_width*5, p_height*5), pixel.V(p_width*35, p_height*35)) // vertices for rect1 (bottom left)
	imd.Rectangle(0)

	imd.Push(pixel.V(p_width*95, p_height*5), pixel.V(p_width*65, p_height*35)) // bottom right
	imd.Rectangle(0)

	imd.Push(pixel.V(p_width*5, p_height*95), pixel.V(p_width*35, p_height*65)) // top left
	imd.Rectangle(0)

	imd.Push(pixel.V(p_width*95, p_height*95), pixel.V(p_width*65, p_height*65)) //(top right)
	imd.Rectangle(0)

	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)

	//Texte bouton encoder
	basicText := text.New(pixel.V(p_width*13, p_height*18), basicAtlas)
	basicText.Color = colornames.Limegreen
	fmt.Fprintln(basicText, "Encoder")

	//Texte bouton charger matrice
	basicText1 := text.New(pixel.V(p_width*13, p_height*75), basicAtlas)
	basicText1.Color = colornames.Limegreen
	fmt.Fprintln(basicText1, "Charger Matrice")

	//Texte bouton charger fichier
	basicText2 := text.New(pixel.V(p_width*75, p_height*75), basicAtlas)
	basicText2.Color = colornames.Limegreen
	fmt.Fprintln(basicText2, "Charger fichier")

	//Texte bouton decoder
	basicText3 := text.New(pixel.V(p_width*75, p_height*18), basicAtlas)
	basicText3.Color = colornames.Limegreen
	fmt.Fprintln(basicText3, "Decoder")

	for !win.Closed() {
		win.Clear(colornames.Aliceblue)

		imd.Draw(win)
		basicText.Draw(win, pixel.IM)
		basicText1.Draw(win, pixel.IM)
		basicText2.Draw(win, pixel.IM)
		basicText3.Draw(win, pixel.IM)

		if win.JustPressed(pixelgl.MouseButtonLeft) {
			button_handler(win)
		}

		win.Update()
	}
}
func main() {
	pixelgl.Run(run)
}

func button_handler(win *pixelgl.Window) {
	pos := win.MousePosition()

	if pos.Y > p_height*5 && pos.Y < p_height*35 { // clic bot

		if pos.X > p_width*5 && pos.X < p_width*35 { // bot left
			//insert matrix
			var index, endex uint8

			// cmd := exec.Command("find", "$HOME", "-name", "matrix.txt")

			// LOCATION_matrix, err := cmd.Output()
			// fmt.Printf("%s", LOCATION_matrix)

			data, err := os.Open("/home/thomas/go/Secu/src/matrix.txt")
			check(err)

			txt := make([]byte, 100)
			txt_len, err := data.Read(txt)
			txt_len++
			check(err)

			index, endex = seekKeyIndex(txt)
			insertMatrix(txt, index, endex)
			reorderMatrix()

		} else if pos.X > p_width*65 && pos.X < p_width*95 { // bot right
			//decrypt

		}

	} else {
		if pos.Y < p_height*95 && pos.Y > p_height*35 { // top

			if pos.X > p_width*5 && pos.X < p_width*35 { //top left
				//insert file

			} else if pos.X > p_width*65 && pos.X < p_width*95 { //top right
				//encrypt

			}
		}
	}

}

func encrypt()    {}
func decrypt()    {}
func selectFile() {}

func insertMatrix(file []byte, index uint8, endex uint8) {

	var i uint8 = index

	var x, y uint8 = 0, 0

	for i < endex {
		if file[i] == 32 { //spacebar

			x++
			y = 0 // newline
			i++   //we jump over the space char
		}
		matrix[x] = append(matrix[x], (file[i])-48)

		y++
		i++
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
func seekKeyIndex(data []byte) (uint8, uint8) {

	var index, endex uint8 = 0, 0
	var i uint8
	for i = 0; int(i) < len(data); i++ {
		if data[i] == 91 { //left bracket ascii code
			index = i + 1

		} else if data[i] == 93 {
			endex = i
		}
	}

	return index, endex

}

func reorderMatrix() {
	var i, j uint8

	var pos_one uint8
	var sum uint8 = 0

	var tmp_matrix [][]uint8 = make([][]uint8, 4, 16)

	for i = 0; int(i) < len(matrix[0]); i++ {
		for j = 0; j < 4; j++ {

			if matrix[j][i] == 1 {
				pos_one = j
				sum += matrix[j][i]

			}
		}
		if sum == 1 {
			tmp_matrix = matrix[pos_one]
		}

	}
}
