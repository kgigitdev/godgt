package godgt

import (
	"fmt"
	"log"
	"strings"

	"github.com/malbrecht/chess"
)

func fmtpsq(piece chess.Piece, square chess.Sq) string {
	return fmt.Sprintf("%c@%s", chess.PieceLetters[piece], square)
}

type MessageProcessor struct {
	Board *chess.Board

	// The first piece lifted; we use this to help us
	// signal special configurations from the board.
	FirstPieceUp *FieldUpdate

	// Could use interface{}, yadda yadda
	PiecesInTheAir map[chess.Sq]chess.Piece
	PiecesDropped  map[chess.Sq]chess.Piece

	// Need to add a move channel here.
}

func NewMessageProcessor() *MessageProcessor {
	return &MessageProcessor{
		PiecesInTheAir: make(map[chess.Sq]chess.Piece),
		PiecesDropped:  make(map[chess.Sq]chess.Piece),
	}
}

func (mp *MessageProcessor) Air() string {
	var elems []string
	for square, piece := range mp.PiecesInTheAir {
		elems = append(elems, fmtpsq(piece, square))
	}
	return "[" + strings.Join(elems, ", ") + "]"
}

func (mp *MessageProcessor) Dropped() string {
	var elems []string
	for square, piece := range mp.PiecesDropped {
		elems = append(elems, fmtpsq(piece, square))
	}
	return "[" + strings.Join(elems, ", ") + "]"
}

func (mp *MessageProcessor) LogAirState() {
	log.Println(mp.AirState())
}

func (mp *MessageProcessor) AirState() string {
	return fmt.Sprintf("{ %s / %s }",
		mp.Air(), mp.Dropped())
}

func (mp *MessageProcessor) ProcessMessage(m *Message) {
	if m.BoardUpdate != nil {
		mp.processBoardUpdate(m)
	} else if m.FieldUpdate != nil {
		mp.processFieldUpdate(m)
	} else if m.TimeUpdate != nil {
		mp.processTimeUpdate(m)
	} else if m.InfoUpdate != nil {
		mp.processInfoUpdate(m)
	} else {
		// Panic? Ignore?
		panic("Received bad message.")
	}
}

func (mp *MessageProcessor) processBoardUpdate(m *Message) {
	// If this is the first time we have received a board update,
	// store it as the first position, and assume that all subsequent
	// updates are relative to it. Note that we can configure non-starting
	// positions using special signalling moves from the board.
	if mp.Board == nil {
		mp.Board = m.BoardUpdate.Board
		log.Println("Received initial board update.")
		log.Println(mp.Board.Fen())
	} else {
		// In future, maybe allow special coded moves to force
		// a board update so we can be sure that our board is
		// correct.
		log.Printf("Ignoring board update.")
	}
}

func (mp *MessageProcessor) processFieldUpdate(m *Message) {
	fieldUpdate := m.FieldUpdate

	if fieldUpdate.Piece == chess.NoPiece {
		mp.processPieceLift(fieldUpdate)
	} else {
		mp.processPieceDrop(fieldUpdate)
	}
}

func (mp *MessageProcessor) processPieceLift(fieldUpdate *FieldUpdate) {
	// It's a lift, so we have to go and look to see
	// what was actually there prior to the lift.

	// piece := fieldUpdate.Piece
	square := fieldUpdate.Square

	pieceLifted := mp.GetPieceLifted(square)

	airUpdate := FieldUpdate{
		Square: square,
		Piece:  pieceLifted,
	}
	log.Println("Saw piece lift")
	mp.PiecesInTheAir[square] = pieceLifted
	// Delete whatever we had detected dropped onto that square;
	// it isn't there any more.
	delete(mp.PiecesDropped, square)
	mp.simplifyAirState()

	// Update FirstPieceUp if it's nil. If there is already a
	// piece lifted, ignore subsequent piece lifts, since that can
	// easily happen as part of a capture, or a castle, or a capture
	// EP, or a fumble.
	if mp.FirstPieceUp == nil {
		mp.FirstPieceUp = &airUpdate
	}
}

