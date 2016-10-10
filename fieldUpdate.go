package godgt

import (
	"fmt"

	"github.com/malbrecht/chess"
)

// FieldUpdate encapsulates a raw field update from a DGT board.
// It's not a full chess move. It consists of detecting a piece
// (including the special piece EMPTY) moving into a square.
type FieldUpdate struct {
	Square chess.Sq
	Piece  chess.Piece
}

func NewFieldUpdate(square chess.Sq, piece chess.Piece) *FieldUpdate {
	return &FieldUpdate{
		Square: square,
		Piece:  piece,
	}
}

func (fu *FieldUpdate) ToString() string {
	return fmt.Sprintf("%c@%s", chess.PieceLetters[fu.Piece], fu.Square)
}

func (fu *FieldUpdate) PieceLetter() string {
	return string(chess.PieceLetters[fu.Piece])
}
