package main

import (
	"bufio"
	"fmt"
	"math"
	"math/bits"
	"os"
	"strconv"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/sqweek/dialog"

	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

//~~~~~~~~~~~~~~GLOBAL VARIABLES~~~~~~~~~~~~~
var err error

// MatrixPath describes path to matrix file
var MatrixPath string

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
var MatrixIDOrder []uint8 = []uint8{4, 1, 2, 3}

//MatrixValuesAsLine each element is a matrix line represented like a byte (easier to manipulate)
var MatrixValuesAsLine []uint8 = []uint8{0, 0, 0, 0}

var arrayMatrixCondition []float64 = []float64{0, 0, 0, 0}

//~~~~~~~~~~~~~~~~~PROGRAM DEBUT~~~~~~~~~~~~~~~~~~~~~~~~~

func run() {

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

	imd.Push(pixel.V(pWidth*5, pHeight*35), pixel.V(pWidth*40, pHeight*75)) // vertices for rect1 (bottom left)
	imd.Rectangle(0)

	//~~~~~~~~~~~~~~~~~~~~~~~ECRITURE DES TEXTES~~~~~~~~~~~~~~~~~~~~~~

	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)

	//Texte bouton charger matrice
	basicText := text.New(pixel.V(pWidth*8, pHeight*55), basicAtlas)
	basicText.Color = colornames.Limegreen
	fmt.Fprintln(basicText, "1 - Charger Matrice")

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
		basicText.Draw(win, pixel.IM.Scaled(basicText.Orig, 2))

		if win.JustPressed(pixelgl.MouseButtonLeft) {
			buttonHandler(win)
		}

		if IsMatrixSelected == true {
			imd.Push(pixel.V(pWidth*95, pHeight*5), pixel.V(pWidth*65, pHeight*35)) // bottom right
			imd.Rectangle(0)

			imd.Push(pixel.V(pWidth*95, pHeight*95), pixel.V(pWidth*65, pHeight*65)) //(top right)
			imd.Rectangle(0)

			basicText2.Draw(win, pixel.IM.Scaled(basicText2.Orig, 2))
			basicText3.Draw(win, pixel.IM.Scaled(basicText3.Orig, 2))
		}

		win.Update()
	}
}
func main() {
	pixelgl.Run(run)
}

