package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/jung-kurt/gofpdf"
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
			text += fmt.Sprintf("%-4s %.2f%s\n", bm.Move,
				bm.Score, marker)
		}
		if !moveSeen {
			text += fmt.Sprintf("\n%-4s %.2f*\n",
				ma.ActualMove.Move, ma.ActualMove.Score)
		}
		// html := pdf.HTMLBasicNew()
		// html.Write(4.0, text)
		pdf.MultiCell(columnWidth, 4.0, text,
			borderStr, alignStr, fill)
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
