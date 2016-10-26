package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
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
	opts             Opts
	pgnFileHandle    *os.File
	outputFileHandle *os.File
	pgntext          string
	db               pgn.DB
	game             *pgn.Game
	engine           *uci.Engine
	engineOptions    map[string]engine.Option
	node             *pgn.Node
	board            *chess.Board
	infoChannel      <-chan engine.Info
	allPvs           map[int]*engine.Pv
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
	g.openOutputFile()
	g.readPgnFile()
	g.createNewEmptyDatabase()
	g.parsePgnText()
	g.createEngine()
	defer g.engine.Quit()
	g.readEngineOptions()
	g.maybePrintEngineOptions()
	g.setMPVOption()
	g.extractRootGameNode()
	g.processAllMoves()
}

func (g *GameRater) openPgnFile() {
	file := g.opts.PgnFile
	fh, err := os.Open(file)

	if err != nil {
		log.Fatal(err)
	}

	g.pgnFileHandle = fh
}

func (g *GameRater) openOutputFile() {
	if g.opts.OutputFile == "-" {
		g.outputFileHandle = os.Stdout
	} else {
		oh, err := os.Create(g.opts.OutputFile)
		if err != nil {
			log.Fatal(err)
		}
		g.outputFileHandle = oh
	}
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

func (g *GameRater) setMPVOption() {
	mpv, ok := g.engineOptions["MultiPV"]
	if !ok {
		log.Fatal("Can't set MultiPV mode")
	}
	// Even though it's an integer field, it takes a string;
	// internally, it performs an Atoi() on it
	mpv.Set(fmt.Sprintf("%d", g.opts.MultiPV))
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
		g.processOneMove()
		if g.gameFinished() {
			break
		}
	}
}

func (g *GameRater) processOneMove() {
	g.extractCurrentBoard()
	g.writeCurrentFen()
	g.writeActualMovePlayed()
	g.sendPositionToEngine()
	g.startSearch()
	g.clearPvs()
	g.processAllEngineMessages()
	g.processEngineResults()
	g.updateGameState()
}

func (g *GameRater) extractCurrentBoard() {
	g.board = g.node.Board
}

func (g *GameRater) writeCurrentFen() {
	fen := g.board.Fen()
	g.write("FEN:    %s\n", fen)
}

func (g *GameRater) writeActualMovePlayed() {
	nextNode := g.node.Next
	if nextNode != nil {
		move := nextNode.Move
		actualSan := move.San(g.board)
		g.write("PLAYED: %s\n", actualSan)
	}
}

func (g *GameRater) sendPositionToEngine() {
	g.engine.SetPosition(g.board)
}

func (g *GameRater) startSearch() {
	var analysisTime = g.opts.TimePerMove
	g.infoChannel = g.engine.SearchTime(time.Second *
		time.Duration(analysisTime))
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
		// actually don't care about the best move; it is
		// duplicated in the MultiPV analysis. We also don't
		// need to indicate that we need to stop processing
		// messages from the engine; the engine will close
		// the channel.
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
	g.write("ANALYSIS:\n")
	for rank, pv := range g.allPvs {
		if pv == nil {
			continue
		}
		if len(pv.Moves) == 0 {
			continue
		}
		move := pv.Moves[0]
		san := move.San(g.board)
		score := pv.Score

		s := fmt.Sprintf("%d %7s %3d\n",
			rank, san, score)
		g.write(s)
	}
	g.write("\n")
}

func (g *GameRater) updateGameState() {
	g.node = g.node.Next
}

func (g *GameRater) gameFinished() bool {
	return g.node == nil
}

func (g *GameRater) clearPvs() {
	g.allPvs = make(map[int]*engine.Pv)
}

func (g *GameRater) write(format string, args ...interface{}) {
	g.outputFileHandle.WriteString(fmt.Sprintf(format, args...))
}
