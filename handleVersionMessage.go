package godgt

import "fmt"

func (dgtboard *DgtBoard) handleVersionMessage(arguments []byte) (*Message, error) {
	major := arguments[0]
	minor := arguments[1]
	info := fmt.Sprintf("DGT_VERSION: %d.%02d\n", int(major), int(minor))
	return NewInfoUpdateMessage(info), nil
}
