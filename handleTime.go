package godgt

import (
	"errors"
	"log"
)

var ERR_CLOCK_ACK = errors.New("Clock ACK (Unhandled)")
var ERR_CLOCK_NOT_RUNNNING = errors.New("Clock Not Running")
var ERR_CLOCK_NOT_CONNECTED = errors.New("Clock Not Connected")

func (dgtboard *DgtBoard) handleTime(arguments []byte) (*Message, error) {
	log.Println("DGT_BWTIME")

	// byteN = "byte N of the whole message", with the first byte
	// after the beader being the 0th argument. There are 10 bytes
	// in the whole message, so the first byte of the arguments is
	// byte3 and the last byte is byte9. Note that in the spec,
	// the Clock Ack message is described twice:
	//
	// "If the (4th byte & 0x0f) equals 0x0a, or if the (7th byte
	// & 0x0f) equals 0x0a, then the message is a Clock Ack
	// message"
	//
	// "(If (byte 3 & 0x0f) is 0x0a, then the msg is a Clock Ack
	// message instead)"
	//
	// Despite the apparent inconsistency, this almost certainly
	// refers to "4th byte of the message, which is byte 3
	// (counting from 0)".
	byte3 := arguments[0]

	if byte3&0x0f == 0x0a {
		return nil, ERR_CLOCK_ACK
	}

	rightPlayerHours := int(byte3 & 0x04)                   // 0b 0000 1111
	rightPlayerFlagFallenAndBlocked := byte3&0x10 == 0x10   // 0b 0001 0000
	rightPlayerTimePerMoveIndicator := byte3&0x20 == 0x20   // 0b 0010 0000
	rightPlayerFlagFallenAndIndicated := byte3&0x40 == 0x40 // 0b 0100 0000

	byte4 := arguments[1]
	rightPlayerMinutes := int(byte4)

	byte5 := arguments[2]
	rightPlayerSeconds := int(byte5)

	byte6 := arguments[3]
	leftPlayerHours := int(byte6 & 0x04)                   // 0b 0000 1111
	leftPlayerFlagFallenAndBlocked := byte6&0x10 == 0x10   // 0b 0001 0000
	leftPlayerTimePerMoveIndicator := byte6&0x20 == 0x20   // 0b 0010 0000
	leftPlayerFlagFallenAndIndicated := byte6&0x40 == 0x40 // 0b 0100 0000

	byte7 := arguments[4]
	leftPlayerMinutes := int(byte7)

	byte8 := arguments[5]
	leftPlayerSeconds := int(byte8)

	byte9 := arguments[6]

	log.Printf("Byte A: 0b%08b\n", byte9)

	clockRunning := byte9&0x01 == 0x01             // 0b 0000 0001
	clockRockerLeftDepressed := byte9&0x02 == 0x02 // 0b 0000 0010
	clockRockerRightDepressed := !clockRockerLeftDepressed
	batteryLow := byte9&0x04 == 0x04        // 0b 0000 0100
	rightPlayersTurn := byte9&0x08 == 0x08  // 0b 0000 1000
	leftPlayersTurn := byte9&0x10 == 0x10   // 0b 0001 0000
	clockNotConnected := byte9&0x20 == 0x20 // 0b 0010 0000

	if clockNotConnected {
		return nil, ERR_CLOCK_NOT_CONNECTED
	}

	if !clockRunning {
		return nil, ERR_CLOCK_NOT_RUNNNING
	}

	log.Printf("%02d:%02d:%02d - %02d:%02d:%02d\n", leftPlayerHours, leftPlayerMinutes, leftPlayerSeconds, rightPlayerHours, rightPlayerMinutes, rightPlayerSeconds)

	log.Printf("%t %t %t %t %t %t %t %t %t %t %t\n", leftPlayerFlagFallenAndBlocked,
		leftPlayerTimePerMoveIndicator,
		leftPlayerFlagFallenAndIndicated,
		rightPlayerFlagFallenAndBlocked,
		rightPlayerTimePerMoveIndicator,
		rightPlayerFlagFallenAndIndicated,
		clockRockerLeftDepressed, clockRockerRightDepressed,
		batteryLow, rightPlayersTurn, leftPlayersTurn)

	// FIXME
	timeUpdate := NewTimeUpdate()
	timeUpdateMessage := NewTimeUpdateMessage(timeUpdate)
	return timeUpdateMessage, nil
}
