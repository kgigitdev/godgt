package main

// Opts contains all the command line options for the utility.
type Opts struct {
	Engine string `long:"engine" short:"e" description:"UCI Engine to run" default:"stockfish"`

	OpeningTimePerMove int `long:"opening-time" short:"t" description:"Time, in seconds, to allocate to each move analysis in the opening."`

	OpeningDepthPerMove int `long:"opening-depth" short:"d" description:"Time, in seconds, to allocate to each move analysis in the opening."`

	OpeningLength int `long:"opening-length" short:"l" description:"Length of the opening, in moves" default:"10"`

	TimePerMove int `long:"time-per-move" short:"T" description:"Time, in seconds, to allocate to each move analysis." default:"30"`

	DepthPerMove int `long:"depth-per-move" short:"D" description:"Depth to analyse each move (overrides time)"`

	PgnFile string `long:"pgn" short:"p" description:"PGN file to analyse." required:"true"`

	MultiPV int `long:"multipv" short:"m" description:"Number of alternative moves to analyse." default:"5" required:"true"`

	OutputFile string `long:"output" short:"o" description:"Output analysis file." default:"-"`

	Verbose bool `long:"verbose" short:"v" description:"Be more verbose."`
}
