package godgt

// TimeUpdate encapsulates a raw time update from a DGT board.
type TimeUpdate struct {
}

func NewTimeUpdate() *TimeUpdate {
	return &TimeUpdate{}
}

func (tu *TimeUpdate) ToString() string {
	return "Time Update"
}
