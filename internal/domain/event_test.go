package domain

import (
	"fmt"
	"testing"
	"time"
)

func TestEvent_Format(t *testing.T) {
	fixedTime := time.Date(2025, 6, 6, 10, 0, 0, 0, time.UTC)
	formattedTime := fixedTime.Format("15:04:05.000")

	testCases := []struct {
		name  string
		event Event
		want  string
	}{
		{
			name:  "CompetitorRegistered",
			event: Event{Time: fixedTime, ID: EventCompetitorRegistered, CompetitorID: 1},
			want:  fmt.Sprintf("[%s] The competitor(1) registered", formattedTime),
		},
		{
			name:  "StartTimeSet",
			event: Event{Time: fixedTime, ID: EventStartTimeSet, CompetitorID: 2, Comments: "10:30:00.000"},
			want:  fmt.Sprintf("[%s] The start time for the competitor(2) was set by a draw to 10:30:00.000", formattedTime),
		},
		{
			name:  "TargetHit",
			event: Event{Time: fixedTime, ID: EventTargetHit, CompetitorID: 3, Comments: "Target5"},
			want:  fmt.Sprintf("[%s] The target(Target5) has been hit by competitor(3)", formattedTime),
		},
		{
			name:  "UnknownEvent",
			event: Event{Time: fixedTime, ID: EventID(99), CompetitorID: 4},
			want:  fmt.Sprintf("[%s] Unknown event(99) for competitor(4)", formattedTime),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.event.Format()
			if got != tc.want {
				t.Errorf("Event.Format() for %s:\ngot:  %q\nwant: %q", tc.name, got, tc.want)
			}
		})
	}
}
