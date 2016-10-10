package main

import (
	"fmt"
	"os"

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

	// FIXME: Make this properly channel based.
	// FIXME: Make this part of the dgtboard class. Users should
	// not need to know that they need to send these.
	dgtboard.WriteCommand(godgt.DGT_SEND_RESET)
	dgtboard.WriteCommand(godgt.DGT_SEND_BRD)
	dgtboard.WriteCommand(godgt.DGT_SEND_UPDATE_BRD)

	mp := godgt.NewMessageProcessor()

	go dgtboard.ReadLoop()

	for {
		select {
		case message := <-dgtboard.MessagesFromBoard:
			mp.ProcessMessage(message)
		}
	}
}

// func handleMessage(dgtboard *godgt.DgtBoard, m *godgt.Message) {
// 	log.Print("Received message from board: ", m.ToString())
// 	if m.BoardUpdate != nil {
// 		go RunEngine(m.ToString())
// 	} else if m.FieldUpdate != nil {
// 		if m.FieldUpdate.Piece() != "." {
// 			dgtboard.WriteCommand(godgt.DGT_SEND_BRD)
// 		}
// 	}
// }
