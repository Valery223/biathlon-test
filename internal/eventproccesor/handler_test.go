package eventproccesor

import (
	"testing"
	"time"

	"github.com/Valery223/biathlon-test/internal/domain"
)

func TestHandleEvent(t *testing.T) {

	baseTime := time.Date(2025, 6, 6, 10, 0, 0, 0, time.UTC)
	lapsCount := 5

	// Create example competitors
	newTestCompetitor := func(id int) *domain.Competitor {
		c := domain.NewCompetitor(id)
		c.ScheduledStart = baseTime
		c.Laps = make([]domain.Lap, lapsCount)
		return c
	}
	t.Run("EventCompetitorRegistered", func(t *testing.T) {
		competitors := make(map[int]*domain.Competitor)
		event := &domain.Event{Time: baseTime, ID: domain.EventCompetitorRegistered, CompetitorID: 1}

		err := HandleEvent(event, competitors, lapsCount)
		if err != nil {
			t.Fatalf("HandleEvent failed: %v", err)
		}

		if _, ok := competitors[1]; !ok {
			t.Errorf("Competitor 1 not found in map after registration")
		}
		if competitors[1].ID != 1 {
			t.Errorf("Competitor ID mismatch: got %d, want 1", competitors[1].ID)
		}
		if competitors[1].Status != domain.StatusNotStarted {
			t.Errorf("Competitor status mismatch: got %v, want %v", competitors[1].Status, domain.StatusNotStarted)
		}
	})

	t.Run("EventStartTimeSet", func(t *testing.T) {
		competitors := make(map[int]*domain.Competitor)
		competitors[1] = newTestCompetitor(1)

		scheduledStartTimeStr := "10:05:00.000"
		expectedScheduledTime, _ := time.Parse("15:04:05.000", scheduledStartTimeStr)

		event := &domain.Event{Time: baseTime, ID: domain.EventStartTimeSet, CompetitorID: 1, Comments: scheduledStartTimeStr}
		err := HandleEvent(event, competitors, lapsCount)
		if err != nil {
			t.Fatalf("HandleEvent failed: %v", err)
		}

		got := competitors[1].ScheduledStart.Format("15:04:05.000")
		want := expectedScheduledTime.Format("15:04:05.000")
		if got != want {
			t.Errorf("ScheduledStart mismatch: got %s, want %s", got, want)
		}
	})
}
