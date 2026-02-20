package ics

import (
	"errors"
	"fmt"
	"time"
)

// DtstampToUnix converts a parsed DTSTAMP (or any date/time property) line
// into a Unix timestamp (seconds since epoch).
//
// Supported value formats per RFC 5545 §3.3.5:
//
//	DATE-TIME (UTC)       20060102T150405Z
//	DATE-TIME (floating)  20060102T150405        resolved as UTC
//	DATE-TIME (local)     20060102T150405        resolved with TZID param if present
//	DATE (all-day)        20060102               interpreted as 00:00:00 UTC on that day
//
// Returns an error if:
//   - the ParsedLine is not a DTSTAMP (or similar date/time property)
//   - the value string cannot be parsed as any known ICS date/time format
func DtstampToUnix(pl ParsedLine) (int64, error) {
	return DateTimeToUnix(pl)
}

// DateTimeToUnix is the generic form of DtstampToUnix. It accepts any
// ParsedLine whose Value contains an ICS DATE-TIME or DATE string and returns
// the corresponding Unix timestamp.
//
// It works correctly for:
//   - DTSTAMP, DTSTART, DTEND, DUE, COMPLETED, CREATED, LAST-MODIFIED,
//     RECURRENCE-ID, EXDATE, RDATE, ACKNOWLEDGED, and any X- property
//     that stores a date/time value.
func DateTimeToUnix(pl ParsedLine) (int64, error) {
	if pl.Value == "" {
		return 0, errors.New("ics: empty value in ParsedLine")
	}

	raw := pl.Value

	// ── 1. Determine the timezone to use ─────────────────────────────────────
	//
	// Priority order (RFC 5545 §3.2.19):
	//   a) Trailing 'Z' suffix → always UTC, ignore any TZID param
	//   b) TZID parameter      → load the named location
	//   c) No suffix, no TZID  → treat as floating / UTC

	var loc *time.Location = time.UTC

	if len(raw) > 0 && raw[len(raw)-1] == 'Z' {
		// UTC — location stays time.UTC, keep the 'Z' in the string so the
		// format patterns below can match it.
		loc = time.UTC
	} else {
		// Look for a TZID parameter.
		tzidVal := paramValue(pl.ParsedParams, "TZID")
		if tzidVal != "" {
			loaded, err := time.LoadLocation(tzidVal)
			if err != nil {
				// Unknown TZID — fall back to UTC rather than hard-failing,
				// because many real-world files use non-IANA names.
				loc = time.UTC
			} else {
				loc = loaded
			}
		}
	}

	// ── 2. Try each known format in order of specificity ─────────────────────

	formats := []string{
		"20060102T150405Z", // DATE-TIME UTC  (with Z)
		"20060102T150405",  // DATE-TIME local / floating
		"20060102",         // DATE only
	}

	for _, layout := range formats {
		if len(raw) != len(layout) {
			continue
		}
		t, err := time.ParseInLocation(layout, raw, loc)
		if err == nil {
			return t.Unix(), nil
		}
	}

	return 0, fmt.Errorf("ics: cannot parse date/time value %q", raw)
}

// ─────────────────────────────────────────────────────────────────────────────
// internal helper
// ─────────────────────────────────────────────────────────────────────────────

// paramValue finds the first parameter with the given key (case-insensitive)
// and returns its value, stripping any surrounding double quotes.
func paramValue(params []Param, key string) string {
	for _, p := range params {
		if p.Key == key {
			v := p.Value
			if len(v) >= 2 && v[0] == '"' && v[len(v)-1] == '"' {
				v = v[1 : len(v)-1]
			}
			return v
		}
	}
	return ""
}
