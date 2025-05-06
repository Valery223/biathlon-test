package config

import (
	"encoding/json"
	"log"
	"os"
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

func MustLoadConfig(configPath string) *Config {

	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}

	// Временная структура для парсинга строковых дат
	type tempConfig struct {
		Laps          int    `json:"laps"`
		LapLength     int    `json:"lapLen"`
		PenaltyLength int    `json:"penaltyLen"`
		FiringLines   int    `json:"firingLines"`
		StartTime     string `json:"start"`
		StartDelta    string `json:"startDelta"`
	}

	var tmp tempConfig
	if err := json.Unmarshal(data, &tmp); err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	startTime, err := time.Parse("15:04:05.000", tmp.StartTime)
	if err != nil {
		log.Fatalf("failed to parse start time: %v", err)
	}

	startDelta, err := time.Parse("15:04:05", tmp.StartDelta)
	if err != nil {
		log.Fatalf("failed to parse start delta: %v", err)
	}

	return &Config{
		Laps:          tmp.Laps,
		LapLength:     tmp.LapLength,
		PenaltyLength: tmp.PenaltyLength,
		FiringLines:   tmp.FiringLines,
		StartTime:     startTime,
		StartDelta:    startDelta,
	}
}
