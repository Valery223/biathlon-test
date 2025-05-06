package scannerEvent

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/Valery223/biathlon-test/internal/domain"
)

// Scanner is a struct that wraps a bufio.Scanner to read and parse event data.
type Scanner struct {
	scanner *bufio.Scanner
}

// NewScanner creates and returns a new Scanner.
// It takes an io.Reader as input.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{scanner: bufio.NewScanner(r)}
}

// Scan reads the next line from the scanner, parses it, and populates the given domain.Event.
// It returns an error if scanning or parsing fails, or io.EOF if the end of the input is reached.
func (s *Scanner) Scan(e *domain.Event) error {
	if !s.scanner.Scan() {
		if err := s.scanner.Err(); err != nil {
			return err
		}
		return io.EOF
	}

	line := s.scanner.Text()
	return s.parseLine(line, e)

}

// parseLine parses a single line of event data and populates the given domain.Event.
// The line is expected to be in the format: "[HH:MM:SS.mmm] EventID ExtraParams [Comments]"
// It returns an error if the line format is invalid or parsing fails.
func (s *Scanner) parseLine(line string, e *domain.Event) error {
	parts := strings.Fields(line) // Split the line into parts by whitespace.
	if len(parts) < 3 {
		return fmt.Errorf("invalid line format: %s", line) // Return an error if the line has fewer than 3 parts.
	}

	timeStr := strings.Trim(parts[0], "[]")                // Extract and trim the time string.
	parsedTime, err := time.Parse("15:04:05.000", timeStr) // Parse the time string.
	if err != nil {
		return fmt.Errorf("invalid time format: %w", err) // Return an error if time parsing fails.
	}
	eventID, err := strconv.Atoi(parts[1]) // Parse the EventID.
	if err != nil {
		return fmt.Errorf("invalid EventID: %w", err) // Return an error if EventID parsing fails.
	}
	extraParams, err := strconv.Atoi(parts[2]) // Parse the ExtraParams.
	if err != nil {
		return fmt.Errorf("invalid ExtraParams: %w", err) // Return an error if ExtraParams parsing fails.
	}
	// Populate the event fields.
	e.Time = parsedTime
	e.ID = domain.EventID(eventID)
	e.CompetitorID = extraParams
	e.Comments = strings.Join(parts[3:], " ") // Join the remaining parts as comments.
	return nil
}