// GetPieceLifted returns the piece that was lifted from a square.  We
// can't look directly into the Board because until we have decoded an
// entire move it doesn't get updated. Image the following sequence of
// meaningless and illegal physical piece moves by White:
//
// e2f3, Ke2, Ke1, f3e2
//
// At the end of the sequence, all the pieces have returned back to
// their starting square, meaning that White's next legal move
// *should* be recognised and accepted. However, when White moves his
// King from e2 back to e1, we need to "know" that the piece being
// lifted is the King, not the pawn originally on e2! Don't forget, we
// don't actually receive a message saying, "White King lifted"; we
// receive a messsage saying, "e2 is now empty", and we have to work
// out what piece was moved.
func (mp *MessageProcessor) GetPieceLifted(square chess.Sq) chess.Piece {
	piece, ok := mp.PiecesDropped[square]
	if ok {
		return piece
	}
	return mp.Board.Piece[square]
}

func (mp *MessageProcessor) processPieceDrop(fieldUpdate *FieldUpdate) {
	// It's a piece drop (that is, a piece placed down on a square).
	// In this case, we have three possibilities:
	//
	// 1. The last piece lifted was dropped back onto the same square.
	//    This is usually a special signal, but only if no other
	//    pieces are in the air.
	//
	// 2. The last piece lifted as dropped back onto another square.
	//    This is assumed to be the move. However, it can only be a valid
	//    move if (a) there is at most one piece in the air (which would
	//    make it a capture), and (b) the move itself is valid.
	//
	// 3. A piece other than the last piece lifted is dropped onto a
	//    square. This is assumed to be noise, and is ignored. This will
	//    always happen during castling; it is assumed that the movement
	//    of the king itself entirely defines the castling move.

	// Note that this will update whatever we had previously detected
	// as being dropped onto this square. But that's OK, because it's
	// true.
	mp.PiecesDropped[fieldUpdate.Square] = fieldUpdate.Piece

	log.Printf("Piece drop: %s\n", fieldUpdate.ToString())
	// Now we need to simplify the air state to detect when
	// accidentally lifted pieces have been returned to their squares;
	// this allows us to recover from a bad state. However, a single
	// piece being lifted and returned to the same square, with no
	// other pieces being lifted in the meantime is a special signal.
	mp.processSpecialSignal()
	mp.simplifyAirState()
	mp.IsMovePseudoLegal()
}

func (mp *MessageProcessor) GetSpecialAirState() (chess.Sq, chess.Piece) {
	if len(mp.PiecesInTheAir) != 1 {
		return chess.NoSquare, chess.NoPiece
	}
	if len(mp.PiecesDropped) != 1 {
		return chess.NoSquare, chess.NoPiece
	}
	for square, uppiece := range mp.PiecesInTheAir {
		downpiece, ok := mp.PiecesDropped[square]
		if !ok {
			// Square wasn't the same
			return chess.NoSquare, chess.NoPiece
		}
		if uppiece != downpiece {
			return chess.NoSquare, chess.NoPiece
		}

		return square, uppiece
	}

	return chess.NoSquare, chess.NoPiece
}

