package godgt

import "log"

func (dgtboard *DgtBoard) ReadLoop() {
	for {
		// log.Println("About to read bytes")
		dgtboard.readBytes()
		// log.Println("About to parse bytes")
		message, err := dgtboard.parseBytes()

		// FIXME: here we are discarding ALL errors,
		// but printing them all out is noisy; lots of
		// the errors are benign (e.g. "not enough bytes in buffer").
		// Maybe return a nil error in those cases and also a nil
		// message?
		if err == nil {
			dgtboard.MessagesFromBoard <- message
		}
	}
}

func (dgtboard *DgtBoard) readBytes() {
	buf := make([]byte, 1024)
	n, err := dgtboard.port.Read(buf)
	if err != nil {
		log.Println("error reading bytes.")
	}
	if n > 0 {
		// fmt.Printf("Read %d bytes\n", n)
		for i := 0; i < n; i++ {
			b := buf[i]
			dgtboard.bytesFromBoard = append(dgtboard.bytesFromBoard, b)
		}
	}
}
