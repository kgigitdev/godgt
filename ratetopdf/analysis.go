package main

import (
	"encoding/json"
	"math"
)

// ScoredMove is a single move, and its score, according to the engine.
type ScoredMove struct {
	Rank  int     `json:"rank,omitempty"`
	Move  string  `json:"move"`
	Score float64 `json:"score"`
}

// MoveAnalysis is the analysis of a single move
type MoveAnalysis struct {
	MoveNumber int          `json:"move_number"`
	Mover      string       `json:"mover"`
	FenBefore  string       `json:"fen_before"`
	FenAfter   string       `json:"fen_after"`
	BestMoves  []ScoredMove `json:"best_moves"`
	ActualMove ScoredMove   `json:"actual_move"`
}

func (ma MoveAnalysis) String() string {
	j, err := json.MarshalIndent(ma, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(j)
}

func (ma MoveAnalysis) BlunderScore() float64 {
	// For white, the blunder is the best move score less
	// the played score. For black, the blunder is the
	// played score less the best move score. This is
	// because all scores are seen from White's point of
	// view, so higher numbers are better for white and worse
	// for black, regardless of who is playing.
	blunder := 0.0
	if ma.Mover == "white" {
		blunder = ma.BestMoves[0].Score - ma.ActualMove.Score
	} else {
		blunder = ma.ActualMove.Score - ma.BestMoves[0].Score
	}

	// We have to clamp this to be no less than zero. This is
	// because if a played move is not in the top MultiPV moves
	// found by the engine, we analyse that move separately. Since
	// explicit analysis is faster than MultiPV analysis, it
	// sometimes finds that the played move is better than the
	// moves it found, since it was able to search deeper in the
	// allocated time. This is usually only a problem when doing
	// very quick game analysis, like 3 seconds per move.
	blunder = math.Max(blunder, 0.0)
	return blunder
}

// GameAnalysis is the analysis of an entire game
type GameAnalysis []MoveAnalysis
