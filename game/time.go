package game

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
}

func NewTime(fps int) *Time {
	return &Time{
		delta:       0,
		fixedDelta:  time.Second / time.Duration(fps),
		accumulator: 0,
		lastUpdate:  time.Now(),
		fixedSteps:  0,
	}
}

func (t *Time) start() {
	t.lastUpdate = time.Now()
}

func (t *Time) tick() {
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
