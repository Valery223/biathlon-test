package eventproccesor

import (
	"fmt"
	"time"

	"github.com/Valery223/biathlon-test/internal/domain"
)

func HandleEvent(e *domain.Event, competitors map[int]*domain.Competitor, lapsCount int) error {
	competitorsID := e.CompetitorID
	competitor, ok := competitors[competitorsID]
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
		//Nothing to do
	case domain.EventCompetitorStarted:
		competitor.ActualStart = e.Time
	case domain.EventCompetitorOnFiringRange:
		competitor.FiringCount++
		// lap := competitor.Laps[competitor.CurrentLap]
		// lap.FiringStart = e.Time
	case domain.EventTargetHit:
		competitor.Shots++
	case domain.EventCompetitorLeftFiringRange:
		// lap := competitor.Laps[competitor.CurrentLap]
		// lap.FiringEnd = e.Time
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
