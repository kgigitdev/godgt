package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"

	"github.com/jung-kurt/gofpdf"
	"github.com/kgigitdev/godgt"
)

func main() {
	var ga GameAnalysis
	infile := os.Args[1]
	outfile := os.Args[2]

	ifh, err := os.Open(infile)
	if err != nil {
		log.Fatal(err)
	}
	analysisJSON, err := ioutil.ReadAll(ifh)
	if err != nil {
		log.Fatal(err)
	}
	ifh.Close()
	json.Unmarshal(analysisJSON, &ga)

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Courier", "", 10)

	pageWidth, pageHeight, _ := pdf.PageSize(0)
	// log.Print(pageWidth, pageHeight, units)

	numColumns := 2
	numColumnsF := float64(numColumns)

	outerMargin := 0.10
	innerMargin := 0.05

	outerMarginWidth := pageWidth * outerMargin
	innerMarginWidth := pageWidth * innerMargin

	writeableWidth := (pageWidth - (2.0 * outerMarginWidth) -
		innerMarginWidth*(numColumnsF-1.0))

	columnWidth := writeableWidth / numColumnsF

	// Now the heights and the y offsets

	numRows := 3
	numRowsF := float64(numRows)

	headerMargin := 0.15
	footerMargin := 0.15

	headerMarginHeight := pageHeight * headerMargin
	footerMarginHeight := pageHeight * footerMargin
	innerMarginHeight := pageHeight * innerMargin

	writeableHeight := (pageHeight -
		headerMarginHeight -
		footerMarginHeight -
		innerMarginHeight*(numRowsF-1))

	rowHeight := writeableHeight / numRowsF

	//  "1" indicates a full border, and one or more of "L", "T",
	//  "R" and "B" indicate the left, top, right and bottom sides
	//  of the border.

	var borderStr = ""

	// ln indicates where the current position should go after the
	// call. Possible values are 0 (to the right), 1 (to the
	// beginning of the next line), and 2 (below). Putting 1 is
	// equivalent to putting 0 and calling Ln() just after.

	// var ln = 1

	// alignStr specifies how the text is to be positioned within
	// the cell. Horizontal alignment is controlled by including
	// "L", "C" or "R" (left, center, right) in alignStr. Vertical
	// alignment is controlled by including "T", "M", "B" or "A"
	// (top, middle, bottom, baseline) in alignStr. The default
	// alignment is left middle.

	var alignStr = "TL"

	// fill is true to paint the cell background or false to
	// leave it transparent.

	var fill = false

	// var link = 0
	// var linkStr = ""

	rowCount := 0
	colCount := 0

	var whiteBlunders sort.Float64Slice
	var blackBlunders sort.Float64Slice

	for _, ma := range ga {
		if ma.Mover == "white" {
			whiteBlunders = append(whiteBlunders,
				ma.BlunderScore())
		} else {
			blackBlunders = append(blackBlunders,
				ma.BlunderScore())
		}
	}

	log.Print(whiteBlunders)
	log.Print(blackBlunders)

	sort.Sort(whiteBlunders)
	sort.Sort(blackBlunders)

	log.Printf("Worst white blunder: %.2f\n", whiteBlunders[0])
	log.Printf("Median white blunder: %.2f\n", getMedian(whiteBlunders))
	log.Printf("Worst black blunder: %.2f\n", blackBlunders[0])
	log.Printf("Median black blunder: %.2f\n", getMedian(blackBlunders))

	// Main drawing loop
	for _, ma := range ga {
		xoffset := (outerMarginWidth +
			innerMarginWidth*(float64(colCount)) +
			float64(colCount)*columnWidth)
		yoffset := (headerMarginHeight +
			innerMarginHeight*(float64(rowCount)) +
			float64(rowCount)*rowHeight)
		pdf.MoveTo(xoffset, yoffset)

		prefix := ""
		if ma.Mover == "black" {
			prefix = "... "
		}

		text := fmt.Sprintf("%d. %s%s\n\n", ma.MoveNumber,
			prefix, ma.ActualMove.Move)
		// text += fmt.Sprintf("Score: %.2f\n", ma.ActualMove.Score)
		moveSeen := false
		for _, bm := range ma.BestMoves {
			marker := ""
			if bm.Move == ma.ActualMove.Move {
				marker = "*"
				moveSeen = true
			}
			text += fmt.Sprintf("%-7s%5.2f%s\n", bm.Move,
				bm.Score, marker)
		}
		if !moveSeen {
			text += fmt.Sprintf("\n%-4s %5.2f*\n",
				ma.ActualMove.Move, ma.ActualMove.Score)
		}
		// html := pdf.HTMLBasicNew()
		// html.Write(4.0, text)
		pdf.MultiCell(columnWidth, 4.0, text,
			borderStr, alignStr, fill)

		blunder := 0.0
		if ma.Mover == "white" {
			blunder = ma.BestMoves[0].Score - ma.ActualMove.Score
		} else {
			blunder = ma.ActualMove.Score - ma.BestMoves[0].Score
		}

		blunderMessage := fmt.Sprintf("Blunder:  %5.2f\n", blunder)
		pdf.MoveTo(xoffset, yoffset+60.0)
		pdf.SetFont("Courier", "B", 12.0)
		setTextColor(pdf, blunder)
		pdf.MultiCell(columnWidth, 4.0, blunderMessage, borderStr, alignStr, fill)
		pdf.SetFont("Courier", "B", 12.0)
		pdf.SetTextColor(0x00, 0x00, 0x00)

		// Draw board
		drawBoard(pdf, ma.FenAfter, xoffset, yoffset)

		// pdf.Write(html)
		// pdf.Cell(columnWidth, rowHeight, text)
		colCount++
		if colCount == numColumns {
			colCount = 0
			rowCount++
			if rowCount == numRows {
				rowCount = 0
				pdf.AddPage()
			}
		}
	}
	err = pdf.OutputFileAndClose(outfile)
	if err != nil {
		log.Fatal(err)
	}
}

