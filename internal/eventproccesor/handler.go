package eventproccesor

import (
	"fmt"
	"time"

	"github.com/Valery223/biathlon-test/internal/domain"
)

// HandleEvent processes a single race event and updates the state of the relevant competitor.
// It takes the event, a map of all competitors, and the race configuration as input.
// It returns an error if the event is invalid or cannot be processed.
func HandleEvent(e *domain.Event, competitors map[int]*domain.Competitor, lapsCount int) error {
	competitorsID := e.CompetitorID
	competitor, ok := competitors[competitorsID]

	// For most events, the competitor must already exist(if event is not CompetitorRegistered)
	if e.ID != domain.EventCompetitorRegistered && !ok {
		return fmt.Errorf("competitor %d not found", competitorsID)
	}
	switch e.ID {
	case domain.EventCompetitorRegistered:
		competitor := domain.NewCompetitor(competitorsID)
		competitor.Laps = append(competitor.Laps, domain.Lap{})
		competitors[competitorsID] = competitor

	case domain.EventStartTimeSet:
		var err error
		competitor.ScheduledStart, err = time.Parse("15:04:05.000", e.Comments)
		if err != nil {
			return fmt.Errorf("error parsing time: %w", err)
		}
		competitor.Laps[competitor.CurrentLap].Start = competitor.ScheduledStart

	case domain.EventCompetitorOnStartLine:
		// Nothing to do
		// No specific state change implemented for this event yet
	case domain.EventCompetitorStarted:
		competitor.ActualStart = e.Time
	case domain.EventCompetitorOnFiringRange:
		competitor.FiringCount++
	case domain.EventTargetHit:
		competitor.Shots++
	case domain.EventCompetitorLeftFiringRange:
		// Nothing to do
		// No specific state change implemented for this event yet
	case domain.EventCompetitorEnteredPenalty:
		competitor.PenaltyLaps = append(competitor.PenaltyLaps, domain.PenaltyLap{
			Start: e.Time,
		})
	case domain.EventCompetitorLeftPenalty:
		competitor.PenaltyLaps[len(competitor.PenaltyLaps)-1].End = e.Time
	case domain.EventCompetitorEndedMainLap:
		competitor.Laps[competitor.CurrentLap].End = e.Time

		if competitor.CurrentLap+1 < lapsCount {
			competitor.CurrentLap++
			competitor.Laps = append(competitor.Laps, domain.Lap{})
			competitor.Laps[competitor.CurrentLap].Start = e.Time
		} else {
			competitor.Status = domain.StatusFinished
		}
	case domain.EventCompetitorCanNotContinue:
		competitor.Status = domain.StatusNotFinished
	default:
		return fmt.Errorf("unknown event %d", e.ID)
	}

	return nil
}
