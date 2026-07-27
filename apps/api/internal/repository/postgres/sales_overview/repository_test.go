package sales_overview

import (
	"testing"
	"time"
)

func TestNormalizeTrendRangeParsesStringDateBoundaries(t *testing.T) {
	start, end := normalizeTrendRange("2023-01-01", "2026-07-27", "monthly")

	if start.Format("2006-01-02 15:04:05") != "2023-01-01 00:00:00" {
		t.Fatalf("expected parsed start date, got %s", start.Format("2006-01-02 15:04:05"))
	}
	if end.Format("2006-01-02 15:04:05") != "2026-07-27 23:59:59" {
		t.Fatalf("expected parsed end date at end of day, got %s", end.Format("2006-01-02 15:04:05"))
	}
}

func TestNormalizeTrendRangeParsesTimeBoundaries(t *testing.T) {
	startInput := time.Date(2024, time.January, 10, 12, 30, 0, 0, time.UTC)
	endInput := time.Date(2024, time.March, 20, 8, 15, 0, 0, time.UTC)

	start, end := normalizeTrendRange(startInput, endInput, "monthly")

	if start.Format("2006-01-02 15:04:05") != "2024-01-10 00:00:00" {
		t.Fatalf("expected normalized start of day, got %s", start.Format("2006-01-02 15:04:05"))
	}
	if end.Format("2006-01-02 15:04:05") != "2024-03-20 23:59:59" {
		t.Fatalf("expected normalized end of day, got %s", end.Format("2006-01-02 15:04:05"))
	}
}
