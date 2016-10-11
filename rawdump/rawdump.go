package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/kgigitdev/godgt"
)

var opts struct {
	Pngs bool   `long:"pngs" description:"Write PNG images of board updates"`
	Port string `short:"p" long:"port" description:"Serial port" default:"/dev/ttyUSB0" env:"DGT_PORT"`

	Size int `short:"s" long:"size" description:"Image size" default:"128"`

	Filename string `short:"f" long:"filename" default:"boardupdate"`
}

func main() {

	_, err := flags.ParseArgs(&opts, os.Args)

	if err != nil {
		os.Exit(1)
	}

	dgtboard := godgt.NewDgtBoard(opts.Port)

	dgtboard.WriteCommand(godgt.DGT_SEND_RESET)
	dgtboard.WriteCommand(godgt.DGT_SEND_BRD)
	dgtboard.WriteCommand(godgt.DGT_SEND_UPDATE_BRD)

	// Ask the board for a complete dump every 10 seconds.
	go func() {
		t := time.NewTicker(time.Second * 10)
		for {
			dgtboard.WriteCommand(godgt.DGT_SEND_BRD)
			<-t.C
		}
	}()

	go dgtboard.ReadLoop()

	var messageCount int

	for {
		select {
		case message := <-dgtboard.MessagesFromBoard:
			messageCount += 1
			writeMessage(message)
			if opts.Pngs && message.BoardUpdate != nil {
				filename := fmt.Sprintf("%s-%04d.png",
					opts.Filename, messageCount)
				fen := message.BoardUpdate.ToString()
				godgt.WritePng(fen, opts.Size, filename)
			}
			if message.FieldUpdate != nil {
				dgtboard.WriteCommand(godgt.DGT_SEND_BRD)
			}
		}
	}
}

func writeMessage(m *godgt.Message) {
	if m.BoardUpdate != nil {
		log.Print("BOARD: ", m.ToString())
		rows := godgt.SimpleBoardFromFen(m.ToString())
		for _, row := range rows {
			log.Print(row)
		}
	} else if m.FieldUpdate != nil {
		log.Print("FIELD: ", m.ToString())
	} else {
		log.Print("OTHER: ", m.ToString())
	}
}
