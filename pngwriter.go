package godgt

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"log"
	"os"
)

func WritePng(fen string, size int, filename string) {
	w, err := os.Create(filename)
	if err != nil {
		log.Fatal("Failed to open file for writing.")
	}
	defer w.Close()
	WriteBoardAsPng(fen, size, w)
}

func GetFigurineName(fenChar string) string {
	// JUST for pieces; we deal with other stuff later,
	// including empty squares and stuff.
	switch fenChar {
	case "P":
		return "WP"
	case "N":
		return "WN"
	case "B":
		return "WB"
	case "R":
		return "WR"
	case "Q":
		return "WQ"
	case "K":
		return "WK"
	case "p":
		return "BP"
	case "n":
		return "BN"
	case "b":
		return "BB"
	case "r":
		return "BR"
	case "q":
		return "BQ"
	case "k":
		return "BK"
	default:
		return ""
	}
}

func GetImageName(fenChar string, iCol int, iRow int) string {
	var squareColour string
	if (iCol+iRow)%2 == 0 {
		// Dark
		squareColour = "D"
	} else {
		// Light
		squareColour = "L"
	}
	figurineName := GetFigurineName(fenChar)
	if figurineName == "" {
		// It's an empty square
		if squareColour == "L" {
			// Nothing to do; we don't draw an image for
			// a light, empty square.
			return ""
		} else {
			return "EMPTYD"
		}
	} else {
		return figurineName + squareColour
	}
}

func GetImagePath(fenChar string, iCol int, iRow int, size int) string {
	imageName := GetImageName(fenChar, iCol, iRow)
	if imageName == "" {
		return ""
	}
	return fmt.Sprintf("/assets/images/%d/%s.png", size, imageName)
}

func WriteBoardAsPng(fen string, size int, w io.Writer) {
	// Get a handle onto the local embedded assets
	fs := FS(false)
	output := image.NewRGBA(image.Rect(0, 0, size*8, size*8))
	white := color.RGBA{255, 255, 255, 255}
	draw.Draw(output, output.Bounds(), &image.Uniform{white},
		image.ZP, draw.Src)

	simpleRows := SimpleBoardFromFen(fen)
	for iRow, simpleRow := range simpleRows {
		for iCol, fenChar := range simpleRow {
			fenString := string(fenChar)
			imagePath := GetImagePath(fenString, iCol, iRow, size)
			if imagePath == "" {
				continue
			}
			file, err := fs.Open(imagePath)
			if err != nil {
				log.Fatal("Error opening " + imagePath +
					": " + err.Error())
			}
			oneImage, err := png.Decode(file)
			if err != nil {
				log.Fatal(err)
			}

			x := size * iCol
			y := size * iRow
			r := image.Rectangle{
				image.Pt(x, y),
				image.Pt(x+size, y+size)}
			draw.Draw(output, r, oneImage, image.ZP, draw.Over)
		}
	}
	png.Encode(w, output)
}
