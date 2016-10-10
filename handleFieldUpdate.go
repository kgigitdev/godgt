package godgt

import (
	"log"

	"github.com/malbrecht/chess"
)

func (dgtboard *DgtBoard) handleFieldUpdate(arguments []byte) (*Message, error) {
	log.Println("Field update")
	fieldNumber := arguments[0]
	gdtPieceCode := arguments[1]

	square := dgtboard.getChessSquareFromGdtFieldNumber(fieldNumber)
	piece := dgtboard.getChessPieceByGdtPieceCode(gdtPieceCode)

	fieldUpdate := NewFieldUpdate(square, piece)
	fieldUpdateMessage := NewFieldUpdateMessage(fieldUpdate)

	return fieldUpdateMessage, nil
}

func (dgtboard *DgtBoard) getChessSquareFromGdtFieldNumber(gdtFieldNumber byte) chess.Sq {
	// The field number is encoded as follows:
	// 0b00rrrfff
	// The 3 "rrr" bits denote the rank (8 values in total).
	// The 3 "fff" bits denote the file (8 values in total).
	// The "rrr" bits can be masked using 0x07.
	// The "fff" bits need to be masked using 0b00111000, or
	// 0x38, then right-shifted 3 bits. Note that r=0 actually
	// corresponds to the 8th rank (from White's POV), so we need
	// to take this into account when computing rankIndex by "flipping"
	// the order.

	// The file index, 0=a, 7=h
	fileIndex := int(gdtFieldNumber & 0x07)

	// The rank index, 0=1st rank, 7=8th rank (notice the flip)
	rankIndex := 7 - int((gdtFieldNumber&0x38)>>3)

	return chess.Square(fileIndex, rankIndex)
}
