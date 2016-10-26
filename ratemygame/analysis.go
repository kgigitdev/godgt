package main

import "encoding/json"

// ScoredMove is a single move, and its score, according to the engine.
type ScoredMove struct {
	San   string  `json:"san"`
	Rank  int     `json:"rank"`
	Score float64 `json:"score"`
}

// MoveAnalysis is the analysis of a single move
type MoveAnalysis struct {
	Move      int    `json:"move"`
	Mover     string `json:"mover"`
	BestMoves []ScoredMove
	Actual    ScoredMove
	Fen       string `json:"fen"`
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
