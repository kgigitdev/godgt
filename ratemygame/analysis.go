package main

import "encoding/json"

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
	Fen        string       `json:"fen"`
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

// GameAnalysis is the analysis of an entire game
type GameAnalysis []MoveAnalysis
