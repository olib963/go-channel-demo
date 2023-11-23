package actor

type Behaviour interface {
	isBehaviour()
}

type stop struct{}

func (s stop) isBehaviour() {}

func Stop() Behaviour {
	return stop{}
}

type same struct{}

func (s same) isBehaviour() {}

func Same() Behaviour {
	return same{}
}

type failed struct {
	err error
}

func (f failed) isBehaviour() {}

func Failed(err error) Behaviour {
	return failed{err}
}