func (mp *MessageProcessor) processSpecialSignal() {
	square, piece := mp.GetSpecialAirState()

	if square == chess.NoSquare || piece == chess.NoPiece {
		return
	}

	// King lifted and dropped onto same square: make that side be
	// the side to play.
	if piece == chess.WK {
		mp.Board.SideToMove = chess.White
		log.Println("Special signal: White to play")
	} else if piece == chess.BK {
		mp.Board.SideToMove = chess.Black
		log.Println("Special signal: Black to play")
	} else if piece == chess.WR && square == chess.A1 {
		// Toggle White's queenside castling rights
		if mp.Board.CastleSq[chess.WhiteOOO] == chess.NoSquare {
			log.Println("White may castle queenside")
			mp.Board.CastleSq[chess.WhiteOOO] = chess.A1
		} else if mp.Board.CastleSq[chess.WhiteOOO] == chess.A1 {
			log.Println("White may NOT castle queenside")
			mp.Board.CastleSq[chess.WhiteOOO] = chess.NoSquare
		}
	} else if piece == chess.WR && square == chess.H1 {
		// Toggle White's kingside castling rights
		if mp.Board.CastleSq[chess.WhiteOO] == chess.NoSquare {
			log.Println("White may castle kingside")
			mp.Board.CastleSq[chess.WhiteOO] = chess.H1
		} else if mp.Board.CastleSq[chess.WhiteOO] == chess.H1 {
			log.Println("White may NOT castle kingside")
			mp.Board.CastleSq[chess.WhiteOO] = chess.NoSquare
		}
	} else if piece == chess.BR && square == chess.A8 {
		// Toggle Black's queenside castling rights
		if mp.Board.CastleSq[chess.BlackOOO] == chess.NoSquare {
			log.Println("Black may castle queenside")
			mp.Board.CastleSq[chess.BlackOOO] = chess.A8
		} else if mp.Board.CastleSq[chess.BlackOOO] == chess.A8 {
			log.Println("Black may NOT castle queenside")
			mp.Board.CastleSq[chess.BlackOOO] = chess.NoSquare
		}
	} else if piece == chess.BR && square == chess.H8 {
		// Toggle Black's kingside castling rights
		if mp.Board.CastleSq[chess.BlackOO] == chess.NoSquare {
			log.Println("Black may castle kingside")
			mp.Board.CastleSq[chess.BlackOO] = chess.H8
		} else if mp.Board.CastleSq[chess.BlackOO] == chess.H8 {
			log.Println("Black may NOT castle kingside")
			mp.Board.CastleSq[chess.BlackOO] = chess.NoSquare
		}
	} else {
		log.Println("Couldn't decode special signal: " + fmtpsq(piece, square))
	}
}

func (mp *MessageProcessor) simplifyAirState() {
	// Delete any field updates common to both up and down,
	// indicating a piece that has been lifted and returned
	// to the same square. This allows us to recover from
	// bad updates, and also allows us to deal with updates
	// received out of order.
	init := mp.AirState()
	for square, uppiece := range mp.PiecesInTheAir {
		downpiece, ok := mp.PiecesDropped[square]
		if !ok {
			// Nothing at that square
			continue
		}
		if uppiece != downpiece {
			// Pieces are at the same square, but are not the same
			// piece.
			continue
		}
		delete(mp.PiecesInTheAir, square)
		delete(mp.PiecesDropped, square)
	}

	for square, uppiece := range mp.PiecesInTheAir {
		if uppiece == chess.NoPiece {
			// No point in holding on to this. Also, this is safe in Go.
			delete(mp.PiecesInTheAir, square)
		}
	}

	for square, downpiece := range mp.PiecesDropped {
		if downpiece == chess.NoPiece {
			// No point in holding on to this. Also, this is safe in Go.
			delete(mp.PiecesDropped, square)
		}
	}

	final := mp.AirState()
	if init != final {
		log.Printf("Simplifying: %s -> %s\n", init, final)
	} else {
		mp.LogAirState()
	}
}

func (mp *MessageProcessor) processUnrelatedPieceDrop(piece chess.Piece, square chess.Sq) {
	log.Println("Detected unrelated piece drop")
}

func (mp *MessageProcessor) processTimeUpdate(m *Message) {
	log.Println("TODO: Implement processTimeUpdate()")
}

func (mp *MessageProcessor) processInfoUpdate(m *Message) {
	log.Println("TODO: Implement processInfoUpdate()")
}

