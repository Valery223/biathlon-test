package task

import (
	"fmt"
	"io"
	"log"

	"github.com/Valery223/biathlon-test/internal/config"
	"github.com/Valery223/biathlon-test/internal/domain"
	"github.com/Valery223/biathlon-test/internal/eventproccesor"
	"github.com/Valery223/biathlon-test/internal/reporting"
)

type ScannerEvent interface {
	Scan(*domain.Event) error
}

type Task struct {
	cfg     *config.Config
	scanner ScannerEvent
}

func NewTask(cfg *config.Config, scanner ScannerEvent) *Task {
	return &Task{
		cfg:     cfg,
		scanner: scanner,
	}
}

// Execute runs the main simulation loop.
// It continuously scans for events, handles them, and once all events are processed,
// it prints a report for each competitor.
// It returns an error if a critical issue occurs during event scanning or processing.
func (t Task) Execute() error {
	mapCompetitors := make(map[int]*domain.Competitor)
	err := t.processAllEvents(mapCompetitors)
	if err != nil {
		return fmt.Errorf("error processing events: %w", err)
	}

	fmt.Println("Final reports")
	for _, competitor := range mapCompetitors {
		report := reporting.CalculateReport(*competitor, t.cfg)
		fmt.Println(report)
	}

	fmt.Println("End of task")
	return nil
}

func (t Task) processAllEvents(mapCompetitors map[int]*domain.Competitor) error {
	for {
		event := &domain.Event{}
		err := t.scanner.Scan(event)
		if err != nil {
			if err == io.EOF {
				log.Println("End of file reached")
				break
			}
			return fmt.Errorf("error scanning event: %w", err)

		}

		err = t.handleAndShowEvent(event, mapCompetitors)

		if err != nil {
			return fmt.Errorf("error handling event: %w", err)
		}

	}

	return nil
}

func (t Task) handleAndShowEvent(event *domain.Event, mapCompetitors map[int]*domain.Competitor) error {
	fmt.Println(event.Format())
	err := eventproccesor.HandleEvent(event, mapCompetitors, t.cfg.Laps)

	return err
}
