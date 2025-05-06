package config

import (
	"log"
	"time"
)

type Config struct {
	Laps          int       `json:"laps"`
	LapLength     int       `json:"lapLen"`
	PenaltyLength int       `json:"penaltyLen"`
	FiringLines   int       `json:"firingLines"`
	StartTime     time.Time `json:"start"`
	StartDelta    time.Time `json:"startDelta"`
}

func MustLoadConfig() *Config {
	startTime, err := time.Parse("15:04:05.000", "10:00:00.000")
	if err != nil {
		log.Fatalf("failed to parse start time: %v", err)
	}

	startDelta, err := time.Parse("15:04:05", "00:00:30")
	if err != nil {
		log.Fatalf("failed to parse start time: %v", err)
	}

	cfg := Config{Laps: 2,
		LapLength:     3500,
		PenaltyLength: 150,
		FiringLines:   2,
		StartTime:     startTime,
		StartDelta:    startDelta,
	}
	return &cfg
}