func (mp *MessageProcessor) IsMovePseudoLegal() bool {
	totalEnemyPiecesLifted := 0
	totalOwnPiecesLifted := 0
	totalEnemyPiecesDropped := 0
	totalOwnPiecesDropped := 0
	for _, piece := range mp.PiecesInTheAir {
		if piece.Color() == mp.Board.SideToMove {
			totalOwnPiecesLifted++
		} else {
			totalEnemyPiecesLifted++
		}
	}

	for _, piece := range mp.PiecesDropped {
		if piece.Color() == mp.Board.SideToMove {
			totalOwnPiecesDropped++
		} else {
			totalEnemyPiecesDropped++
		}
	}

	// For a move to be "even pseudolegal":

	// 1. Exactly one (moving, capturing) or two (castling) own pieces
	// must have been lifted, and the same number dropped.
	if totalOwnPiecesLifted == 0 {
		log.Printf("Not a move: no pieces lifted.")
		return false
	}

	if totalOwnPiecesLifted > 2 {
		log.Printf("Not a move: more than two pieces lifted.")
		return false
	}

	if totalOwnPiecesLifted > totalOwnPiecesDropped {
		log.Printf("Not a move: some pieces are still in the air")
		return false
	}

	if totalOwnPiecesLifted < totalOwnPiecesDropped {
		log.Printf("Not a move: some pieces have appeared from thin air")
		return false
	}

	// 2. Zero enemy pieces must have been dropped.

	if totalEnemyPiecesDropped > 0 {
		log.Printf("Not a move: enemy piece dropped.")
		return false
	}

	// 3. Exactly ZERO or ONE enemy pieces may have been lifted (captured).

	if totalEnemyPiecesLifted > 1 {
		log.Printf("Not a move: multiple enemy pieces lifted.")
		return false
	}

	// 4. In the case of a castling move, only kings and rooks are
	// allowed to have been moved. Note that we don't bother checking
	// for the legality of the actual castling move; we let the chess
	// library do that. We just check to see if it *looks* like a
	// castling move.

	if totalOwnPiecesDropped == 2 {
		// Looks like a castling move
		if mp.Board.SideToMove == chess.White {
			// White to move: The king must have been lifted from E8
			piece, ok := mp.PiecesInTheAir[chess.E1]
			if !ok {
				log.Printf("Not a castling move: nothing moved from e1.")
				return false
			}
			if piece != chess.WK {
				log.Printf("Not a castling move: piece lifted from e1 != WK")
				return false
			}
			g1drop, kingside := mp.PiecesDropped[chess.G1]
			c1drop, queenside := mp.PiecesDropped[chess.C1]
			if kingside {
				if g1drop != chess.WK {
					log.Printf("Not a castling move: piece dropped on g1 != WK")
					return false
				}
				h1lift, ok := mp.PiecesInTheAir[chess.H1]
				if !ok {
					log.Printf("Not a castling move: no piece taken from h1")
					return false
				}
				if h1lift != chess.WR {
					log.Printf("Not a castling move: piece taken from h1 != WR")
					return false
				}
				f1drop, ok := mp.PiecesDropped[chess.F1]
				if !ok {
					log.Printf("Not a castling move: no piece dropped on f1")
					return false
				}
				if f1drop != chess.WR {
					log.Printf("Not a castling move: piece dropped on f1 != WR")
					return false
				}
				log.Println("Detecting White O-O")
				return true
			} else if queenside {
				if c1drop != chess.WK {
					log.Printf("Not a castling move: piece dropped on c1 != WK")
					return false
				}
				a1lift, ok := mp.PiecesInTheAir[chess.A1]
				if !ok {
					log.Printf("Not a castling move: no piece taken from a1")
					return false
				}
				if a1lift != chess.WR {
					log.Printf("Not a castling move: piece taken from a1 != WR")
					return false
				}
				d1drop, ok := mp.PiecesDropped[chess.D1]
				if !ok {
					log.Printf("Not a castling move: no piece dropped on d1")
					return false
				}
				if d1drop != chess.WR {
					log.Printf("Not a castling move: piece dropped on d1 != WR")
					return false
				}
				log.Println("Detecting White O-O-O")
				return true
			}
		} else {
			// Black to move: The king must have been lifted from E8
			piece, ok := mp.PiecesInTheAir[chess.E8]
			if !ok {
				log.Printf("Not a castling move: nothing moved from e8.")
				return false
			}
			if piece != chess.BK {
				log.Printf("Not a castling move: piece lifted from e8 != WK")
				return false
			}
			g8drop, kingside := mp.PiecesDropped[chess.G8]
			c8drop, queenside := mp.PiecesDropped[chess.C8]
			if kingside {
				if g8drop != chess.BK {
					log.Printf("Not a castling move: piece dropped on g8 != BK")
					return false
				}
				h8lift, ok := mp.PiecesInTheAir[chess.H8]
				if !ok {
					log.Printf("Not a castling move: no piece taken from h8")
					return false
				}
				if h8lift != chess.BR {
					log.Printf("Not a castling move: piece taken from h8 != BR")
					return false
				}
				f8drop, ok := mp.PiecesDropped[chess.F8]
				if !ok {
					log.Printf("Not a castling move: no piece dropped on f8")
					return false
				}
				if f8drop != chess.BR {
					log.Printf("Not a castling move: piece dropped on f8 != BR")
					return false
				}
				log.Println("Detecting Black O-O")
				return true
			} else if queenside {
				if c8drop != chess.BK {
					log.Printf("Not a castling move: piece dropped on c8 != BK")
					return false
				}
				a8lift, ok := mp.PiecesInTheAir[chess.A8]
				if !ok {
					log.Printf("Not a castling move: no piece taken from a8")
					return false
				}
				if a8lift != chess.BR {
					log.Printf("Not a castling move: piece taken from a8 != BR")
					return false
				}
				d8drop, ok := mp.PiecesDropped[chess.D8]
				if !ok {
					log.Printf("Not a castling move: no piece dropped on d8")
					return false
				}
				if d8drop != chess.BR {
					log.Printf("Not a castling move: piece dropped on d8 != BR")
					return false
				}
				log.Println("Detecting Black O-O-O")
				return true
			}
		}
	}

	// 2. If we're not castling, exactly ONE piece needs to have been
	// moved, which means ONE lift square and ONE drop square, and the
	// piece needs to be the same piece.

	if totalOwnPiecesDropped != 1 {
		log.Println("Reached unreachable code")
		return false
	}

	var sourceSquare chess.Sq
	var sourcePiece chess.Piece
	var targetSquare chess.Sq
	var targetPiece chess.Piece

	for square, piece := range mp.PiecesInTheAir {
		if piece.Color() != mp.Board.SideToMove {
			continue
		}
		sourceSquare = square
		sourcePiece = piece
	}

	for square, piece := range mp.PiecesDropped {
		if piece.Color() != mp.Board.SideToMove {
			continue
		}
		targetSquare = square
		targetPiece = piece
	}

	if sourcePiece != targetPiece {
		log.Println("Piece changed in midair")
		return false
	}

	log.Printf("Detected move: %s -> %s\n", fmtpsq(sourcePiece,
		sourceSquare), fmtpsq(targetPiece, targetSquare))

	uciMove := fmt.Sprintf("%s%s", sourceSquare, targetSquare)
	log.Printf("Submitting move [%s]\n", uciMove)
	parsedMove, err := mp.Board.ParseMove(uciMove)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Move accepted!")
		log.Println(parsedMove.San(mp.Board))

		log.Println("Clearing down field updates ...")
		mp.PiecesInTheAir = make(map[chess.Sq]chess.Piece)
		mp.PiecesDropped = make(map[chess.Sq]chess.Piece)

		log.Println("Switching player ...")
		if mp.Board.SideToMove == chess.White {
			mp.Board.SideToMove = chess.Black
		} else {
			mp.Board.SideToMove = chess.White
		}
	}

	return true
}

func (mp *MessageProcessor) processMove(piece chess.Piece, square chess.Sq) {
	log.Println("Detected move")
}
