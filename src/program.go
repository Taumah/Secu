package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/text"

	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

//~~~~~~~~~~~~~~GLOBAL VARIABLES~~~~~~~~~~~~~
var err error

var WORKING_DIRECTORY string
var USED_FILE string

var SELECT_file bool = false
var SELECT_matrix bool = false

var matrix [][]uint8 = make([][]uint8, 4, 16)

var WIDTH float64 = 1100.0
var HEIGHT float64 = 700.0
var p_width float64 = WIDTH / 100
var p_height float64 = HEIGHT / 100

var my_var string = ""

//~~~~~~~~~~~~~~~~~PROGRAM DEBUT~~~~~~~~~~~~~~~~~~~~~~~~~

func run() {

	WORKING_DIRECTORY, err = os.Executable()
	check(err)
	WORKING_DIRECTORY = filepath.Dir(WORKING_DIRECTORY)

	imd := imdraw.New(nil)

	cfg := pixelgl.WindowConfig{
		Title:  "Chiffrement",
		Bounds: pixel.R(0, 0, WIDTH, HEIGHT),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	//~~~~~~~~~~~~~~~~~~~~~~~DESSIN DES BOUTONS~~~~~~~~~~~~~~~~~~~~~~~
	imd.Color = colornames.Navy

	imd.Push(pixel.V(p_width*5, p_height*5), pixel.V(p_width*35, p_height*35)) // vertices for rect1 (bottom left)
	imd.Rectangle(0)

	imd.Push(pixel.V(p_width*95, p_height*5), pixel.V(p_width*65, p_height*35)) // bottom right
	imd.Rectangle(0)

	imd.Push(pixel.V(p_width*5, p_height*95), pixel.V(p_width*35, p_height*65)) // top left
	imd.Rectangle(0)

	imd.Push(pixel.V(p_width*95, p_height*95), pixel.V(p_width*65, p_height*65)) //(top right)
	imd.Rectangle(0)
	//~~~~~~~~~~~~~~~~~~~~~~~ECRITURE DES TEXTES~~~~~~~~~~~~~~~~~~~~~~

	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)

	//Texte bouton charger matrice
	basicText := text.New(pixel.V(p_width*13, p_height*18), basicAtlas)
	basicText.Color = colornames.Limegreen
	fmt.Fprintln(basicText, "Charger Matrice")

	//Texte bouton Décoder
	basicText1 := text.New(pixel.V(p_width*13, p_height*75), basicAtlas)
	basicText1.Color = colornames.Limegreen
	fmt.Fprintln(basicText1, "Charger fichier")

	//Texte bouton encoder
	basicText2 := text.New(pixel.V(p_width*75, p_height*75), basicAtlas)
	basicText2.Color = colornames.Limegreen
	fmt.Fprintln(basicText2, "Encoder")

	//Texte bouton Charger fichier
	basicText3 := text.New(pixel.V(p_width*75, p_height*18), basicAtlas)
	basicText3.Color = colornames.Limegreen
	fmt.Fprintln(basicText3, "Decoder")

	//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

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
			USED_FILE = WORKING_DIRECTORY + "/matrix.txt"

			data, err := os.Open(USED_FILE)
			check(err)

			txt := make([]byte, 100)
			txt_len, err := data.Read(txt)
			txt_len++
			check(err)

			index, endex = seekKeyIndex(txt)
			insertMatrix(txt, index, endex)
			// reorderMatrix()

			SELECT_matrix = true
		} else if pos.X > p_width*65 && pos.X < p_width*95 { // bot right
			//decrypt
			if SELECT_file && SELECT_matrix {
				decrypt_file()
				//TODO
			}
		}
	} else {
		if pos.Y < p_height*95 && pos.Y > p_height*35 { // top

			if pos.X > p_width*5 && pos.X < p_width*35 { //top left
				//insert file
				//by default file.txt
				selectFile()
			} else if pos.X > p_width*65 && pos.X < p_width*95 { //top right
				//encrypt
				if SELECT_file && SELECT_matrix {
					encrypt_file()

				}
			}
		}
	}
}

func selectFile() {
	data, err := os.Open(WORKING_DIRECTORY + "/file.txt")
	check(err)

	txt := make([]byte, 100)
	_, err = data.Read(txt)
	SELECT_file = true
}

func insertMatrix(file []byte, index uint8, endex uint8) {
	clearMatrix()
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

	var tmp_matrix []uint8 = make([]uint8, len(matrix[0]))
	var tmp_col []uint8 = make([]uint8, 4)

	for i = 0; int(i) < len(matrix[0]); i++ { // we know the index of identity matrix cols
		sum = 0
		for j = 0; j < 4; j++ {

			if matrix[j][i] == 1 {
				pos_one = j
				sum += matrix[j][i]
			}

		}
		if sum == 1 {
			tmp_matrix[pos_one] = i
		}
	}

	tmp_col = tmp_col[:0]
	for j = 0; int(j) < 4; j++ { // 4 : G4 ID matrix length
		for i = 0; i < 4; i++ {

			tmp_col = append(tmp_col, matrix[i][j]) //  saving current col in order to swap them in right order
		}

		for i = 0; i < 4; i++ {

			matrix[i][j] = matrix[i][tmp_matrix[j]]
		}

		for i = 0; i < 4; i++ {

			matrix[i][tmp_matrix[j]] = tmp_col[i]
		}

		tmp_col = tmp_col[:0]

	}

	// var matrix is now in order
}

