package main

import (
	"bufio"
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

// WorkingDirectory created on start , describes the folder where the program is
var WorkingDirectory string

// UsedFile describes path to matrix file
var UsedFile string

// IsFileSelected  => if statement on graphic app , make sure a file is selected before decrypting it
var IsFileSelected bool = false

// IsMatrixSelected => if statement on graphic app , make sure a matrix is selected before decrypting files
var IsMatrixSelected bool = false

var matrix [][]uint8 = make([][]uint8, 4, 16)

//WIDTH screen's width
var WIDTH float64 = 1100.0

//HEIGHT screen's height
var HEIGHT float64 = 700.0

var pWidth float64 = WIDTH / 100
var pHeight float64 = HEIGHT / 100

//MatrixIDOrder array representing which bits to extract from byte
var MatrixIDOrder []int = []int{5, 2, 3, 4}

//~~~~~~~~~~~~~~~~~PROGRAM DEBUT~~~~~~~~~~~~~~~~~~~~~~~~~

func run() {

	WorkingDirectory, err = os.Executable()
	check(err)
	WorkingDirectory = filepath.Dir(WorkingDirectory)

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

	imd.Push(pixel.V(pWidth*5, pHeight*5), pixel.V(pWidth*35, pHeight*35)) // vertices for rect1 (bottom left)
	imd.Rectangle(0)

	imd.Push(pixel.V(pWidth*95, pHeight*5), pixel.V(pWidth*65, pHeight*35)) // bottom right
	imd.Rectangle(0)

	imd.Push(pixel.V(pWidth*5, pHeight*95), pixel.V(pWidth*35, pHeight*65)) // top left
	imd.Rectangle(0)

	imd.Push(pixel.V(pWidth*95, pHeight*95), pixel.V(pWidth*65, pHeight*65)) //(top right)
	imd.Rectangle(0)
	//~~~~~~~~~~~~~~~~~~~~~~~ECRITURE DES TEXTES~~~~~~~~~~~~~~~~~~~~~~

	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)

	//Texte bouton charger matrice
	basicText := text.New(pixel.V(pWidth*13, pHeight*18), basicAtlas)
	basicText.Color = colornames.Limegreen
	fmt.Fprintln(basicText, "Charger Matrice")

	//Texte bouton Décoder
	basicText1 := text.New(pixel.V(pWidth*13, pHeight*75), basicAtlas)
	basicText1.Color = colornames.Limegreen
	fmt.Fprintln(basicText1, "Charger fichier")

	//Texte bouton encoder
	basicText2 := text.New(pixel.V(pWidth*75, pHeight*75), basicAtlas)
	basicText2.Color = colornames.Limegreen
	fmt.Fprintln(basicText2, "Encoder")

	//Texte bouton Charger fichier
	basicText3 := text.New(pixel.V(pWidth*75, pHeight*18), basicAtlas)
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
			buttonHandler(win)
		}

		win.Update()
	}
}
func main() {

	pixelgl.Run(run)
}

func buttonHandler(win *pixelgl.Window) {

	pos := win.MousePosition()

	if pos.Y > pHeight*5 && pos.Y < pHeight*35 { // clic bot

		if pos.X > pWidth*5 && pos.X < pWidth*35 { // bot left
			//insert matrix
			var index, endex uint8
			UsedFile = WorkingDirectory + "/matrix.txt"

			data, err := os.Open(UsedFile)
			check(err)
			fmt.Printf("Matrice insérée \n")
			txt := make([]byte, 100)
			_, err = data.Read(txt)
			check(err)

			index, endex = seekKeyIndex(txt)
			insertMatrix(txt, index, endex)
			// reorderMatrix()

			IsMatrixSelected = true
		} else if pos.X > pWidth*65 && pos.X < pWidth*95 { // bot right
			//decrypt
			if IsFileSelected && IsMatrixSelected {
				decryptFile()
				//TODO
			}
		}
	} else {
		if pos.Y < pHeight*95 && pos.Y > pHeight*35 { // top

			if pos.X > pWidth*5 && pos.X < pWidth*35 { //top left
				//insert file
				//by default file.txt
				selectFile()
			} else if pos.X > pWidth*65 && pos.X < pWidth*95 { //top right
				//encrypt
				if IsFileSelected && IsMatrixSelected {
					encryptFile()

				}
			}
		}
	}
}

