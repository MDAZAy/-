package validators

import (
	"testing"
	"time"
)

func TestIsOverlapping(t *testing.T) {
	startA := time.Date(2026, 5, 1, 10, 0, 0, 0, time.UTC)
	endA := startA.Add(time.Hour)
	startB := time.Date(2026, 5, 1, 10, 30, 0, 0, time.UTC)
	endB := startB.Add(time.Hour)

	if !IsOverlapping(startA, endA, startB, endB) {
		t.Fatal("expected ranges to overlap")
	}
}