func clearMatrix() {
	var i uint8 = 0

	for i < 4 {
		matrix[i] = matrix[i][:0]
		i++
	}
}

func encrypt_byte(the_bytes []byte, length int) []byte {
	j, byte_num, loop := 0, 0, 0
	var i int

	var result string = ""
	var copy_byte string = ""
	var bytes [2]string
	var sum int
	var str_result []string
	var result_int uint64
	var s_int []byte

	for byte_num < length {
		tmp_byte := strconv.FormatInt(int64(the_bytes[byte_num]), 2)
		i = len(tmp_byte)

		for i < 8 {
			copy_byte += "0" // all bytes set to length 8
			i++
		}
		tmp_byte = copy_byte + tmp_byte
		//text translated to binary OK
		bytes[0] = tmp_byte[0:4]
		bytes[1] = tmp_byte[4:8] // G4C matrix, every int coded on 8bites -> twice 4bits
		copy_byte = ""
		loop = 0
		for loop < 2 { // len(bytes) = 2 at any run
			for i = 0; i < len(matrix[0]); i++ {

				sum = 0
				for j = 0; j < 4; j++ {
					if matrix[j][i] == 1 && bytes[loop][j] == 49 {
						sum++ // simulating XOR
					}
				}
				copy_byte += strconv.Itoa(sum % 2) // end XOR
			}

			loop++
			copy_byte += " "
		}
		// fmt.Printf("%d apres passage donne %s ", the_bytes, copy_byte)

		result += copy_byte
		copy_byte = ""

		byte_num++
	}
	my_var += result
	str_result = strings.Split(result, " ")
	i = 0
	for i < len(str_result)-1 { // removing last element because ther's nothing in
		fmt.Printf("%s\n", str_result[i])
		// result += longStringToIntString(str_result[i])
		result_int, err = strconv.ParseUint(str_result[i], 2, 64)
		check(err)
		s_int = append(s_int, byte(result_int))
		fmt.Printf("%d\n", result_int)

		i++
	}
	return s_int

}

func encrypt_file() {
	var err = os.Remove(WORKING_DIRECTORY + "/file.txtc") // in case it already exists
	var write_tab []byte
	newfile, err := os.Create(WORKING_DIRECTORY + "/file.txtc")

	file, err := os.Open(WORKING_DIRECTORY + "/file.txt") // read from a file write into another
	check(err)

	current_byte := make([]byte, 1)
	for {
		//lecture d'un byte
		read_byte, err := file.Read(current_byte)
		if err != nil {
			break
		}

		//cryptage d'un byte
		write_tab = encrypt_byte(current_byte, read_byte)
		//ecriture d'un byte
		fmt.Printf("%v", write_tab)
		_, err = newfile.Write(write_tab)
		check(err)

	}
	// fmt.Println(my_var)
	fmt.Println("file encrypted")
}

func decrypt_file() {
	var err = os.Remove(WORKING_DIRECTORY + "/file.txtd") // in case it already exists
	newfile, err := os.Create(WORKING_DIRECTORY + "/file.txtd")

	var write_byte string
	file, err := os.Open(WORKING_DIRECTORY + "/file.txtc")
	check(err)

	read_byte := make([]byte, len(matrix[0])*2)

	for {
		//lecture d'un byte
		for i := range read_byte {
			read_byte[i] = 0
		}

		decrypt_bytes, err := file.Read(read_byte)
		if err != nil {
			break //reading until we can't anymore (EOF)
		}
		write_byte = decrypt_byte(read_byte, decrypt_bytes)

		_, err = newfile.WriteString(write_byte)
		check(err)

	}
	fmt.Println("file decrypted")
}

func decrypt_byte(the_bytes []byte, length int) string {
	var i int = 0
	var concat_bins string = ""
	var concat_result string = ""
	for i < length {
		concat_bins += fmt.Sprintf("%08b", the_bytes[i])
		i++
	}

	// fmt.Println(concat_bins)

	i = 0
	for i < length {

		for j := 0; j < 4; j++ {
			concat_result += fmt.Sprintf("%c", concat_bins[i*8+j])
		}
		i++
	}
	result := longStringToIntString(concat_result)
	// fmt.Println(result)

	return result
}

func parseBinToChar(s string) string { //smartest result from Stack
	ui, err := strconv.ParseUint(s, 2, 64)
	check(err)

	return fmt.Sprintf("%c", ui)
}
func parseIntToBin(Int int64) string { //smartest result from Stack
	var format string = "%08b"
	ui := strconv.FormatInt(Int, 2)
	check(err)
	//ex : %016b indicates base 2, zero padded, with 16 characters
	return fmt.Sprintf(format, ui)
}

func longStringToIntString(binary string) string {
	var qty_loops int = len(binary) / 8
	var i, j int = 0, 0

	var result string = ""
	var copy string = ""
	var tmp string
	for i < qty_loops {
		copy = ""
		j = 0
		for j < 8 {
			copy += string(binary[i*8+j])
			j++
		}
		tmp = parseBinToChar(copy)
		// fmt.Printf(tmp)
		result += tmp

		i++
	}
	// fmt.Println()
	return result
}
