package combine

type Bridge struct {
	CallTuple  [2]string
	AtInstance string
	AtEpoch    int
}

// func (b *Bridge) Merge(others ...*Bridge) *Link {
// 	l := &Link{
// 		Bridges: make([]*Bridge, 0),
// 	}
// 	l.Bridges = append(l.Bridges, b)
// 	l.Bridges = append(l.Bridges, others...)
// 	return l
// }
