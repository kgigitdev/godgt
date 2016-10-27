package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"text/tabwriter"
	"time"

	"github.com/malbrecht/chess"
	"github.com/malbrecht/chess/engine"
	"github.com/malbrecht/chess/engine/uci"
	"github.com/malbrecht/chess/pgn"
)

// GameRater is the outer game rating class (really just a wrapper for
// lower level things like a PGN reader and a UCI interface to an
// engine)
type GameRater struct {
	opts          Opts
	pgnFileHandle *os.File
	pgntext       string
	db            pgn.DB
	game          *pgn.Game
	engine        *uci.Engine
	engineOptions map[string]engine.Option
	node          *pgn.Node
	board         *chess.Board
	infoChannel   <-chan engine.Info
	allPvs        map[int]*engine.Pv
	analysis      GameAnalysis
	bestMove      float64
	searchStart   time.Time
}

// NewGameRater creates and returns a pointer to a new GameRater
func NewGameRater(opts Opts) *GameRater {
	g := &GameRater{
		opts: opts,
	}
	g.clearPvs()
	return g
}

// Run is the main entry point for GameRater
func (g *GameRater) Run() {
	g.openPgnFile()
	g.readPgnFile()
	g.createNewEmptyDatabase()
	g.parsePgnText()
	g.createEngine()
	defer g.engine.Quit()
	g.readEngineOptions()
	g.maybePrintEngineOptions()
	g.setMPVOption(g.opts.MultiPV)
	g.extractRootGameNode()
	g.processAllMoves()
	g.writeOutputFile()
}

func (g *GameRater) openPgnFile() {
	file := g.opts.PgnFile
	fh, err := os.Open(file)

	if err != nil {
		log.Fatal(err)
	}

	g.pgnFileHandle = fh
}

func (g *GameRater) readPgnFile() {
	pgntext, err := ioutil.ReadAll(g.pgnFileHandle)
	if err != nil {
		log.Fatal(err)
	}
	g.pgntext = string(pgntext)
}

func (g *GameRater) createNewEmptyDatabase() {
	g.db = pgn.DB{}
}

func (g *GameRater) parsePgnText() {
	// Note that pgn.DB can contain any number of games; we
	// are only ever interested in the first.
	errors := g.db.Parse(g.pgntext)
	if len(errors) > 0 {
		for _, err := range errors {
			log.Println(err)
		}
		os.Exit(1)
	}
	g.game = g.db.Games[0]
	err := g.db.ParseMoves(g.game)

	if err != nil {
		log.Fatal(err)
	}
}

func (g *GameRater) createEngine() {
	var sfargs []string
	// logger := log.New(os.Stdout, "", log.LstdFlags)
	var logger *log.Logger
	engine, err := uci.Run(g.opts.Engine, sfargs, logger)
	if err != nil {
		log.Fatal(err)
	}
	g.engine = engine
}

func (g *GameRater) readEngineOptions() {
	g.engineOptions = g.engine.Options()
}

func (g *GameRater) maybePrintEngineOptions() {
	if !g.opts.Verbose {
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 1, 8, 0, ' ', 0)
	for k, v := range g.engineOptions {
		fmt.Fprintln(w, k, "\t", v)
	}
	w.Flush()
}

func (g *GameRater) setMPVOption(mpv int) {
	mpvOption, ok := g.engineOptions["MultiPV"]
	if !ok {
		log.Fatal("Can't set MultiPV mode")
	}
	// Even though it's an integer field, it takes a string;
	// internally, it performs an Atoi() on it
	mpvOption.Set(fmt.Sprintf("%d", mpv))
}

func (g *GameRater) extractRootGameNode() {
	// The root node is the node after "zero" moves; it contains
	// the initial starting position, and a null move. We
	// therefore need to pull out the next move, so we can compare
	// what the engine would do in the current position with what
	// the player actually played (which is stored on the *next*
	// node)
	g.node = g.game.Root
}

func (g *GameRater) processAllMoves() {
	for {
		// Note that even though it might appear that it's would be
		// possible to move some of the method calls into
		// processOneMove(), we can't. This is because we also call
		// processOneMove() from another place, when scoring very
		// bad moves.
		if g.gameFinished() {
			break
		}
		g.extractCurrentBoard()
		g.processOneMove()
		g.processEngineResults()
		g.updateGameState()
		if g.gameFinished() {
			break
		}
	}
}

func (g *GameRater) processOneMove() {
	g.sendPositionToEngine()
	g.startSearch()
	g.clearPvs()
	g.processAllEngineMessages()
}

func (g *GameRater) extractCurrentBoard() {
	g.board = g.node.Board
}

func (g *GameRater) sendPositionToEngine() {
	g.engine.SetPosition(g.board)
}

func (g *GameRater) startSearch() {
	var depthPerMove int
	var timePerMove int

	depthPerMove = g.opts.DepthPerMove
	timePerMove = g.opts.TimePerMove

	if g.board.MoveNr < g.opts.OpeningLength {
		// We're still in the opening, so maybe use different
		// search values.
		g.debug("In opening")
		if g.opts.OpeningDepthPerMove > 0 || g.opts.OpeningTimePerMove > 0 {
			depthPerMove = g.opts.OpeningDepthPerMove
			timePerMove = g.opts.OpeningTimePerMove
		}
	} else {
		g.debug("Out of opening")
	}

	// Sanity check: if both are zero, set the time to something
	// not entirely insane.
	if depthPerMove == 0 && timePerMove == 0 {
		timePerMove = 5
	}

	g.searchStart = time.Now()
	if depthPerMove > 0 {
		g.debug("Searching to a depth of %d", depthPerMove)
		g.infoChannel = g.engine.SearchDepth(depthPerMove)
	} else {
		g.debug("Searching for %d seconds", timePerMove)
		g.infoChannel = g.engine.SearchTime(time.Second *
			time.Duration(timePerMove))
	}
}

