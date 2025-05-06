package domain

import (
	"fmt"
	"time"
)

type EventID int

const (
	EventUnknown                   EventID = 0
	EventCompetitorRegistered      EventID = 1
	EventStartTimeSet              EventID = 2
	EventCompetitorOnStartLine     EventID = 3
	EventCompetitorStarted         EventID = 4
	EventCompetitorOnFiringRange   EventID = 5
	EventTargetHit                 EventID = 6
	EventCompetitorLeftFiringRange EventID = 7
	EventCompetitorEnteredPenalty  EventID = 8
	EventCompetitorLeftPenalty     EventID = 9
	EventCompetitorEndedMainLap    EventID = 10
	EventCompetitorCanNotContinue  EventID = 11
)

type ScannerEvent interface {
	Scan(*Event) error
}

type Event struct {
	Time         time.Time
	ID           EventID
	CompetitorID int
	Comments     string
}

func (e *Event) Format() string {
	timestamp := e.Time.Format("15:04:05.000")
	switch e.ID {
	case EventCompetitorRegistered:
		return fmt.Sprintf("[%s] The competitor(%d) registered", timestamp, e.CompetitorID)
	case EventStartTimeSet:
		return fmt.Sprintf("[%s] The start time for the competitor(%d) was set by a draw to %s", timestamp, e.CompetitorID, e.Comments)
	case EventCompetitorOnStartLine:
		return fmt.Sprintf("[%s] The competitor(%d) is on the start line", timestamp, e.CompetitorID)
	case EventCompetitorStarted:
		return fmt.Sprintf("[%s] The competitor(%d) has started", timestamp, e.CompetitorID)
	case EventCompetitorOnFiringRange:
		return fmt.Sprintf("[%s] The competitor(%d) is on the firing range(%s)", timestamp, e.CompetitorID, e.Comments)
	case EventTargetHit:
		return fmt.Sprintf("[%s] The target(%s) has been hit by competitor(%d)", timestamp, e.Comments, e.CompetitorID)
	case EventCompetitorLeftFiringRange:
		return fmt.Sprintf("[%s] The competitor(%d) left the firing range", timestamp, e.CompetitorID)
	case EventCompetitorEnteredPenalty:
		return fmt.Sprintf("[%s] The competitor(%d) entered the penalty laps", timestamp, e.CompetitorID)
	case EventCompetitorLeftPenalty:
		return fmt.Sprintf("[%s] The competitor(%d) left the penalty laps", timestamp, e.CompetitorID)
	case EventCompetitorEndedMainLap:
		return fmt.Sprintf("[%s] The competitor(%d) ended the main lap", timestamp, e.CompetitorID)
	case EventCompetitorCanNotContinue:
		return fmt.Sprintf("[%s] The competitor(%d) can`t continue: %s", timestamp, e.CompetitorID, e.Comments)
	default:
		return fmt.Sprintf("[%s] Unknown event(%d) for competitor(%d)", timestamp, e.ID, e.CompetitorID)
	}
}
