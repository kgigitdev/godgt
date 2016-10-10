package godgt

// The public API of DgtBoard.

import "io"

type DgtBoard struct {
	port           io.ReadWriteCloser
	bytesFromBoard []byte

	// A channel for reading messages from the board.
	MessagesFromBoard chan *Message

	// A channel for sending commands to the board.
	CommandsToBoard chan *Command
}

func (dgtboard *DgtBoard) WriteBytes(bytes []byte) (int, error) {
	return dgtboard.port.Write(bytes)
}

func (dgtboard *DgtBoard) WriteCommand(command byte) (int, error) {
	bytes := []byte{command}
	return dgtboard.port.Write(bytes)
}

func (dgtboard *DgtBoard) Close() {
	if dgtboard.port != nil {
		dgtboard.port.Close()
	}
}

// NewBoard() ...
func NewDgtBoard(portName string) *DgtBoard {
	port, err := CreatePort(portName)

	if err != nil {
		panic(err)
	}

	// What values here are sane?
	messagesFromBoard := make(chan *Message, 1024)
	commandsToBoard := make(chan *Command, 1024)

	return &DgtBoard{
		port:              port,
		MessagesFromBoard: messagesFromBoard,
		CommandsToBoard:   commandsToBoard,
	}
}
