package godgt

import (
	"log"

	"github.com/malbrecht/chess"
)

func (dgtboard *DgtBoard) handleBoardDump(arguments []byte) (*Message, error) {
	log.Println("DGT_BOARD_DUMP")
	board := &chess.Board{}
	for squareIndex, gdtPieceCode := range arguments {
		chessPiece := dgtboard.getChessPieceByGdtPieceCode(gdtPieceCode)
		square := dgtboard.getChessSquareFromIndex(squareIndex)
		board.Piece[square] = chessPiece
	}

	// Assume that we can castle unless it's clear that we can't.
	if board.Piece[chess.E1] == chess.WK {
		if board.Piece[chess.A1] == chess.WR {
			board.CastleSq[chess.WhiteOOO] = chess.A1
		} else {
			board.CastleSq[chess.WhiteOOO] = chess.NoSquare
		}
		if board.Piece[chess.H1] == chess.WR {
			board.CastleSq[chess.WhiteOO] = chess.H1
		} else {
			board.CastleSq[chess.WhiteOO] = chess.NoSquare
		}
	} else {
		board.CastleSq[chess.WhiteOOO] = chess.NoSquare
		board.CastleSq[chess.WhiteOO] = chess.NoSquare
	}

	if board.Piece[chess.E8] == chess.WK {
		if board.Piece[chess.A8] == chess.WR {
			board.CastleSq[chess.BlackOOO] = chess.A8
		} else {
			board.CastleSq[chess.BlackOOO] = chess.NoSquare
		}
		if board.Piece[chess.H8] == chess.WR {
			board.CastleSq[chess.BlackOO] = chess.H8
		} else {
			board.CastleSq[chess.BlackOO] = chess.NoSquare
		}
	} else {
		board.CastleSq[chess.BlackOOO] = chess.NoSquare
		board.CastleSq[chess.BlackOO] = chess.NoSquare
	}

	boardUpdate := NewBoardUpdate(board)
	return NewBoardUpdateMessage(boardUpdate), nil
}

func (dgtboard *DgtBoard) getChessSquareFromIndex(dgtSquare int) chess.Sq {
	// chess.Square numbers a1=0, b1=1, ..., h1=7, a2=8, ...,
	// h8=63. DGT update messages number them a8=0, h8=7, a7=8,
	// h1=63, etc.
	// Compute the file (0-7)
	fileIndex := dgtSquare % 8

	// Compute the rank, from white's POV, as 0-7, which is the
	// same approach as used by chess.Square
	rankIndex := 7 - ((dgtSquare - fileIndex) / 8)

	return chess.Square(fileIndex, rankIndex)
}

func (dgtboard *DgtBoard) getChessPieceByGdtPieceCode(gdtPieceCode byte) chess.Piece {
	switch gdtPieceCode {
	case WPAWN:
		return chess.WP
	case WKNIGHT:
		return chess.WN
	case WBISHOP:
		return chess.WB
	case WROOK:
		return chess.WR
	case WQUEEN:
		return chess.WQ
	case WKING:
		return chess.WK
	case BPAWN:
		return chess.BP
	case BKNIGHT:
		return chess.BN
	case BBISHOP:
		return chess.BB
	case BROOK:
		return chess.BR
	case BQUEEN:
		return chess.BQ
	case BKING:
		return chess.BK
	case EMPTY:
		return 0
	default:
		panic("Bad piece")
	}
}
