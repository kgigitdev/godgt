package godgt

import (
	"strconv"
	"strings"
)

// CharToUnicode converts a simple "FEN" character (that is, Q=white queen,
// p=black pawn, etc) into a single-rune unicode string containing the
// unicode figurine for that piece.
func FenCharToFigurine(piece string) string {
	switch piece {
	case "P":
		return "♙"
	case "N":
		return "♘"
	case "B":
		return "♗"
	case "R":
		return "♖"
	case "Q":
		return "♕"
	case "K":
		return "♔"
	case "p":
		return "♟"
	case "n":
		return "♞"
	case "b":
		return "♝"
	case "r":
		return "♜"
	case "q":
		return "♛"
	case "k":
		return "♚"
	case " ":
		return " "
	default:
		return "?"
	}
}

func GetSquare(iCol int, iRow int) string {
	// iCol, iRow both counted from zero
	if (iCol+iRow)%2 == 0 {
		// return "."
		return "▨"
	} else {
		return " "
	}
}

func SimpleBoardFromFen(fen string) []string {
	fen = strings.TrimSpace(fen)
	if strings.Contains(fen, " ") {
		// It's a "full FEN"; that is, with all the extra
		// things like castling rights, en passant and so
		// forth. We don't care about these for simple
		// board representation.
		elems := strings.Fields(fen)
		fen = elems[0]
	}
	rows := strings.Split(fen, "/")
	var figurineRows []string
	for _, row := range rows {
		var figurineRow string
		for _, fenChar := range row {
			skip, err := strconv.Atoi(string(fenChar))
			if err == nil {
				// It's a number; skip that many spaces.
				for i := 0; i < skip; i++ {
					figurineRow = figurineRow + " "
				}
			} else {
				// It's a piece; add that piece.
				figurine := string(fenChar)
				figurineRow = figurineRow + figurine
			}
		}
		if len(figurineRow) < 8 {
			// Fill in the end of the row with spaces.
			for i := 0; i < 8-len(figurineRow); i++ {
				figurineRow = figurineRow + " "
			}
		}
		figurineRows = append(figurineRows, figurineRow)
	}
	return figurineRows
}

func UnicodeBoardFromFen(fen string) []string {
	simpleRows := SimpleBoardFromFen(fen)
	var unicodeRows []string
	for iRow, simpleRow := range simpleRows {
		var unicodeRow string
		for iCol, fenChar := range simpleRow {
			fenString := string(fenChar)
			if fenString == " " {
				fenString = GetSquare(iCol, iRow)
			}
			unicodeRow = unicodeRow + fenString
		}
		unicodeRows = append(unicodeRows, unicodeRow)
	}
	return unicodeRows
}