func selectFile() {
	_, err := os.Open(WorkingDirectory + "/file.txt")
	check(err)

	IsFileSelected = true
	fmt.Printf("file selected \n")
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
	// fmt.Printf("%v", matrix)

	// var j uint8 = 0
	// var k uint8 = 1
	// var result uint8
	// for k < 4 {
	// 	for j < 8 {
	// 		result = matrix[0][j] ^ matrix[k][j]
	// 		matrix[k][j] = result
	// 		j++
	// 	}
	// 	k++
	// 	j = 0
	// }
	// fmt.Printf("\n%v", matrix)
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

func clearMatrix() {
	var i uint8 = 0

	for i < 4 {
		matrix[i] = matrix[i][:0]
		i++
	}
}

func encryptByte(theBytes []byte, length int) []byte {
	j, byteNum, loop := 0, 0, 0
	var i int

	var result string = ""
	var copyByte string = ""
	var bytes [2]string
	var sum int
	var strResult []string
	var resultInt uint64
	var sInt []byte

	for byteNum < length {
		tmpByte := strconv.FormatInt(int64(theBytes[byteNum]), 2)
		i = len(tmpByte)

		for i < 8 {
			copyByte += "0" // all bytes set to length 8
			i++
		}
		tmpByte = copyByte + tmpByte
		//text translated to binary OK
		bytes[0] = tmpByte[0:4]
		bytes[1] = tmpByte[4:8] // G4C matrix, every int coded on 8bites -> twice 4bits
		copyByte = ""
		loop = 0
		for loop < 2 { // len(bytes) = 2 at any run
			for i = 0; i < len(matrix[0]); i++ {

				sum = 0
				for j = 0; j < 4; j++ {
					if matrix[j][i] == 1 && bytes[loop][j] == 49 {
						sum++ // simulating XOR
					}
				}
				copyByte += strconv.Itoa(sum % 2) // end XOR
			}

			loop++
			copyByte += " "
		}
		// fmt.Printf("%d apres passage donne %s ", the_bytes, copy_byte)

		result += copyByte
		copyByte = ""

		byteNum++
	}
	strResult = strings.Split(result, " ")
	i = 0
	for i < len(strResult)-1 { // removing last element because ther's nothing in
		// fmt.Printf("%s\n", strResult[i])
		// result += longStringToIntString(strResult[i])
		resultInt, err = strconv.ParseUint(strResult[i], 2, 64)
		check(err)
		sInt = append(sInt, byte(resultInt))
		// fmt.Printf("%d\n", resultInt)

		i++
	}
	return sInt

}

func encryptFile() {
	var err = os.Remove(WorkingDirectory + "/file.txtc") // in case it already exists
	var writeTab []byte
	newfile, err := os.Create(WorkingDirectory + "/file.txtc")

	file, err := os.Open(WorkingDirectory + "/file.txt") // read from a file write into another
	check(err)

	/* 	write_tab = []byte{126}
	   	_, err = newfile.Write(write_tab) */ //permet de rentrer à la main 1 char , pour voir les effets
	currentByte := make([]byte, 1)
	for {
		//lecture d'un byte
		readByte, err := file.Read(currentByte)
		if err != nil {
			break
		}

		//cryptage d'un byte
		writeTab = encryptByte(currentByte, readByte)
		//ecriture d'un byte
		// fmt.Printf("%v", write_tab)
		_, err = newfile.Write(writeTab)
		check(err)

	}
	fmt.Println("file encrypted")
}

func decryptFile() {
	var err = os.Remove(WorkingDirectory + "/file.txtd") // in case it already exists
	// var write_byte []byte
	file, err := os.Open(WorkingDirectory + "/file.txtc")

	fi, err := file.Stat()
	fmt.Printf("file size : %d\n", fi.Size())

	check(err)

	readByte := bufio.NewReaderSize(file, int(fi.Size()))

	decryptByte(readByte, int(fi.Size()))

	// _, err = newfile.Write(write_byte)
	// check(err)

	fmt.Println("file decrypted")
	file.Close()

}

func decryptByte(reader *bufio.Reader, size int) {

	newfile, err := os.Create(WorkingDirectory + "/file.txtd")
	writeBytes := bufio.NewWriterSize(newfile, size/2)

	var i int = 0
	var concatResult string = ""
	var tmpByte string

	var leByte byte

	var writtenByte uint8
	for i < size {
		leByte, err = reader.ReadByte()

		tmpByte = fmt.Sprintf("%08b", leByte)

		tmpByte = string(tmpByte[4:5]) + string(tmpByte[1:2]) + string(tmpByte[2:3]) + string(tmpByte[3:4])

		concatResult += tmpByte
		if i%2 == 1 {
			writtenByte, _ = parseByte(concatResult)

			writeBytes.WriteByte(writtenByte)
			check(err)
			concatResult = ""

		}
		i++
	}
	// strResult := strings.Split(concat_result, " ")
	// i = 0
	// for i < len(strResult)-1 { // removing last element because ther's nothing in
	// 	resultInt, err := strconv.ParseUint(strResult[i], 2, 64)
	// 	check(err)
	// 	s_int = append(s_int, byte(resultInt))
	// 	i++
	// }
	// fmt.Println(s_int)
	// _, err = newfile.Write(s_int)
	// check(err)
	writeBytes.Flush()
	newfile.Close()

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

func parseByte(intStr string) (retV uint8, err error) {
	var value int64
	value, err = strconv.ParseInt(intStr, 2, 8)
	return byte(value), err
}

// func longStringToIntString(binary string) string {
// 	var qty_loops int = len(binary) / 8
// 	var i, j int = 0, 0

// 	var result string = ""
// 	var copy string = ""
// 	var tmp string
// 	for i < qty_loops {
// 		copy = ""
// 		j = 0
// 		for j < 8 {
// 			copy += string(binary[i*8+j])
// 			j++
// 		}
// 		tmp = parseBinToChar(copy)
// 		// fmt.Printf(tmp)
// 		result += tmp

// 		i++
// 	}
// 	// fmt.Println()
// 	return result
// }

// func reorderMatrix() {
// 	var i, j uint8

// 	var pos_one uint8
// 	var sum uint8 = 0

// 	var tmp_matrix []uint8 = make([]uint8, len(matrix[0]))
// 	var tmp_col []uint8 = make([]uint8, 4)

// 	for i = 0; int(i) < len(matrix[0]); i++ { // we know the index of identity matrix cols
// 		sum = 0
// 		for j = 0; j < 4; j++ {

// 			if matrix[j][i] == 1 {
// 				pos_one = j
// 				sum += matrix[j][i]
// 			}

// 		}
// 		if sum == 1 {
// 			tmp_matrix[pos_one] = i
// 		}
// 	}

// 	tmp_col = tmp_col[:0]
// 	for j = 0; int(j) < 4; j++ { // 4 : G4 ID matrix length
// 		for i = 0; i < 4; i++ {

// 			tmp_col = append(tmp_col, matrix[i][j]) //  saving current col in order to swap them in right order
// 		}

// 		for i = 0; i < 4; i++ {

// 			matrix[i][j] = matrix[i][tmp_matrix[j]]
// 		}

// 		for i = 0; i < 4; i++ {

// 			matrix[i][tmp_matrix[j]] = tmp_col[i]
// 		}

// 		tmp_col = tmp_col[:0]

// 	}

// 	// var matrix is now in order
// }
