package ics_test

import (
	"testing"
	"time"

	"github.com/Durelius/next-week/internal/ics"
)

// helper: parse a line and call DtstampToUnix, failing on any error
func mustDtstamp(t *testing.T, line string) int64 {
	t.Helper()
	pl, err := ics.ParseLine(line)
	if err != nil {
		t.Fatalf("ParseLine(%q): %v", line, err)
	}
	unix, err := ics.DtstampToUnix(pl)
	if err != nil {
		t.Fatalf("DtstampToUnix(%q): %v", line, err)
	}
	return unix
}

// helper: build the expected Unix timestamp from a UTC time string
func utcUnix(t *testing.T, layout, value string) int64 {
	t.Helper()
	tm, err := time.ParseInLocation(layout, value, time.UTC)
	if err != nil {
		t.Fatalf("utcUnix parse %q: %v", value, err)
	}
	return tm.Unix()
}

// ─────────────────────────────────────────────────────────────────────────────

func TestDtstampToUnix_UTC(t *testing.T) {
	got := mustDtstamp(t, "DTSTAMP:20240115T103000Z")
	want := utcUnix(t, "20060102T150405Z", "20240115T103000Z")
	if got != want {
		t.Errorf("UTC: want %d, got %d", want, got)
	}
}

func TestDtstampToUnix_Epoch(t *testing.T) {
	// Unix epoch itself
	got := mustDtstamp(t, "DTSTAMP:19700101T000000Z")
	if got != 0 {
		t.Errorf("epoch: want 0, got %d", got)
	}
}

func TestDtstampToUnix_FloatingNoTZID(t *testing.T) {
	// No Z suffix, no TZID → treated as UTC
	got := mustDtstamp(t, "DTSTAMP:20240115T103000")
	want := utcUnix(t, "20060102T150405", "20240115T103000")
	if got != want {
		t.Errorf("floating: want %d, got %d", want, got)
	}
}

func TestDtstampToUnix_WithTZIDParam(t *testing.T) {
	// TZID=Europe/Stockholm is UTC+1 in January (CET)
	got := mustDtstamp(t, "DTSTART;TZID=Europe/Stockholm:20240115T110000")

	loc, err := time.LoadLocation("Europe/Stockholm")
	if err != nil {
		t.Skip("Europe/Stockholm timezone not available on this system")
	}
	want, err := time.ParseInLocation("20060102T150405", "20240115T110000", loc)
	if err != nil {
		t.Fatal(err)
	}
	if got != want.Unix() {
		t.Errorf("TZID Stockholm: want %d, got %d (diff %ds)", want.Unix(), got, got-want.Unix())
	}
}

func TestDtstampToUnix_DateOnly(t *testing.T) {
	// DATE-only value — should be midnight UTC on that day
	got := mustDtstamp(t, "DTSTART;VALUE=DATE:20240115")
	want := utcUnix(t, "20060102", "20240115")
	if got != want {
		t.Errorf("date-only: want %d, got %d", want, got)
	}
}

func TestDtstampToUnix_CreatedProperty(t *testing.T) {
	got := mustDtstamp(t, "CREATED:20230601T080000Z")
	want := utcUnix(t, "20060102T150405Z", "20230601T080000Z")
	if got != want {
		t.Errorf("CREATED: want %d, got %d", want, got)
	}
}

func TestDtstampToUnix_LastModified(t *testing.T) {
	got := mustDtstamp(t, "LAST-MODIFIED:20231231T235959Z")
	want := utcUnix(t, "20060102T150405Z", "20231231T235959Z")
	if got != want {
		t.Errorf("LAST-MODIFIED: want %d, got %d", want, got)
	}
}

func TestDtstampToUnix_FarFuture(t *testing.T) {
	got := mustDtstamp(t, "DTEND:20991231T235959Z")
	want := utcUnix(t, "20060102T150405Z", "20991231T235959Z")
	if got != want {
		t.Errorf("far future: want %d, got %d", want, got)
	}
}

