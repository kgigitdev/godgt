package godgt

// Message is really just a way to multiplex three different types of message
// onto a single channel.
type Message struct {
	BoardUpdate *BoardUpdate
	FieldUpdate *FieldUpdate
	TimeUpdate  *TimeUpdate
	InfoUpdate  *InfoUpdate
}

// Note, not implementing Stringer interface as you can't implement
// that on pointer receiver.
func (m *Message) ToString() string {
	if m.BoardUpdate != nil {
		return m.BoardUpdate.ToString()
	} else if m.FieldUpdate != nil {
		return m.FieldUpdate.ToString()
	} else if m.TimeUpdate != nil {
		return m.TimeUpdate.ToString()
	} else if m.InfoUpdate != nil {
		return m.InfoUpdate.ToString()
	} else {
		return ""
	}
}

func NewBoardUpdateMessage(boardUpdate *BoardUpdate) *Message {
	return &Message{
		BoardUpdate: boardUpdate,
	}
}

func NewFieldUpdateMessage(fieldUpdate *FieldUpdate) *Message {
	return &Message{
		FieldUpdate: fieldUpdate,
	}
}

func NewTimeUpdateMessage(timeUpdate *TimeUpdate) *Message {
	return &Message{
		TimeUpdate: timeUpdate,
	}
}

func NewInfoUpdateMessage(info string) *Message {
	infoUpdate := NewInfoUpdate(info)
	return &Message{
		InfoUpdate: infoUpdate,
	}
}
