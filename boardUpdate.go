package godgt

import "github.com/malbrecht/chess"

type BoardUpdate struct {
	Board *chess.Board
}

// BoardUpdate contains a full position update from the DGT board.
// Note that instead of making the Message rely directly on the board
// representation, we hide the specifics of the representation inside
// a small wrapper (BoardUpdate), which helps insulate the code from
// changes to the board representation.
func NewBoardUpdate(board *chess.Board) *BoardUpdate {
	return &BoardUpdate{
		Board: board,
	}
}

func (b *BoardUpdate) ToString() string {
	return b.Board.Fen()
}
