package task

import (
	"fmt"
	"io"
	"log"
	"time"
)

type ScannerEvent interface {
	Scan(*Event) error
}

type Event struct {
	Time        time.Time
	EventID     int
	ExtraParams int
	Comments    string
}
type Status int

const (
	StatusFinished Status = iota
	StatusNotStarted
	StatusNotFinished
)

type Competitor struct {
	ID             int
	Status         Status
	ScheduledStart time.Time
	ActualStart    time.Time
	Laps           []Lap
	PenaltyLaps    []PenaltyLap
	Hits           int
	Shots          int
	CurrentLap     int
}

func NewCompetitor(id int, cfg *Config) *Competitor {
	return &Competitor{
		ID:          id,
		Status:      StatusNotStarted,
		Laps:        make([]Lap, cfg.Laps),
		PenaltyLaps: make([]PenaltyLap, 0),
		Hits:        0,
		Shots:       0,
	}
}

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

type Config struct {
	Laps          int       `json:"laps"`
	LapLength     int       `json:"lapLen"`
	PenaltyLength int       `json:"penaltyLen"`
	FiringLines   int       `json:"firingLines"`
	StartTime     time.Time `json:"start"`
	StartDelta    time.Time `json:"startDelta"`
}

type Task struct {
	cfg     *Config
	scanner ScannerEvent
}

func NewTask(cfg *Config, scanner ScannerEvent) *Task {
	return &Task{
		cfg:     cfg,
		scanner: scanner,
	}
}

func (t Task) Execute() error {
	mapCompetitors := make(map[int]*Competitor)
	for {
		event := &Event{}
		err := t.scanner.Scan(event)
		if err != nil {
			if err == io.EOF {
				log.Println("End of file reached")
				break
			}
			return fmt.Errorf("error scanning event: %w", err)

		}
		event.print()
		err = t.handleEvent(event, mapCompetitors)
		if err != nil {
			return fmt.Errorf("error handling event: %w", err)
		}

	}

	for _, competitor := range mapCompetitors {
		report := competitor.GetReport(t.cfg)
		fmt.Printf("Competitor %d: Total Time: %s, Laps: [", report.CompetitorID, formatDuration(report.TotalTime))
		for i, lap := range report.LapsStatistics {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Printf("{Duration: %s, AVG: %.2f} ", formatDuration(lap.Duraction), lap.AVG)
		}
		fmt.Println("]")

		if report.PenaltyLapStatictic.Duraction == 0 {
			fmt.Print("Penalty Laps: {0, -}")
		} else {
			fmt.Printf("Penalty Laps: {Duration: %s, AVG: %.2f}", formatDuration(report.PenaltyLapStatictic.Duraction), report.PenaltyLapStatictic.AVG)
		}
		fmt.Printf(" Hits/Shots: %d/%d\n", report.Hits, report.Shots)
	}

	fmt.Println("End of task")
	return nil
}

type LapTime struct {
	Duraction time.Duration
	AVG       float64
}

type Report struct {
	CompetitorID        int
	TotalTime           time.Duration
	LapsStatistics      []LapTime
	PenaltyLapStatictic LapTime
	Hits                int
	Shots               int
}

func (c *Competitor) GetReport(cfg *Config) Report {
	r := Report{}
	r.CompetitorID = c.ID
	r.TotalTime = c.Laps[len(c.Laps)-1].End.Sub(c.ScheduledStart)
	r.LapsStatistics = make([]LapTime, 0)
	startLap := c.ScheduledStart
	for _, lap := range c.Laps {
		var lapTime LapTime
		lapTime.Duraction = lap.End.Sub(startLap)
		lapTime.AVG = float64(cfg.LapLength) / float64(lapTime.Duraction.Seconds())
		r.LapsStatistics = append(r.LapsStatistics, lapTime)

		startLap = lap.End
	}

	var penaltyLapAbsoluteTime time.Duration = 0
	for _, lap := range c.PenaltyLaps {
		penaltyLapAbsoluteTime += lap.End.Sub(lap.Start)
	}
	r.PenaltyLapStatictic.Duraction = penaltyLapAbsoluteTime
	if len(c.PenaltyLaps) == 0 {
		r.PenaltyLapStatictic.AVG = 0
	} else {
		r.PenaltyLapStatictic.AVG = float64(cfg.PenaltyLength) / float64(penaltyLapAbsoluteTime.Seconds())
	}

	r.Hits = c.Hits
	r.Shots = len(c.Laps) * 5
	return r
}

func formatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	ms := int(d.Milliseconds()) % 1000
	return fmt.Sprintf("%02d:%02d:%02d.%03d", h, m, s, ms)
}

func (e *Event) print() {
	if e.ExtraParams != 1 {
		return
	}
	switch e.EventID {
	case 1:
		fmt.Printf("[%s] The competitor(%d) registered\n", e.Time.Format("15:04:05.000"), e.ExtraParams)
	case 2:
		fmt.Printf("[%s] The start time for the competitor(%d) was set by a draw to %s\n", e.Time.Format("15:04:05.000"), e.ExtraParams, e.Comments)
	case 3:
		fmt.Printf("[%s] The competitor(%d) is on the start line\n", e.Time.Format("15:04:05.000"), e.ExtraParams)
	case 4:
		fmt.Printf("[%s] The competitor(%d) has started\n", e.Time.Format("15:04:05.000"), e.ExtraParams)
	case 5:
		fmt.Printf("[%s] The competitor(%d) is on the firing range(%s)\n", e.Time.Format("15:04:05.000"), e.ExtraParams, e.Comments)
	case 6:
		fmt.Printf("[%s] The target(%s) has been hit by competitor(%d)\n", e.Time.Format("15:04:05.000"), e.Comments, e.ExtraParams)
	case 7:
		fmt.Printf("[%s] The competitor(%d) left the firing range\n", e.Time.Format("15:04:05.000"), e.ExtraParams)
	case 8:
		fmt.Printf("[%s] The competitor(%d) entered the penalty laps\n", e.Time.Format("15:04:05.000"), e.ExtraParams)
	case 9:
		fmt.Printf("[%s] The competitor(%d) left the penalty laps\n", e.Time.Format("15:04:05.000"), e.ExtraParams)
	case 10:
		fmt.Printf("[%s] The competitor(%d) ended the main lap\n", e.Time.Format("15:04:05.000"), e.ExtraParams)
	case 11:
		fmt.Printf("[%s] The competitor(%d) can`t continue: %s\n", e.Time.Format("15:04:05.000"), e.ExtraParams, e.Comments)
	default:
		fmt.Printf("[%s] Unknown event(%d)\n", e.Time.Format("15:04:05.000"), e.EventID)
	}

}

func (t *Task) handleEvent(e *Event, competitors map[int]*Competitor) error {
	competitorsID := e.ExtraParams
	competitor, ok := competitors[competitorsID]
	if e.EventID != 1 && !ok {
		return fmt.Errorf("Competitor %d not found", e.ExtraParams)
	}

	switch e.EventID {
	case 1:
		competitor := NewCompetitor(competitorsID, t.cfg)
		competitors[e.ExtraParams] = competitor
	case 2:
		var err error
		competitor.ScheduledStart, err = time.Parse("15:04:05.000", e.Comments)
		if err != nil {
			return fmt.Errorf("error parsing time: %w", err)
		}
	case 3:
		//Nothing to do
	case 4:
		competitor.ActualStart = e.Time
	case 5:
		// lap := competitor.Laps[competitor.CurrentLap]
		// lap.FiringStart = e.Time
	case 6:
		competitor.Hits++
	case 7:
		// lap := competitor.Laps[competitor.CurrentLap]
		// lap.FiringEnd = e.Time
	case 8:
		competitor.PenaltyLaps = append(competitor.PenaltyLaps, PenaltyLap{
			Start: e.Time,
		})
	case 9:
		competitor.PenaltyLaps[len(competitor.PenaltyLaps)-1].End = e.Time
	case 10:
		competitor.Laps[competitor.CurrentLap].End = e.Time

		if competitor.CurrentLap+1 < len(competitor.Laps) {
			competitor.CurrentLap++
			competitor.Laps[competitor.CurrentLap].Start = e.Time
		}
	case 11:
		competitor.Status = StatusFinished
	default:
		return fmt.Errorf("unknown event %d", e.EventID)
	}

	return nil
}
