package reporting

import (
	"fmt"
	"time"

	"github.com/Valery223/biathlon-test/internal/config"
	"github.com/Valery223/biathlon-test/internal/domain"
)

// ShotsPerFiring defines the number of shots fired per firing session in a biathlon.
const ShotsPerFiring = 5

// LapTime holds statistics for a single lap (main or penalty).
type LapStat struct {
	Duration     time.Duration
	AverageSpeed float64
}

// Report aggregates statistics for a single competitor's performance in the race.
type Report struct {
	CompetitorID        int
	Status              domain.Status
	TotalTime           time.Duration
	LapsStatistics      []LapStat
	PenaltyLapStatictic LapStat
	Shots               int
	PossibleShots       int
}

// CalculateReport generates a performance Report for a given competitor based on their race data and the configuration.
func CalculateReport(c domain.Competitor, cfg *config.Config) Report {
	r := Report{
		CompetitorID:   c.ID,
		LapsStatistics: make([]LapStat, 0, len(c.Laps)),
		Shots:          c.Shots,
	}

	r.Status = c.Status

	if len(c.Laps) == 0 || c.Laps[0].End.IsZero() {
		r.TotalTime = 0
	} else {
		r.TotalTime = c.Laps[len(c.Laps)-1].End.Sub(c.ScheduledStart)
	}

	r.PossibleShots = c.FiringCount * ShotsPerFiring

	currentLapStartTime := c.ScheduledStart
	for _, lap := range c.Laps {
		var lapStat LapStat
		lapStat.Duration = lap.End.Sub(currentLapStartTime)
		if lapStat.Duration > 0 {
			lapStat.AverageSpeed = float64(cfg.LapLength) / lapStat.Duration.Seconds()
		} else {
			lapStat.AverageSpeed = 0
		}
		r.LapsStatistics = append(r.LapsStatistics, lapStat)

		currentLapStartTime = lap.End
	}

	var tottalPenaltyLapTime time.Duration = 0
	for _, pLap := range c.PenaltyLaps {
		if pLap.Start.IsZero() || pLap.End.IsZero() {
			continue
		}
		tottalPenaltyLapTime += pLap.End.Sub(pLap.Start)
	}

	r.PenaltyLapStatictic.Duration = tottalPenaltyLapTime
	if len(c.PenaltyLaps) == 0 || tottalPenaltyLapTime == 0 {
		r.PenaltyLapStatictic.AverageSpeed = 0
	} else {
		r.PenaltyLapStatictic.AverageSpeed = float64(cfg.PenaltyLength) / tottalPenaltyLapTime.Seconds()
	}

	return r
}

// FormatDuration formats a time.Duration into a "HH:MM:SS.mmm" string.
func FormatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	ms := int(d.Milliseconds()) % 1000
	return fmt.Sprintf("%02d:%02d:%02d.%03d", h, m, s, ms)
}

// String provides a compact string representation of the Report
// conforming to the specified output format:
// total_time id [{time_lap, avg}, ...] {time_penalty_lap, avg_penalty} shots/PossibleShots
func (r Report) String() string {
	var result string
	switch r.Status {
	case domain.StatusNotFinished:
		result += "[NotFinished]"
	case domain.StatusNotStarted:
		result += "[NotStarted]"
	case domain.StatusFinished:
		result += FormatDuration(r.TotalTime)
	}
	result += " "

	result += fmt.Sprintf("%d ", r.CompetitorID)
	result += " "

	result += "["
	for i, lapStat := range r.LapsStatistics {
		if i > 0 {
			result += ", "
		}
		result += fmt.Sprintf("{%s %.3f}", FormatDuration(lapStat.Duration), lapStat.AverageSpeed)
	}
	result += "]"
	result += " "

	result += fmt.Sprintf("{%s, %.2f}", FormatDuration(r.PenaltyLapStatictic.Duration), r.PenaltyLapStatictic.AverageSpeed)
	result += " "

	result += fmt.Sprintf("%d/%d", r.Shots, r.PossibleShots)

	return result
}
