package godgt

import (
	"errors"
	"fmt"
	"log"
)

var ERR_NOT_ENOUGH_DATA = errors.New("Not enough data")
var ERR_NONE_COMMAND = errors.New("NONE Command")
var ERR_PARSE_FAILED = errors.New("Failed to parse bytes")

func (dgtboard *DgtBoard) parseBytes() (*Message, error) {
	if len(dgtboard.bytesFromBoard) < 3 {
		// Since a well-formed header is always 3 bytes,
		// if we haven't read at least three bytes, there's
		// no point in even trying.
		return nil, ERR_NOT_ENOUGH_DATA
	}
	b0 := dgtboard.bytesFromBoard[0]
	b1 := dgtboard.bytesFromBoard[1]
	b2 := dgtboard.bytesFromBoard[2]
	// log.Printf("B0 : 0x%02x 0b%08b\n", b0, b0)
	// log.Printf("B1 : 0x%02x 0b%08b\n", b1, b1)
	// log.Printf("B2 : 0x%02x 0b%08b\n", b2, b2)

	if (b0 | MESSAGE_BIT) == 0 {
		// FIXME: Handle this properly (maybe consume characters until
		// the next message header?)
		err := errors.New(fmt.Sprintf("Received corrupt message header: %0b08b\n", b0))
		return nil, err
	}

	// Mask out the MSB to find out what the message is. mx = "xth
	// masked byte".
	m0 := b0 & MESSAGE_MASK
	// log.Printf("M0 : 0x%02x 0b%08b\n", m0, m0)
	// This should always be the case anyway, but do this just
	// case.
	m1 := b1 & MESSAGE_MASK
	m2 := b2 & MESSAGE_MASK

	// Combine the two bytes to determine the message length. Note
	// that each byte only contains 7 bits of the message length,
	// so the message length can be a maximum of 14 bits. We therefore
	// also need to multiply the first byte
	length := (m1 << 7) + m2
	// log.Printf("Decoded message length: %d bytes.\n", length)

	// Pull off the bytes from the front of the bytes slice;
	// subsequent messages (if any) will be read the next time
	// this method is invoked.
	if len(dgtboard.bytesFromBoard) < int(length) {
		// We haven't read the complete command yet.
		// log.Printf("Detected incomplete command in the buffer.")
		return nil, ERR_NOT_ENOUGH_DATA
	}

	completeCommand := dgtboard.bytesFromBoard[0:length]

	// Turns out we have to have this earlier. A DGT_NONE message
	// has a length of 0, with no bytes beyond the 3 header bytes,
	// so we cannot pull out the arguments with "[3:]", since index
	// 3 is already out of bounds. FIXME: When we start getting
	// these, we tend to get lots of them, apparently without end.
	// Maybe reset the board automatically if we receive lots of these?
	if m0 == DGT_NONE {
		log.Println("DGT_NONE")
		return nil, ERR_NONE_COMMAND
	}

	// The length argument always includes the 3 header bytes;
	// we can therefore remove the first three bytes and consider
	// the remainder to be the arguments to the command.
	arguments := completeCommand[3:]
	dgtboard.bytesFromBoard = dgtboard.bytesFromBoard[length:]

	// Now handle the command.
	switch m0 {
	case DGT_NONE:
		return dgtboard.defaultUnhandler(arguments)
	case DGT_BOARD_DUMP:
		return dgtboard.handleBoardDump(arguments)
	case DGT_BWTIME:
		return dgtboard.handleTime(arguments)
	case DGT_FIELD_UPDATE:
		return dgtboard.handleFieldUpdate(arguments)
	case DGT_EE_MOVES:
		return dgtboard.defaultUnhandler(arguments)
	case DGT_BUSADRES:
		return dgtboard.defaultUnhandler(arguments)
	case DGT_SERIALNR:
		return dgtboard.defaultUnhandler(arguments)
	case DGT_TRADEMARK:
		return dgtboard.handleTrademarkMessage(arguments)
	case DGT_VERSION:
		return dgtboard.handleVersionMessage(arguments)
	default:
		return nil, ERR_PARSE_FAILED
	}
}