func TestDtstampToUnix_FarPast(t *testing.T) {
	got := mustDtstamp(t, "DTSTART:19000101T000000Z")
	want := utcUnix(t, "20060102T150405Z", "19000101T000000Z")
	if got != want {
		t.Errorf("far past: want %d, got %d", want, got)
	}
}

func TestDtstampToUnix_Midnight(t *testing.T) {
	got := mustDtstamp(t, "DTSTAMP:20240101T000000Z")
	want := utcUnix(t, "20060102T150405Z", "20240101T000000Z")
	if got != want {
		t.Errorf("midnight: want %d, got %d", want, got)
	}
}

func TestDtstampToUnix_EndOfDay(t *testing.T) {
	got := mustDtstamp(t, "DTSTAMP:20240101T235959Z")
	want := utcUnix(t, "20060102T150405Z", "20240101T235959Z")
	if got != want {
		t.Errorf("end of day: want %d, got %d", want, got)
	}
}

func TestDtstampToUnix_LeapDay(t *testing.T) {
	got := mustDtstamp(t, "DTSTART:20240229T120000Z")
	want := utcUnix(t, "20060102T150405Z", "20240229T120000Z")
	if got != want {
		t.Errorf("leap day: want %d, got %d", want, got)
	}
}

func TestDtstampToUnix_UnknownTZIDFallsBackToUTC(t *testing.T) {
	// Should not error — falls back to UTC silently
	pl, err := ics.ParseLine("DTSTART;TZID=Not/A/Real/Zone:20240115T110000")
	if err != nil {
		t.Fatal(err)
	}
	unix, err := ics.DtstampToUnix(pl)
	if err != nil {
		t.Fatalf("unknown TZID should not error, got: %v", err)
	}
	want := utcUnix(t, "20060102T150405", "20240115T110000")
	if unix != want {
		t.Errorf("unknown TZID fallback: want %d (UTC), got %d", want, unix)
	}
}

func TestDtstampToUnix_EmptyValue(t *testing.T) {
	pl, _ := ics.ParseLine("DTSTAMP:")
	_, err := ics.DtstampToUnix(pl)
	if err == nil {
		t.Error("expected error for empty value")
	}
}

func TestDtstampToUnix_InvalidFormat(t *testing.T) {
	pl, _ := ics.ParseLine("DTSTAMP:not-a-date")
	_, err := ics.DtstampToUnix(pl)
	if err == nil {
		t.Error("expected error for invalid date format")
	}
}

func TestDtstampToUnix_PartialDate(t *testing.T) {
	pl, _ := ics.ParseLine("DTSTAMP:2024")
	_, err := ics.DtstampToUnix(pl)
	if err == nil {
		t.Error("expected error for partial date")
	}
}

func TestDtstampToUnix_ZSuffixTakesPrecedenceOverTZID(t *testing.T) {
	// Z wins: result must be pure UTC regardless of TZID
	line := "DTSTART;TZID=America/New_York:20240115T120000Z"
	pl, err := ics.ParseLine(line)
	if err != nil {
		t.Fatal(err)
	}
	got, err := ics.DtstampToUnix(pl)
	if err != nil {
		t.Fatal(err)
	}
	want := utcUnix(t, "20060102T150405Z", "20240115T120000Z")
	if got != want {
		t.Errorf("Z suffix must beat TZID: want %d, got %d", want, got)
	}
}

func TestDtstampToUnix_Roundtrip(t *testing.T) {
	// Parse a known timestamp, convert back to time.Time, compare fields
	unix := mustDtstamp(t, "DTSTAMP:20240615T143000Z")
	tm := time.Unix(unix, 0).UTC()
	if tm.Year() != 2024 || tm.Month() != 6 || tm.Day() != 15 {
		t.Errorf("date mismatch: got %v", tm)
	}
	if tm.Hour() != 14 || tm.Minute() != 30 || tm.Second() != 0 {
		t.Errorf("time mismatch: got %v", tm)
	}
}