func (g *GameRater) processAllEngineMessages() {
	for info := range g.infoChannel {
		g.processOneEngineMessage(info)
	}
}

func (g *GameRater) processOneEngineMessage(info engine.Info) {
	if info.Err() != nil {
		log.Fatalf("%s", info.Err())
	}
	if _, ok := info.BestMove(); ok {
		// This indicates that the thinking time is over. We
		// don't need to indicate that we need to stop
		// processing messages from the engine; the engine
		// will close the channel.
		searchEnd := time.Now()
		g.debug("Move analysis took %.2f seconds.",
			searchEnd.Sub(g.searchStart).Seconds())
		return
	} else if pv := info.Pv(); pv != nil {
		g.processOnePV(pv)
		// log.Println("stats:", info.Stats())
	} else {
		// log.Println("stats:", info.Stats())
	}
}

func (g *GameRater) processOnePV(pv *engine.Pv) {
	g.allPvs[pv.Rank] = pv
}

func (g *GameRater) processEngineResults() {
	nextNode := g.node.Next
	if nextNode == nil {
		return
	}
	ma := MoveAnalysis{}

	// Since g.allPvs is actually a map, the ranks have no
	// guaranteed order. However, it's nice to have them properly
	// ordered in the JSON output, which is a list. We also need
	// to guard against the possibility that the actual played
	// move is so bad that it doesn't feature in any of the top
	// moves from the engine, in which case we need to do
	// something special to score it.
	var ranks []int
	moveToScore := make(map[string]float64)
	for rank := range g.allPvs {
		ranks = append(ranks, rank)
	}
	sort.Ints(ranks)
	for rank := range ranks {
		pv, ok := g.allPvs[rank]
		if !ok || pv == nil {
			continue
		}
		if len(pv.Moves) == 0 {
			continue
		}
		engineMove := pv.Moves[0]
		san := engineMove.San(g.board)
		score := pv.Score
		fscore := float64(score) / 100.0

		sm := ScoredMove{
			// Add 1 so the best move has rank 1. Also to
			// prevent rank 0 from being reaped by the
			// "omitempty" JSON directive.
			Rank:  rank + 1,
			Move:  san,
			Score: fscore,
		}
		ma.BestMoves = append(ma.BestMoves, sm)
		moveToScore[san] = fscore
	}

	// Now add the actual move played.
	move := nextNode.Move
	actualSan := move.San(g.board)
	fen := g.board.Fen()
	if g.board.SideToMove == 0 {
		ma.Mover = "white"
	} else {
		ma.Mover = "black"
	}
	ma.MoveNumber = g.board.MoveNr
	ma.Fen = fen

	actualScore, ok := moveToScore[actualSan]
	if !ok {
		actualScore = g.computeExplicitScore(nextNode)
	}

	actualMove := ScoredMove{
		Move:  actualSan,
		Score: actualScore,
	}
	ma.ActualMove = actualMove

	g.analysis = append(g.analysis, ma)
	g.printMoveSummary(ma)
}

func (g *GameRater) printMoveSummary(ma MoveAnalysis) {
	actualMove := ma.ActualMove
	bestMove := ma.BestMoves[0]
	var score string
	if bestMove.Move == actualMove.Move {
		score = fmt.Sprintf("%6.2f", actualMove.Score)
	} else {
		score = fmt.Sprintf("%6.2f (%-6s = %6.2f)", actualMove.Score,
			bestMove.Move, bestMove.Score)
	}
	if ma.Mover == "white" {
		log.Printf("%2d. %-6s %-6s %s\n", ma.MoveNumber,
			ma.ActualMove.Move, "", score)
	} else {
		log.Printf("%2d. %-6s %-6s %s\n", ma.MoveNumber,
			"...", ma.ActualMove.Move, score)
	}
}

func (g *GameRater) computeExplicitScore(node *pgn.Node) float64 {
	// We need to process a bad move explicitly. The way we do
	// this is to feed the engine the position AFTER the move and
	// look at the score for the very best move by the *opponent*.
	// Since we only want the very best score, we can set MultiPV
	// to 1 for this.
	log.Println("Computing explicit score")
	g.setMPVOption(1)
	g.board = node.Board
	g.processOneMove()
	g.setMPVOption(g.opts.MultiPV)
	bestPv, ok := g.allPvs[0]
	if !ok {
		// What to do here?
		return -999.0
	}
	return float64(bestPv.Score) / 100.0
}

func (g *GameRater) updateGameState() {
	if g.node != nil {
		g.node = g.node.Next
	}
}

func (g *GameRater) gameFinished() bool {
	return g.node == nil
}

func (g *GameRater) clearPvs() {
	g.allPvs = make(map[int]*engine.Pv)
}

func (g *GameRater) writeOutputFile() {
	j, err := json.MarshalIndent(g.analysis, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	var oh *os.File
	if g.opts.OutputFile == "-" {
		oh = os.Stdout
	} else {
		oh, err = os.Create(g.opts.OutputFile)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer oh.Close()
	oh.Write(j)
	oh.WriteString("\n")
}

func (g *GameRater) debug(format string, args ...interface{}) {
	format = format + "\n"
	if g.opts.Verbose {
		log.Printf(format, args...)
	}
}
