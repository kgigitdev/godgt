package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"text/tabwriter"
	"time"

	// "github.com/freeeve/uci"
	"github.com/malbrecht/chess/engine"
	"github.com/malbrecht/chess/engine/uci"
	"github.com/malbrecht/chess/pgn"

	flags "github.com/jessevdk/go-flags"
)

var opts struct {
	TimePerOpeningMove int `long:"opening-time" short:"t" description:"Time, in seconds, to allocate to each move analysis in the opening."`

	OpeningLength int `long:"opening-length" short:"l" description:"Length of the opening, in moves" default:"10"`

	TimePerMove int `long:"time-per-move" short:"T" description:"Time, in seconds, to allocate to each move analysis." default:"30"`

	DepthPerMove int `long:"depth-per-move" short:"d" description:"Depth to analyse each move (overrides time)"`

	PgnFile string `long:"pgn" short:"p" description:"PGN file to analyse." required:"true"`

	MultiPV int `long:"multipv" short:"m" description:"Number of alternative moves to analyse." default:"5" required:"true"`

	OutputFile string `long:"output" short:"o" description:"Output analysis file." default:"-"`

	Verbose bool `long:"verbose" short:"v" description:"Be more verbose."`
}

func main() {
	parser := flags.NewParser(&opts, flags.Default)
	_, err := parser.ParseArgs(os.Args)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	file := opts.PgnFile
	fh, err := os.Open(file)

	if err != nil {
		log.Fatal(err)
	}

	var oh *os.File
	if opts.OutputFile == "-" {
		oh = os.Stdout
	} else {
		oh, err = os.Create(opts.OutputFile)
		if err != nil {
			log.Fatal(err)
		}
	}

	pgntext, err := ioutil.ReadAll(fh)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	db := pgn.DB{}
	errors := db.Parse(string(pgntext))

	if len(errors) > 0 {
		for _, err := range errors {
			log.Println(err)
		}
		os.Exit(1)
	}

	game := db.Games[0]
	err = db.ParseMoves(game)

	if err != nil {
		log.Fatal(err)
	}

	var sfargs []string
	// logger := log.New(os.Stdout, "", log.LstdFlags)
	var logger *log.Logger
	e, err := uci.Run("/usr/games/stockfish", sfargs, logger)
	// eng, err := uci.NewEngine("/usr/games/stockfish")
	if err != nil {
		log.Fatal(err)
	}
	defer e.Quit()

	opt := e.Options()

	if opts.Verbose {
		w := tabwriter.NewWriter(os.Stdout, 1, 8, 0, ' ', 0)
		for k, v := range opt {
			fmt.Fprintln(w, k, "\t", v)
		}
		w.Flush()
	}

	mpv, ok := opt["MultiPV"]
	if !ok {
		log.Fatal("Can't set MultiPV mode")
	}
	mpv.Set(fmt.Sprintf("%d", opts.MultiPV))

	// The root node is the node after "zero" moves; it contains
	// the initial starting position, and a null move. We
	// therefore need to pull out the next move, so we can compare
	// what the engine would do in the current position with what
	// the player actually played (which is stored on the *next*
	// node)
	node := game.Root
	for {
		board := node.Board
		fen := board.Fen()
		oh.WriteString(fmt.Sprintf("FEN:    %s\n", fen))

		nextNode := node.Next
		if nextNode != nil {
			move := nextNode.Move
			actualSan := move.San(board)
			oh.WriteString(fmt.Sprintf("PLAYED: %s\n", actualSan))
		}

		e.SetPosition(board)

		var allPvs [5]*engine.Pv

		var analysisTime = opts.TimePerMove

		for info := range e.SearchTime(time.Second *
			time.Duration(analysisTime)) {
			if info.Err() != nil {
				log.Fatalf("%s", info.Err())
			}
			if _, ok := info.BestMove(); ok {
				// This indicates that the thinking
				// time is over.
				oh.WriteString("ANALYSIS:\n")
				for rank, pv := range allPvs {
					if pv == nil {
						continue
					}
					if len(pv.Moves) == 0 {
						continue
					}
					move := pv.Moves[0]
					san := move.San(board)
					score := pv.Score

					s := fmt.Sprintf("%d %7s %3d\n",
						rank, san, score)
					oh.WriteString(s)
				}
				oh.WriteString("\n")

			} else if pv := info.Pv(); pv != nil {
				allPvs[pv.Rank] = pv
				// log.Println("stats:", info.Stats())
			} else {
				// log.Println("stats:", info.Stats())
			}
		}
		node = node.Next
		if node == nil {
			break
		}
	}
}
