package game

import "time"

const (
	MaxAccumulator = time.Second / 4
	MaxDeltaTime   = time.Millisecond * 50
)

type Time struct {
	last        time.Time
	delta       time.Duration
	accumulator time.Duration
	fixedDelta  time.Duration
	fixedSteps  int
}

func NewTime(fps int) *Time {
	return &Time{
		fixedDelta:  time.Second / time.Duration(fps),
		delta:       0,
		accumulator: 0,
		fixedSteps:  0,
	}
}

func (t *Time) Tick() {
	now := time.Now()

	if t.last.IsZero() {
		t.last = now
		return
	}

	t.delta = now.Sub(t.last)
	t.last = now

	if t.delta <= 0 {
		return
	}

	if t.delta > MaxDeltaTime {
		t.delta = MaxDeltaTime
	}

	t.accumulator = min(t.accumulator+t.delta, MaxAccumulator)
	t.fixedSteps = 0

	for t.accumulator >= time.Duration(t.fixedDelta) {
		t.accumulator -= time.Duration(t.fixedDelta)
		t.fixedSteps++
	}
}

func (t *Time) Delta() float64 {
	return t.delta.Seconds()
}

func (t *Time) FixedDelta() float64 {
	return t.fixedDelta.Seconds()
}
