package main

import (
	"log"
	"os"
	"time"

	scannerEvent "github.com/Valery223/biathlon-test/internal/scanner"
	"github.com/Valery223/biathlon-test/internal/task"
)

func main() {

	// TODO: read config file

	// TODO: struct for readinh line from file

	// TODO: task run

	// f, err := os.Open("sunny_5_skiers/events")
	f, err := os.Open("my_events")
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer f.Close()

	sc := scannerEvent.NewScanner(f)
	startTime, err := time.Parse("15:04:05.000", "09:30:00.000")
	if err != nil {
		log.Fatalf("failed to parse start time: %v", err)
	}

	startDelta, err := time.Parse("15:04:05", "00:00:30")
	if err != nil {
		log.Fatalf("failed to parse start time: %v", err)
	}

	cfg := task.Config{Laps: 1,
		LapLength:     3651,
		PenaltyLength: 50,
		FiringLines:   2,
		StartTime:     startTime,
		StartDelta:    startDelta,
	}

	task := task.NewTask(&cfg, sc)
	err = task.Execute()
	if err != nil {
		log.Fatalf("failed to run task: %v", err)
	}

}
