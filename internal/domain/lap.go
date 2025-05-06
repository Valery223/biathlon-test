package domain

import "time"

type Lap struct {
	Start     time.Time
	End       time.Time
	TargetHit int
	// FiringStart time.Time
	// FiringEnd   time.Time
}

type PenaltyLap struct {
	Start time.Time
	End   time.Time
}
