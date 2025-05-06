package reporting

import (
	"fmt"
	"time"

	"github.com/Valery223/biathlon-test/internal/config"
	"github.com/Valery223/biathlon-test/internal/domain"
)

const ShotsPerFiring = 5

type LapStat struct {
	Duraction    time.Duration
	AverageSpeed float64
}

type Report struct {
	CompetitorID        int
	TotalTime           time.Duration
	LapsStatistics      []LapStat
	PenaltyLapStatictic LapStat
	Shots               int
	PossibleShots       int
}

func CalculateReport(c domain.Competitor, cfg *config.Config) Report {
	r := Report{
		CompetitorID:   c.ID,
		LapsStatistics: make([]LapStat, 0, len(c.Laps)),
		Shots:          c.Shots,
	}

	if len(c.Laps) == 0 || c.Laps[0].End.IsZero() {
		r.TotalTime = 0
	} else {
		r.TotalTime = c.Laps[len(c.Laps)-1].End.Sub(c.ScheduledStart)
	}

	r.PossibleShots = len(c.Laps) * ShotsPerFiring
	fmt.Println("--------,", len(c.Laps), r.PossibleShots)

	currentLapStartTime := c.ScheduledStart
	for _, lap := range c.Laps {
		var lapStat LapStat
		lapStat.Duraction = lap.End.Sub(currentLapStartTime)
		if lapStat.Duraction > 0 {
			lapStat.AverageSpeed = float64(cfg.LapLength) / lapStat.Duraction.Seconds()
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

	r.PenaltyLapStatictic.Duraction = tottalPenaltyLapTime
	if len(c.PenaltyLaps) == 0 || tottalPenaltyLapTime == 0 {
		r.PenaltyLapStatictic.AverageSpeed = 0
	} else {
		r.PenaltyLapStatictic.AverageSpeed = float64(cfg.PenaltyLength) / tottalPenaltyLapTime.Seconds()
	}

	return r
}

func FormatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	ms := int(d.Milliseconds()) % 1000
	return fmt.Sprintf("%02d:%02d:%02d.%03d", h, m, s, ms)
}

func (r Report) String() string {
	var result string
	result += fmt.Sprintf("Competitor ID: %d ", r.CompetitorID)
	result += fmt.Sprintf("Total Time: %s ", FormatDuration(r.TotalTime))

	result += "Lap Statistics:\n"
	result += "["
	for i, lap := range r.LapsStatistics {
		if i >= 1 {
			result += ", "
		}
		result += "{"

		result += fmt.Sprintf("Lap %d: Duration: %s, Average Speed: %.2f", i+1, FormatDuration(lap.Duraction), lap.AverageSpeed)
		result += "}"
	}
	result += "]\n"

	result += "{"
	result += fmt.Sprintf("Penalty Lap Time: %s ", FormatDuration(r.PenaltyLapStatictic.Duraction))
	result += fmt.Sprintf("Penalty Lap Average Speed: %.2f", r.PenaltyLapStatictic.AverageSpeed)
	result += "}\n"

	result += fmt.Sprintf(" Shots/Possible: %d/%d", r.Shots, r.PossibleShots)
	return result
}
