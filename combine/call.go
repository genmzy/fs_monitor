package combine

type Call struct {
	LegTuple [2]string
	CallId   string
}

func (c *Call) With(another *Call) (b *Bridge) {
	// b := &Bridge{Calls: make([]*Call, 0)}
	// b.AtInstance = ""
	// b.AtEpoch = 0
	// b.Calls = append(b.Calls, another)
	return
}
