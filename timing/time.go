package timing

import "time"

const (
	MaxDeltaTime       = time.Second / 4
	MaxAccumulatedTime = time.Second * 2
)

type Time struct {
	delta       time.Duration
	fixedDelta  time.Duration
	accumulator time.Duration
	lastUpdate  time.Time
	fixedSteps  int
	frame       uint64
}

func NewTime(fps int) *Time {
	return &Time{
		delta:       0,
		fixedDelta:  time.Second / time.Duration(fps),
		accumulator: 0,
		lastUpdate:  time.Now(),
		fixedSteps:  0,
		frame:       0,
	}
}

func (t *Time) Start() {
	t.lastUpdate = time.Now()
	t.frame = 0
}

func (t *Time) Tick() {
	t.frame++

	now := time.Now()

	t.delta = now.Sub(t.lastUpdate)
	t.lastUpdate = now

	if t.delta > MaxDeltaTime {
		t.delta = MaxDeltaTime
	}

	t.accumulator += t.delta
	t.fixedSteps = 0

	if t.accumulator > MaxAccumulatedTime {
		t.accumulator = MaxAccumulatedTime
	}

	for t.accumulator >= t.fixedDelta {
		t.accumulator -= t.fixedDelta
		t.fixedSteps++
	}
}

func (t *Time) Delta() float64 {
	return t.delta.Seconds()
}

func (t *Time) Delta32() float32 {
	return float32(t.delta.Seconds())
}

func (t *Time) FixedSteps() int {
	return t.fixedSteps
}

func (t *Time) FixedDelta() float64 {
	return t.fixedDelta.Seconds()
}

func (t *Time) FixedDelta32() float32 {
	return float32(t.fixedDelta.Seconds())
}

func (t *Time) Frame() uint64 {
	return t.frame
}
