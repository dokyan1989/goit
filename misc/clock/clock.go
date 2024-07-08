package clock

import "time"

type Clock interface {
	Now() time.Time
}

func Real() *realClock {
	return &realClock{}
}

type realClock struct{}

func (*realClock) Now() time.Time {
	return time.Now()
}

type fakeClock struct {
	now time.Time
}

func Fake(opts ...Option) *fakeClock {
	fk := &fakeClock{
		now: time.Now(),
	}
	for _, opt := range opts {
		opt(fk)
	}

	return fk
}

func (fk *fakeClock) Now() time.Time {
	return fk.now
}

type Option func(*fakeClock)

func WithNow(now time.Time) func(*fakeClock) {
	return func(fk *fakeClock) {
		fk.now = now
	}
}
