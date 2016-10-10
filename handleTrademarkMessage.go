package godgt

func (dgtboard *DgtBoard) handleTrademarkMessage(arguments []byte) (*Message, error) {
	info := string(arguments)
	return NewInfoUpdateMessage(info), nil
}
