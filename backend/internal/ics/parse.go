package ics

import (
	"errors"
	"strings"
)

// ─────────────────────────────────────────────────────────────────────────────
// ParsedLine holds the result of parsing one logical ICS content line.
// ─────────────────────────────────────────────────────────────────────────────

// ParsedLine is the result of parsing one logical ICS content line.
//
// A content line has the form:
//
//	NAME[;param1=v1;param2=v2]:value
//
// Examples:
//
//	DTSTART;TZID=America/New_York:20240101T090000
//	SUMMARY:Team standup
//	BEGIN:VEVENT
//	ATTENDEE;ROLE=REQ-PARTICIPANT;RSVP=TRUE:mailto:alice@example.com
type ParsedLine struct {
	// Property is the matched Property iota.
	// -1 if the tag is not a known property (e.g. BEGIN/END or an unrecognised X- tag).
	Property Property

	// IsBegin is true when the line is "BEGIN:..."
	IsBegin bool
	// IsEnd is true when the line is "END:..."
	IsEnd bool

	// Component is set (and valid) when IsBegin or IsEnd is true.
	// -1 if the component name is unrecognised.
	Component Component

	// RawName is the property name exactly as it appeared in the line (upper-cased).
	// For "DTSTART;TZID=America/New_York:…" this is "DTSTART".
	RawName string

	// Params holds the raw parameter string, i.e. everything between the first
	// ';' and the ':'.  Empty string when no parameters are present.
	// Example: "TZID=America/New_York;LANGUAGE=en"
	Params string

	// ParsedParams is Params split into individual key=value pairs.
	// The key is upper-cased; the value is kept as-is.
	ParsedParams []Param

	// Value is the raw value string after the ':' (not unescaped).
	Value string

	// IsUnknown is true when the property name was not found in the lookup table.
	IsUnknown bool
}

// Param is a single ;KEY=VALUE parameter on a content line.
type Param struct {
	Key   string // upper-cased parameter name
	Value string // raw value (may be quoted)
}

// ─────────────────────────────────────────────────────────────────────────────
// Sentinel values for "not found"
// ─────────────────────────────────────────────────────────────────────────────

const (
	UnknownProperty  Property  = -1
	UnknownComponent Component = -1
)

// ─────────────────────────────────────────────────────────────────────────────
// Lookup tables — built once at init time
// ─────────────────────────────────────────────────────────────────────────────

var propertyByName map[string]Property
var componentByName map[string]Component

