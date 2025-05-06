package domain

import "time"

type Competitor struct {
	ID             int
	Status         Status
	ScheduledStart time.Time
	ActualStart    time.Time
	Laps           []Lap
	PenaltyLaps    []PenaltyLap
	Shots          int
	FiringCount    int // Number of times the competitor was on the firing line
	CurrentLap     int
}

type Status int

const (
	StatusFinished Status = iota
	StatusNotStarted
	StatusNotFinished
)

func NewCompetitor(id int) *Competitor {
	return &Competitor{
		ID:          id,
		Status:      StatusNotStarted,
		Laps:        make([]Lap, 0),
		PenaltyLaps: make([]PenaltyLap, 0),
		Shots:       0,
	}
}
