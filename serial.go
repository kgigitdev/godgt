package godgt

import (
	"io"

	"github.com/jacobsa/go-serial/serial"
)

func CreatePort(portName string) (io.ReadWriteCloser, error) {
	options := serial.OpenOptions{
		PortName:        portName,
		BaudRate:        9600,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 1,
	}

	return serial.Open(options)
}