func setTextColor(pdf *gofpdf.Fpdf, blunder float64) {
	if blunder == 0.0 {
		pdf.SetTextColor(0x00, 0x80, 0x00)
	} else if blunder < 0.2 {
		pdf.SetTextColor(0x20, 0x70, 0x00)
	} else if blunder < 0.4 {
		pdf.SetTextColor(0x40, 0x60, 0x00)
	} else if blunder < 0.6 {
		pdf.SetTextColor(0x60, 0x50, 0x00)
	} else if blunder < 0.8 {
		pdf.SetTextColor(0x80, 0x40, 0x00)
	} else if blunder < 1.0 {
		pdf.SetTextColor(0xa0, 0x30, 0x00)
	} else if blunder < 1.2 {
		pdf.SetTextColor(0xc0, 0x20, 0x00)
	} else if blunder < 1.4 {
		pdf.SetTextColor(0xe0, 0x10, 0x00)
	} else {
		pdf.SetTextColor(0xff, 0x00, 0x00)
	}

}

func drawBoard(pdf *gofpdf.Fpdf, fen string, xbase float64, ybase float64) {
	fs := godgt.FS(false)
	size := 128
	flow := false

	// output := image.NewRGBA(image.Rect(0, 0, size*8, size*8))
	// white := color.RGBA{255, 255, 255, 255}
	// draw.Draw(output, output.Bounds(), &image.Uniform{white},
	// image.ZP, draw.Src)

	imageOptions := gofpdf.ImageOptions{
		ImageType: "PNG",
	}

	simpleRows := godgt.SimpleBoardFromFen(fen)
	for iRow, simpleRow := range simpleRows {
		for iCol, fenChar := range simpleRow {
			fenString := string(fenChar)
			imagePath := godgt.GetImagePath(fenString, iCol, iRow, size)
			if imagePath == "" {
				continue
			}

			file, err := fs.Open(imagePath)
			if err != nil {
				log.Fatal("Error opening " + imagePath +
					": " + err.Error())
			}

			pdf.RegisterImageOptionsReader(imagePath, imageOptions, file)

			// Magic numbers determined by trial and error;
			// need to work out how to compute these properly.
			x := xbase + 35.0 + float64(iCol)*5.75
			y := ybase + float64(iRow)*5.75
			// w := 16.0
			// h := 16.0
			pdf.MoveTo(x, y)
			pdf.ImageOptions(imagePath, x, y, -600, -600, flow, imageOptions, 0, "")
		}
	}
}

func getMedian(values sort.Float64Slice) float64 {
	numEntries := len(values)
	// FIXME: Deal with empty slices, etc.
	var median float64
	if numEntries%2 == 1 {
		// If it's an odd length, subtract one, halve it, and
		// use half+1 as the index.
		idx := ((numEntries - 1) / 2) + 1
		median = values[idx]
	} else {
		// If it's an even length, halve it, take the indices
		// at the half and half+1, and return the average of
		// the values at those indices.
		idx1 := numEntries / 2
		idx2 := idx1 + 1
		median = (values[idx1] + values[idx2]) / 2.0
	}
	return median
}
