package main

import (
	"log"
	"os"

	"github.com/Valery223/biathlon-test/internal/config"
	scannerEvent "github.com/Valery223/biathlon-test/internal/scanner"
	"github.com/Valery223/biathlon-test/internal/task"
)

func main() {

	f, err := os.Open("sunny_5_skiers/events")
	// f, err := os.Open("my_events")
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer f.Close()

	sc := scannerEvent.NewScanner(f)

	cfg := config.MustLoadConfig()

	task := task.NewTask(cfg, sc)
	err = task.Execute()
	if err != nil {
		log.Fatalf("failed to run task: %v", err)
	}

}
