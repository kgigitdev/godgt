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
)

func main() {

	file := os.Args[1]
	fh, err := os.Open(file)

	if err != nil {
		log.Fatal(err)
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

	var args []string
	// logger := log.New(os.Stdout, "", log.LstdFlags)
	var logger *log.Logger
	e, err := uci.Run("/usr/games/stockfish", args, logger)
	// eng, err := uci.NewEngine("/usr/games/stockfish")
	if err != nil {
		log.Fatal(err)
	}
	defer e.Quit()

	opt := e.Options()
	w := tabwriter.NewWriter(os.Stdout, 1, 8, 0, ' ', 0)
	for k, v := range opt {
		fmt.Fprintln(w, k, "\t", v)
	}
	w.Flush()

	mpv, ok := opt["MultiPV"]
	if !ok {
		log.Fatal("Can't set MultiPV mode")
	}
	mpv.Set("5")

	// set some engine options
	/*
		eng.SetOptions(uci.Options{
			Hash:    1024, // Hash size in MB
			Ponder:  false,
			OwnBook: true,
			MultiPV: 5,
		})
	*/

	node := game.Root

	for {
		board := node.Board
		fen := board.Fen()
		log.Println(fen)
		e.SetPosition(board)

		// set some result filter options
		// resultOpts := uci.HighestDepthOnly | uci.IncludeUpperbounds | uci.IncludeLowerbounds
		// results, err := eng.GoDepth(15, resultOpts)
		var allPvs [5]*engine.Pv
		for info := range e.SearchTime(time.Second * 30) {
			if info.Err() != nil {
				log.Fatalf("%s", info.Err())
			}
			if _, ok := info.BestMove(); ok {
				// This indicates that the thinking
				// time is over.
				for rank, pv := range allPvs {
					move := pv.Moves[0]
					san := move.San(board)
					score := pv.Score
					log.Printf("%d %7s %3d\n",
						rank, san, score)
				}
				// log.Println("bestmove:", m.San(board))
			} else if pv := info.Pv(); pv != nil {
				allPvs[pv.Rank] = pv
				// log.Println("stats:", info.Stats())
			} else {
				// log.Println("stats:", info.Stats())
			}
		}

		/*
			if err != nil {
				log.Fatal(err)
			}
		*/

		// print it (String() goes to pretty JSON for now)

		if node.Next == nil {
			break
		}
		node = node.Next
	}

	os.Exit(1)

	// set the starting position
	// eng.SetFEN("rnb4r/ppp1k1pp/3bp3/1N3p2/1P2n3/P3BN2/2P1PPPP/R3KB1R b KQ - 4 11")

}
