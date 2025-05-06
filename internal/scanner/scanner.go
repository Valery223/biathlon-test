package scannerEvent

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/Valery223/biathlon-test/internal/task"
)

type Scanner struct {
	scanner *bufio.Scanner
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{scanner: bufio.NewScanner(r)}
}

func (s *Scanner) Scan(e *task.Event) error {
	if !s.scanner.Scan() {
		if err := s.scanner.Err(); err != nil {
			return err
		}
		return io.EOF
	}

	line := s.scanner.Text()
	return s.parseLine(line, e)

}

func (s *Scanner) parseLine(line string, e *task.Event) error {
	parts := strings.Fields(line)
	if len(parts) < 3 {
		return fmt.Errorf("invalid line format: %s", line)
	}

	timeStr := strings.Trim(parts[0], "[]")
	parsedTime, err := time.Parse("15:04:05.000", timeStr)
	if err != nil {
		return fmt.Errorf("invalid time format: %w", err)
	}
	eventID, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("invalid EventID: %w", err)
	}
	extraParams, err := strconv.Atoi(parts[2])
	if err != nil {
		return fmt.Errorf("invalid ExtraParams: %w", err)
	}
	e.Time = parsedTime
	e.EventID = eventID
	e.ExtraParams = extraParams
	e.Comments = strings.Join(parts[3:], " ")
	return nil
}