func init() {
	// ── Property lookup ───────────────────────────────────────────────────────
	propertyByName = map[string]Property{
		// Calendar
		"CALSCALE": PropCalscale,
		"METHOD":   PropMethod,
		"PRODID":   PropProdid,
		"VERSION":  PropVersion,
		// Descriptive
		"ATTACH":           PropAttach,
		"CATEGORIES":       PropCategories,
		"CLASS":            PropClass,
		"COMMENT":          PropComment,
		"DESCRIPTION":      PropDescription,
		"GEO":              PropGeo,
		"LOCATION":         PropLocation,
		"PERCENT-COMPLETE": PropPercentComplete,
		"PRIORITY":         PropPriority,
		"RESOURCES":        PropResources,
		"STATUS":           PropStatus,
		"SUMMARY":          PropSummary,
		// Date/time
		"COMPLETED": PropCompleted,
		"DTEND":     PropDtend,
		"DUE":       PropDue,
		"DTSTART":   PropDtstart,
		"DURATION":  PropDuration,
		"FREEBUSY":  PropFreebusy,
		"TRANSP":    PropTransp,
		// Timezone
		"TZID":         PropTzid,
		"TZNAME":       PropTzname,
		"TZOFFSETFROM": PropTzoffsetfrom,
		"TZOFFSETTO":   PropTzoffsetto,
		"TZURL":        PropTzurl,
		// Relationship
		"ATTENDEE":      PropAttendee,
		"CONTACT":       PropContact,
		"ORGANIZER":     PropOrganizer,
		"RECURRENCE-ID": PropRecurrenceid,
		"RELATED-TO":    PropRelatedto,
		"URL":           PropUrl,
		"UID":           PropUid,
		// Recurrence
		"EXDATE": PropExdate,
		"RDATE":  PropRdate,
		"RRULE":  PropRrule,
		// Alarm
		"ACTION":  PropAction,
		"REPEAT":  PropRepeat,
		"TRIGGER": PropTrigger,
		// Change management
		"CREATED":       PropCreated,
		"DTSTAMP":       PropDtstamp,
		"LAST-MODIFIED": PropLastModified,
		"SEQUENCE":      PropSequence,
		// Request status
		"REQUEST-STATUS": PropRequestStatus,
		// RFC 7986
		"NAME":             PropName,
		"REFRESH-INTERVAL": PropRefreshInterval,
		"SOURCE":           PropSource,
		"COLOR":            PropColor,
		"IMAGE":            PropImage,
		// RFC 9073
		"CONFERENCE":         PropConference,
		"PARTICIPANT":        PropParticipant,
		"STRUCTURED-DATA":    PropStructuredData,
		"STYLED-DESCRIPTION": PropStyledDescription,
		"LOCATION-TYPE":      PropLocationType,
		"RESOURCE-TYPE":      PropResourceType,
		// RFC 9074
		"ACKNOWLEDGED": PropAcknowledged,
		"PROXIMITY":    PropProximity,
		// RFC 7953
		"BUSYTYPE": PropBusyType,
		// Vendor extensions
		"X-WR-CALNAME":                     PropXWRCalname,
		"X-WR-CALDESC":                     PropXWRCaldesc,
		"X-WR-TIMEZONE":                    PropXWRTimezone,
		"X-WR-RELCALID":                    PropXWRRelcalid,
		"X-APPLE-STRUCTURED-LOCATION":      PropXAppleStructuredLocation,
		"X-APPLE-TRAVEL-ADVISORY-BEHAVIOR": PropXAppleTravelAdvisoryBehavior,
		"X-APPLE-DEFAULT-ALARM":            PropXAppleDefaultAlarm,
		"X-APPLE-OMIT-FROM-SYNC":           PropXAppleOmitFromSync,
		"X-MICROSOFT-CDO-BUSYSTATUS":       PropXMicrosoftCdoBusystatus,
		"X-MICROSOFT-CDO-INTENDEDSTATUS":   PropXMicrosoftCdoIntendedstatus,
		"X-MICROSOFT-CDO-ALLDAYEVENT":      PropXMicrosoftCdoAlldayevent,
		"X-MICROSOFT-CDO-IMPORTANCE":       PropXMicrosoftCdoImportance,
		"X-MICROSOFT-CDO-INSTTYPE":         PropXMicrosoftCdoInsttype,
		"X-MICROSOFT-CDO-OWNER-APPT-ID":    PropXMicrosoftCdoOwnerapptid,
		"X-MICROSOFT-CDO-APPT-SEQUENCE":    PropXMicrosoftCdoApptsequence,
		"X-GOOGLE-CONFERENCE":              PropXGoogleConference,
		"X-GOOGLE-STRUCTURED-LOCATION":     PropXGoogleStructuredLocation,
		"X-GOOGLE-CALENDAR-ID":             PropXGoogleCalendar,
	}

	// ── Component lookup ──────────────────────────────────────────────────────
	componentByName = map[string]Component{
		"VCALENDAR":     ComponentVCalendar,
		"VEVENT":        ComponentVEvent,
		"VTODO":         ComponentVTodo,
		"VJOURNAL":      ComponentVJournal,
		"VFREEBUSY":     ComponentVFreeBusy,
		"VTIMEZONE":     ComponentVTimezone,
		"VALARM":        ComponentVAlarm,
		"VAVAILABILITY": ComponentVAvailability,
		"AVAILABLE":     ComponentAvailable,
		"STANDARD":      ComponentStandard,
		"DAYLIGHT":      ComponentDaylight,
		"VLOCATION":     ComponentVLocation,
		"VRESOURCE":     ComponentVResource,
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// ParseLine parses a single logical ICS content line.
//
// A "logical line" may have been assembled from multiple physical lines by the
// caller via RFC 5545 line-unfolding (joining CRLF + SPACE/TAB continuations)
// before calling this function.
//
// Returns an error only for structurally invalid lines (no ':' separator).
// Unrecognised property names are returned with IsUnknown=true rather than an
// error, so the caller can still access RawName and Value.
// ─────────────────────────────────────────────────────────────────────────────
func ParseLine(line string) (ParsedLine, error) {
	line = strings.TrimRight(line, "\r\n")
	if line == "" {
		return ParsedLine{}, errors.New("ics: empty line")
	}

	// ── 1. Split name+params from value on the FIRST unquoted ':' ─────────────
	nameAndParams, value, found := splitOnFirstUnquotedColon(line)
	if !found {
		return ParsedLine{}, errors.New("ics: no ':' separator found in line: " + line)
	}

	// ── 2. Split name from parameters on the FIRST ';' ────────────────────────
	var rawName, rawParams string
	if idx := strings.IndexByte(nameAndParams, ';'); idx >= 0 {
		rawName = nameAndParams[:idx]
		rawParams = nameAndParams[idx+1:]
	} else {
		rawName = nameAndParams
	}
	rawName = strings.ToUpper(strings.TrimSpace(rawName))

	result := ParsedLine{
		Property:  UnknownProperty,
		Component: UnknownComponent,
		RawName:   rawName,
		Params:    rawParams,
		Value:     value,
	}

	// ── 3. Parse individual parameters ───────────────────────────────────────
	if rawParams != "" {
		result.ParsedParams = parseParams(rawParams)
	}

	// ── 4. Handle BEGIN / END specially ───────────────────────────────────────
	switch rawName {
	case "BEGIN":
		result.IsBegin = true
		compName := strings.ToUpper(strings.TrimSpace(value))
		if c, ok := componentByName[compName]; ok {
			result.Component = c
		}
		return result, nil

	case "END":
		result.IsEnd = true
		compName := strings.ToUpper(strings.TrimSpace(value))
		if c, ok := componentByName[compName]; ok {
			result.Component = c
		}
		return result, nil
	}

	// ── 5. Look up the property iota ─────────────────────────────────────────
	if prop, ok := propertyByName[rawName]; ok {
		result.Property = prop
	} else {
		result.IsUnknown = true
	}

	return result, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// ParseProperty is a convenience wrapper that returns just the Property iota
// and the raw value string.
//
//	prop, value, err := ics.ParseProperty("DTSTART;TZID=Europe/Stockholm:20240515T083000")
//
// Returns UnknownProperty (-1) when the tag is not recognised.
// Returns an error for structurally invalid lines.
// ─────────────────────────────────────────────────────────────────────────────
func ParseProperty(line string) (Property, string, error) {
	pl, err := ParseLine(line)
	if err != nil {
		return UnknownProperty, "", err
	}
	return pl.Property, pl.Value, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// internal helpers
// ─────────────────────────────────────────────────────────────────────────────

// splitOnFirstUnquotedColon splits s at the first ':' that is not inside a
// double-quoted string.  RFC 5545 allows parameter values like:
//
//	ALTREP="http://example.com/calendar":some text
//
// so we must skip over quoted regions.
func splitOnFirstUnquotedColon(s string) (before, after string, found bool) {
	inQuote := false
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '"':
			inQuote = !inQuote
		case ':':
			if !inQuote {
				return s[:i], s[i+1:], true
			}
		}
	}
	return "", "", false
}

// parseParams splits a raw parameter string (everything after the first ';' up
// to the ':') into individual Param{Key, Value} pairs.
//
// Input example:  TZID=America/New_York;LANGUAGE=en;CN="Alice Smith"
func parseParams(raw string) []Param {
	var params []Param
	// Split on ';' that are not inside quotes.
	parts := splitUnquoted(raw, ';')
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		idx := strings.IndexByte(part, '=')
		if idx < 0 {
			// Bare parameter name with no value — store it as-is.
			params = append(params, Param{Key: strings.ToUpper(part), Value: ""})
			continue
		}
		key := strings.ToUpper(strings.TrimSpace(part[:idx]))
		val := strings.TrimSpace(part[idx+1:])
		params = append(params, Param{Key: key, Value: val})
	}
	return params
}

// splitUnquoted splits s on sep characters that are not inside double quotes.
func splitUnquoted(s string, sep byte) []string {
	var parts []string
	inQuote := false
	start := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '"':
			inQuote = !inQuote
		default:
			if s[i] == sep && !inQuote {
				parts = append(parts, s[start:i])
				start = i + 1
			}
		}
	}
	parts = append(parts, s[start:])
	return parts
}
