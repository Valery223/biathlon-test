package domain

import "time"

// Lap represents a single main lap in the race.
type Lap struct {
	Start     time.Time
	End       time.Time
	TargetHit int
	// FiringStart time.Time // Not used
	// FiringEnd   time.Time // Not used
}

// PenaltyLap represents a single penalty lap.
type PenaltyLap struct {
	Start time.Time
	End   time.Time
}
