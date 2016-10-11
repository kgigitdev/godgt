package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kgigitdev/godgt"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s /path/to/usbdevice\n", os.Args[0])
		fmt.Printf("e.g.:  %s /dev/ttyUSB0\n", os.Args[0])
		os.Exit(1)
	}

	portName := os.Args[1]
	dgtboard := godgt.NewDgtBoard(portName)

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

	for {
		select {
		case message := <-dgtboard.MessagesFromBoard:
			writeMessage(message)
			if message.FieldUpdate != nil {
				dgtboard.WriteCommand(godgt.DGT_SEND_BRD)
			}
		}
	}
}

func writeMessage(m *godgt.Message) {
	if m.BoardUpdate != nil {
		log.Print("BOARD: ", m.ToString())
		rows := godgt.BoardFromFen(m.ToString())
		for _, row := range rows {
			log.Print(row)
		}
	} else if m.FieldUpdate != nil {
		log.Print("FIELD: ", m.ToString())
	} else {
		log.Print("OTHER: ", m.ToString())
	}
}
