package reporting

import (
	"testing"
	"time"

	"github.com/Valery223/biathlon-test/internal/config"
	"github.com/Valery223/biathlon-test/internal/domain"
)

func TestCalculateReport(t *testing.T) {
	cfg := &config.Config{
		LapLength:     3000,
		PenaltyLength: 100,
	}

	scheduledStart := time.Date(2025, 6, 6, 10, 0, 0, 0, time.UTC)

	testCases := []struct {
		name       string
		competitor *domain.Competitor
		wantReport Report
	}{
		{
			name: "finish with penalty",
			competitor: &domain.Competitor{
				ID:             1,
				Status:         domain.StatusFinished,
				ScheduledStart: scheduledStart,
				Laps: []domain.Lap{
					{End: scheduledStart.Add(10 * time.Minute)},
					{End: scheduledStart.Add(21 * time.Minute)},
				},
				PenaltyLaps: []domain.PenaltyLap{
					{
						Start: scheduledStart.Add(7*time.Minute - 45*time.Second),
						End:   scheduledStart.Add(7 * time.Minute)},
					{
						Start: scheduledStart.Add(19*time.Minute - 60*time.Second),
						End:   scheduledStart.Add(19 * time.Minute)},
				},
				Shots:       8,
				FiringCount: 2,
			},
			wantReport: Report{
				CompetitorID: 1,
				Status:       domain.StatusFinished,
				TotalTime:    21 * time.Minute,
				LapsStatistics: []LapStat{
					{Duration: 10 * time.Minute, AverageSpeed: float64(cfg.LapLength) / (10 * 60)},
					{Duration: 11 * time.Minute, AverageSpeed: float64(cfg.LapLength) / (11 * 60)},
				},
				PenaltyLapStatictic: LapStat{
					Duration:     1*time.Minute + 45*time.Second,
					AverageSpeed: float64(cfg.PenaltyLength) / (1*60 + 45),
				},
				Shots:         8,
				PossibleShots: 10,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			got := CalculateReport(*tc.competitor, cfg)

			if got.CompetitorID != tc.wantReport.CompetitorID {
				t.Errorf("CompetitorID: got %d, want %d", got.CompetitorID, tc.wantReport.CompetitorID)
			}
			if got.TotalTime != tc.wantReport.TotalTime {
				t.Errorf("TotalTime: got %v, want %v", got.TotalTime, tc.wantReport.TotalTime)
			}
			if len(got.LapsStatistics) != len(tc.wantReport.LapsStatistics) {
				t.Fatalf("LapsStatistics length: got %d, want %d. Got: %+v", len(got.LapsStatistics), len(tc.wantReport.LapsStatistics), got.LapsStatistics)
			}
			for i := range got.LapsStatistics {
				if got.LapsStatistics[i].Duration != tc.wantReport.LapsStatistics[i].Duration {
					t.Errorf("LapsStatistics[%d].Duration: got %v, want %v", i, got.LapsStatistics[i].Duration, tc.wantReport.LapsStatistics[i].Duration)
				}

				if !floatEquals(got.LapsStatistics[i].AverageSpeed, tc.wantReport.LapsStatistics[i].AverageSpeed, 0.001) {
					t.Errorf("LapsStatistics[%d].AverageSpeed: got %.3f, want %.3f", i, got.LapsStatistics[i].AverageSpeed, tc.wantReport.LapsStatistics[i].AverageSpeed)
				}
			}
			if got.PenaltyLapStatictic.Duration != tc.wantReport.PenaltyLapStatictic.Duration {
				t.Errorf("PenaltyLapStatistic.Duration: got %v, want %v", got.PenaltyLapStatictic.Duration, tc.wantReport.PenaltyLapStatictic.Duration)
			}
			if !floatEquals(got.PenaltyLapStatictic.AverageSpeed, tc.wantReport.PenaltyLapStatictic.AverageSpeed, 0.001) {
				t.Errorf("PenaltyLapStatistic.AverageSpeed: got %.3f, want %.3f", got.PenaltyLapStatictic.AverageSpeed, tc.wantReport.PenaltyLapStatictic.AverageSpeed)
			}
			if got.Shots != tc.wantReport.Shots {
				t.Errorf("Shots: got %d, want %d", got.Shots, tc.wantReport.Shots)
			}
			if got.PossibleShots != tc.wantReport.PossibleShots {
				t.Errorf("PossibleShots: got %d, want %d", got.PossibleShots, tc.wantReport.PossibleShots)
			}

		})
	}
}

func floatEquals(a, b, tolerance float64) bool {
	if (a-b) < tolerance && (b-a) < tolerance {
		return true
	}
	return false
}