func buttonHandler(win *pixelgl.Window) {

	pos := win.MousePosition()

	if pos.Y > pHeight*35 && pos.Y < pHeight*75 { // clic bot

		if pos.X > pWidth*5 && pos.X < pWidth*40 { // bot left
			//insert matrix
			var index, endex uint8
			MatrixPath, err := dialog.File().Filter("Fichier Texte", "txt").Load()
			check(err)
			data, err := os.Open(MatrixPath)
			check(err)
			txt := make([]byte, 100)
			_, err = data.Read(txt)
			check(err)

			index, endex = seekKeyIndex(txt)
			if endex == 0 {
				dialog.Message("%s", "Mauvais format de matrice, veuillez ressayer").Title("Aie aie aie").Info()
				IsMatrixSelected = false
			} else {
				insertMatrix(txt, index, endex)
				// reorderMatrix()
				fillMatrixIDOrder()
				fillMatrixValueLine(txt, index, endex)
				dialog.Message("%s", "Matrix Loaded").Title("Success !!").Info()
				IsMatrixSelected = true
			}
		}
	} else if pos.Y > pHeight*5 && pos.Y < pHeight*35 {
		if pos.X > pWidth*65 && pos.X < pWidth*95 { // bot right
			//decrypt
			if IsMatrixSelected {
				decryptFile()
			}
		}
	} else {
		if pos.Y < pHeight*95 && pos.Y > pHeight*35 { // top
			if pos.X > pWidth*65 && pos.X < pWidth*95 { //top right
				//encrypt
				if IsMatrixSelected {
					encryptFile()
				}
			}
		}
	}
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
	if endex-index != 35 {
		return index, 0
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

func encryptByte(bytes *bufio.Reader, length int, pathToCreate string) {
	var i, j int = 0, 0
	var byteNum int = 0
	// var futureByte [2]byte

	var cutByte [2]byte
	var condition uint8 = 0 //use to know if any bit is on/off
	newfile, _ := os.Create(pathToCreate)
	writeBytes := bufio.NewWriter(newfile)

	for byteNum < length {
		// futureByte = 0
		theReadByte, err := bytes.ReadByte()
		check(err)
		cutByte[0] = 0
		cutByte[1] = 0
		i = 0

		condition = 0b10000000
		for i < 2 {
			j = 0
			for j < 4 {
				// fmt.Printf("%d & %d : %b \n", theReadByte, condition, (theReadByte&condition) == condition)
				// (theReadByte & condition) == condition

				if (theReadByte & condition) == condition {
					cutByte[i] ^= MatrixValuesAsLine[j]
					// fmt.Printf("%d cutByte apres %08b (%d)\n", byteNum, cutByte[i], cutByte[i])
					// fmt.Printf("%v \n", cutByte)

				}
				condition >>= 1
				j++
			}
			// fmt.Println()

			// fmt.Println("écriture:  ", cutByte[i])
			writeBytes.WriteByte(cutByte[i])
			i++

		}

		byteNum++
	}
	writeBytes.Flush()

}

func encryptFile() {

	// var writeTab []byte
	filename, err := dialog.File().Title("Chose a file to encrypt").Load()
	name := getFileName(filename)
	exten := getFileExt(filename)
	pathCryptedFile, err := dialog.Directory().Title("Chose a directory to save your encrypted file").Browse()
	// err = os.Remove(pathCryptedFile + "/encryptedFile.txtc") // in case it already exists
	newfile, err := os.Create(pathCryptedFile + "/" + name + "." + exten + "c")
	file, err := os.Open(filename) // read from a file write into another
	check(err)

	fi, err := file.Stat()
	check(err)

	bufferReader := bufio.NewReaderSize(file, int(fi.Size()))

	newfile.Close()

	path := pathCryptedFile + "/" + name + "." + exten + "c"
	start := time.Now()
	encryptByte(bufferReader, int(fi.Size()), path)

	t := time.Now()
	elapsed := t.Sub(start)

	dialog.Message("%s", "File encrypted").Title("Success !!").Info()
	fmt.Printf("%.0f", elapsed.Seconds())
}

func decryptFile() {

	filename, err := dialog.File().Title("Chose a file to decrypt").Load()

	name := getFileName(filename)
	exten := getFileExt(filename)
	exten = string(exten[0 : len(exten)-1])

	pathDecryptedFile, err := dialog.Directory().Title("Chose a directory to save your decrypted file").Browse()
	check(err)
	file, err := os.Open(filename)
	fi, err := file.Stat()
	// fmt.Printf("file size : %d\n", fi.Size())
	check(err)

	readByte := bufio.NewReaderSize(file, int(fi.Size()))
	start := time.Now()
	decryptByte(readByte, fi.Size(), pathDecryptedFile, name, exten)
	t := time.Now()
	elapsed := t.Sub(start)
	dialog.Message("%s", "File decrypted").Title("Success !!").Info()
	fmt.Printf("%.0f", elapsed.Seconds())
	file.Close()

}

func decryptByte(bytes *bufio.Reader, size int64, path string, name string, exten string) {

	newfile, err := os.Create(path + "/" + name + "." + exten + "d")
	writeBytes := bufio.NewWriter(newfile)

	var k int = 0
	var i, j int64 = 0, 0
	var bitPos float64
	//var tmpByte string
	var leByteDecomp uint8
	var leByteRead uint8

	// var writtenByte byte
	// fmt.Println(arrayMatrixCondition)
	for i < size-1 {

		leByteDecomp, j, bitPos = 0, 0, 0
		for j < 2 {
			k = 0
			leByteRead, _ = bytes.ReadByte()
			for k < 4 { //id matrix length
				condition := uint8(arrayMatrixCondition[k])
				if leByteRead&condition == condition {
					leByteDecomp += uint8(math.Pow(2, bitPos))
				}
				k++
				bitPos++
			}

			j++
		}
		leByteDecomp = bits.Reverse8(leByteDecomp)

		writeBytes.WriteByte(leByteDecomp)
		check(err)
		i += 2
	}

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
	var value uint64
	value, _ = strconv.ParseUint(intStr, 2, 8)
	return uint8(value), err
}

func fillMatrixIDOrder() {

	var i, j uint8
	var sum, posOne uint8
	for i = 0; int(i) < len(matrix[0]); i++ { // we know the index of identity matrix cols
		sum = 0
		for j = 0; j < 4; j++ {
			if matrix[j][i] == 1 {
				posOne = j
				sum += matrix[j][i]
			}
		}
		if sum == 1 {
			MatrixIDOrder[posOne] = i
			arrayMatrixCondition[posOne] = math.Pow(2, float64(8-i-1))
			//fmt.Println(arrayMatrixCondition[posOne])
			// fmt.Println("result : ", MatrixIDOrder, "new ", posOne)
		}
	}
}

func getFileName(filename string) string {
	var i int
	var lastSlash int
	var beforeName int

	for i = 0; i < len(filename); i++ {
		if filename[i] == 92 {
			lastSlash = i
		}
	}

	for i = len(filename) - 1; i > lastSlash; i-- {
		if filename[i] == 46 {
			beforeName = i
		}
	}
	name := string(filename[lastSlash+1 : beforeName])
	return name
}

func getFileExt(filename string) string {
	var i int
	var lastDot int

	for i = len(filename) - 1; i > len(filename)-7; i-- {
		if filename[i] == 46 {
			lastDot = i
			break
		}
	}
	ext := string(filename[lastDot+1 : len(filename)])
	return ext
}

func fillMatrixValueLine(file []byte, index uint8, endex uint8) {

	i := index
	var caseNb uint8 = 0
	for i < endex-1 {
		MatrixValuesAsLine[caseNb], err = parseByte(fmt.Sprintf("%s", file[i:i+8]))
		check(err)
		i += 9
		caseNb++
	}
}
