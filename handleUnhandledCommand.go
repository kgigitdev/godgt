package godgt

import "errors"

var ERR_UNHANDLED = errors.New("Unhandled")

func (dgtboard *DgtBoard) defaultUnhandler(arguments []byte) (*Message, error) {
	return nil, ERR_UNHANDLED
}
