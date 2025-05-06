package main

import (
	"flag"
	"log"
	"os"

	"github.com/Valery223/biathlon-test/internal/config"
	scannerEvent "github.com/Valery223/biathlon-test/internal/scanner"
	"github.com/Valery223/biathlon-test/internal/task"
)

const defaultConfigPath = "sunny_5_skiers/config.json"
const defaultEventPath = "sunny_5_skiers/events"

func main() {

	var configPath string
	var eventPath string

	flag.StringVar(&configPath, "config", defaultConfigPath, "path to config file")
	flag.StringVar(&eventPath, "events", defaultEventPath, "path to events file")
	flag.Parse()

	f, err := os.Open(eventPath)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer f.Close()
	sc := scannerEvent.NewScanner(f)

	cfg := config.MustLoadConfig(configPath)

	task := task.NewTask(cfg, sc)
	err = task.Execute()
	if err != nil {
		log.Fatalf("failed to run task: %v", err)
	}

}
