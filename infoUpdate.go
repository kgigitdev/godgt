package godgt

type InfoUpdate struct {
	info string
}

func NewInfoUpdate(info string) *InfoUpdate {
	return &InfoUpdate{
		info: info,
	}
}

func (iu *InfoUpdate) ToString() string {
	return iu.info
}
