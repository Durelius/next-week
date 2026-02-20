package ics_test

import (
	"strings"
	"testing"

	"github.com/Durelius/next-week/internal/ics"
)

func TestParseProperty_SimpleValue(t *testing.T) {
	prop, value, err := ics.ParseProperty("SUMMARY:Team standup")
	if err != nil {
		t.Fatal(err)
	}
	if prop != ics.PropSummary {
		t.Errorf("prop: want PropSummary (%d), got %d", ics.PropSummary, prop)
	}
	if value != "Team standup" {
		t.Errorf("value: want %q, got %q", "Team standup", value)
	}
}

func TestParseProperty_WithParams(t *testing.T) {
	prop, value, err := ics.ParseProperty("DTSTART;TZID=America/New_York:20240101T090000")
	if err != nil {
		t.Fatal(err)
	}
	if prop != ics.PropDtstart {
		t.Errorf("prop: want PropDtstart (%d), got %d", ics.PropDtstart, prop)
	}
	if value != "20240101T090000" {
		t.Errorf("value: want %q, got %q", "20240101T090000", value)
	}
}

func TestParseProperty_ColonInValue(t *testing.T) {
	// URL values contain colons — only the FIRST one is the separator
	prop, value, err := ics.ParseProperty("URL:https://example.com/calendar?foo=bar")
	if err != nil {
		t.Fatal(err)
	}
	if prop != ics.PropUrl {
		t.Errorf("prop: want PropUrl, got %d", prop)
	}
	if value != "https://example.com/calendar?foo=bar" {
		t.Errorf("value: want full URL, got %q", value)
	}
}

func TestParseProperty_QuotedColonInParam(t *testing.T) {
	// ALTREP param value contains a colon inside quotes — must not split there
	line := `DESCRIPTION;ALTREP="http://example.com/desc:full":Short description`
	pl, err := ics.ParseLine(line)
	if err != nil {
		t.Fatal(err)
	}
	if pl.Property != ics.PropDescription {
		t.Errorf("prop: want PropDescription, got %d", pl.Property)
	}
	if pl.Value != "Short description" {
		t.Errorf("value: want %q, got %q", "Short description", pl.Value)
	}
}

func TestParseProperty_BEGIN(t *testing.T) {
	pl, err := ics.ParseLine("BEGIN:VEVENT")
	if err != nil {
		t.Fatal(err)
	}
	if !pl.IsBegin {
		t.Error("IsBegin: want true")
	}
	if pl.IsEnd {
		t.Error("IsEnd: want false")
	}
	if pl.Component != ics.ComponentVEvent {
		t.Errorf("Component: want ComponentVEvent, got %d", pl.Component)
	}
	if pl.Property != ics.UnknownProperty {
		t.Errorf("Property: want UnknownProperty for BEGIN line, got %d", pl.Property)
	}
}

func TestParseProperty_END(t *testing.T) {
	pl, err := ics.ParseLine("END:VCALENDAR")
	if err != nil {
		t.Fatal(err)
	}
	if !pl.IsEnd {
		t.Error("IsEnd: want true")
	}
	if pl.Component != ics.ComponentVCalendar {
		t.Errorf("Component: want ComponentVCalendar, got %d", pl.Component)
	}
}

func TestParseProperty_UnknownComponent(t *testing.T) {
	pl, err := ics.ParseLine("BEGIN:VCUSTOM")
	if err != nil {
		t.Fatal(err)
	}
	if pl.Component != ics.UnknownComponent {
		t.Errorf("Component: want UnknownComponent (-1) for unknown component name, got %d", pl.Component)
	}
}

func TestParseProperty_UnknownProperty(t *testing.T) {
	pl, err := ics.ParseLine("X-CUSTOM-TAG:some value")
	if err != nil {
		t.Fatal(err)
	}
	if !pl.IsUnknown {
		t.Error("IsUnknown: want true for unrecognised X- property")
	}
	if pl.Property != ics.UnknownProperty {
		t.Errorf("Property: want UnknownProperty (-1), got %d", pl.Property)
	}
	if pl.Value != "some value" {
		t.Errorf("Value: want %q, got %q", "some value", pl.Value)
	}
	if pl.RawName != "X-CUSTOM-TAG" {
		t.Errorf("RawName: want %q, got %q", "X-CUSTOM-TAG", pl.RawName)
	}
}

func TestParseProperty_NoColon(t *testing.T) {
	_, _, err := ics.ParseProperty("THISISNOTVALID")
	if err == nil {
		t.Error("expected error for line with no ':'")
	}
}

func TestParseProperty_EmptyLine(t *testing.T) {
	_, _, err := ics.ParseProperty("")
	if err == nil {
		t.Error("expected error for empty line")
	}
}

func TestParseProperty_MultipleParams(t *testing.T) {
	line := "ATTENDEE;ROLE=REQ-PARTICIPANT;RSVP=TRUE;CN=Alice:mailto:alice@example.com"
	pl, err := ics.ParseLine(line)
	if err != nil {
		t.Fatal(err)
	}
	if pl.Property != ics.PropAttendee {
		t.Errorf("prop: want PropAttendee, got %d", pl.Property)
	}
	if pl.Value != "mailto:alice@example.com" {
		t.Errorf("value: want mailto URI, got %q", pl.Value)
	}
	if len(pl.ParsedParams) != 3 {
		t.Fatalf("ParsedParams: want 3, got %d: %v", len(pl.ParsedParams), pl.ParsedParams)
	}
	byKey := map[string]string{}
	for _, p := range pl.ParsedParams {
		byKey[p.Key] = p.Value
	}
	if byKey["ROLE"] != "REQ-PARTICIPANT" {
		t.Errorf("ROLE param: want REQ-PARTICIPANT, got %q", byKey["ROLE"])
	}
	if byKey["RSVP"] != "TRUE" {
		t.Errorf("RSVP param: want TRUE, got %q", byKey["RSVP"])
	}
	if byKey["CN"] != "Alice" {
		t.Errorf("CN param: want Alice, got %q", byKey["CN"])
	}
}

func TestParseProperty_VendorExtension(t *testing.T) {
	prop, value, err := ics.ParseProperty("X-WR-CALNAME:My Calendar")
	if err != nil {
		t.Fatal(err)
	}
	if prop != ics.PropXWRCalname {
		t.Errorf("prop: want PropXWRCalname, got %d", prop)
	}
	if value != "My Calendar" {
		t.Errorf("value: want %q, got %q", "My Calendar", value)
	}
}

func TestParseProperty_CRLFStripped(t *testing.T) {
	prop, value, err := ics.ParseProperty("SUMMARY:Standup\r\n")
	if err != nil {
		t.Fatal(err)
	}
	if prop != ics.PropSummary {
		t.Errorf("prop: want PropSummary, got %d", prop)
	}
	if value != "Standup" {
		t.Errorf("value should have CRLF stripped, got %q", value)
	}
}

func TestParseProperty_EmptyValue(t *testing.T) {
	prop, value, err := ics.ParseProperty("COMMENT:")
	if err != nil {
		t.Fatal(err)
	}
	if prop != ics.PropComment {
		t.Errorf("prop: want PropComment, got %d", prop)
	}
	if value != "" {
		t.Errorf("value: want empty string, got %q", value)
	}
}

func TestParseProperty_AllKnownProperties(t *testing.T) {
	// Every known property should round-trip through ParseLine.
	cases := []struct {
		line string
		want ics.Property
	}{
		{"CALSCALE:GREGORIAN", ics.PropCalscale},
		{"METHOD:PUBLISH", ics.PropMethod},
		{"PRODID:-//Example//EN", ics.PropProdid},
		{"VERSION:2.0", ics.PropVersion},
		{"ATTACH:https://example.com/doc.pdf", ics.PropAttach},
		{"CATEGORIES:WORK,MEETING", ics.PropCategories},
		{"CLASS:PUBLIC", ics.PropClass},
		{"COMMENT:A comment", ics.PropComment},
		{"DESCRIPTION:Event description", ics.PropDescription},
		{"GEO:37.386013;-122.082932", ics.PropGeo},
		{"LOCATION:Conference Room 1", ics.PropLocation},
		{"PERCENT-COMPLETE:50", ics.PropPercentComplete},
		{"PRIORITY:1", ics.PropPriority},
		{"RESOURCES:Projector", ics.PropResources},
		{"STATUS:CONFIRMED", ics.PropStatus},
		{"SUMMARY:Monthly review", ics.PropSummary},
		{"COMPLETED:20240101T120000Z", ics.PropCompleted},
		{"DTEND:20240101T170000Z", ics.PropDtend},
		{"DUE:20240201T000000Z", ics.PropDue},
		{"DTSTART:20240101T090000Z", ics.PropDtstart},
		{"DURATION:PT1H", ics.PropDuration},
		{"FREEBUSY:20240101T090000Z/PT1H", ics.PropFreebusy},
		{"TRANSP:OPAQUE", ics.PropTransp},
		{"TZID:Europe/Stockholm", ics.PropTzid},
		{"TZNAME:CET", ics.PropTzname},
		{"TZOFFSETFROM:+0100", ics.PropTzoffsetfrom},
		{"TZOFFSETTO:+0200", ics.PropTzoffsetto},
		{"TZURL:https://tzurl.org/zoneinfo/Europe/Stockholm", ics.PropTzurl},
		{"ATTENDEE:mailto:bob@example.com", ics.PropAttendee},
		{"CONTACT:Bob <bob@example.com>", ics.PropContact},
		{"ORGANIZER:mailto:alice@example.com", ics.PropOrganizer},
		{"RECURRENCE-ID:20240101T090000Z", ics.PropRecurrenceid},
		{"RELATED-TO:uid-12345", ics.PropRelatedto},
		{"URL:https://example.com", ics.PropUrl},
		{"UID:abc123@example.com", ics.PropUid},
		{"EXDATE:20240201T090000Z", ics.PropExdate},
		{"RDATE:20240301T090000Z", ics.PropRdate},
		{"RRULE:FREQ=WEEKLY;BYDAY=MO", ics.PropRrule},
		{"ACTION:EMAIL", ics.PropAction},
		{"REPEAT:2", ics.PropRepeat},
		{"TRIGGER:-PT15M", ics.PropTrigger},
		{"CREATED:20231201T000000Z", ics.PropCreated},
		{"DTSTAMP:20240101T000000Z", ics.PropDtstamp},
		{"LAST-MODIFIED:20240101T000000Z", ics.PropLastModified},
		{"SEQUENCE:1", ics.PropSequence},
		{"REQUEST-STATUS:2.0;Success", ics.PropRequestStatus},
		{"NAME:My Calendar", ics.PropName},
		{"REFRESH-INTERVAL;VALUE=DURATION:PT12H", ics.PropRefreshInterval},
		{"SOURCE:https://example.com/calendar.ics", ics.PropSource},
		{"COLOR:navy", ics.PropColor},
		{"IMAGE:https://example.com/image.png", ics.PropImage},
		{"CONFERENCE:https://meet.example.com/room", ics.PropConference},
		{"STRUCTURED-DATA:some data", ics.PropStructuredData},
		{"STYLED-DESCRIPTION:some styled text", ics.PropStyledDescription},
		{"LOCATION-TYPE:ONLINE", ics.PropLocationType},
		{"RESOURCE-TYPE:ROOM", ics.PropResourceType},
		{"ACKNOWLEDGED:20240101T090000Z", ics.PropAcknowledged},
		{"PROXIMITY:ARRIVE", ics.PropProximity},
		{"BUSYTYPE:BUSY", ics.PropBusyType},
		{"X-WR-CALNAME:Work Cal", ics.PropXWRCalname},
		{"X-WR-CALDESC:Work calendar", ics.PropXWRCaldesc},
		{"X-WR-TIMEZONE:America/New_York", ics.PropXWRTimezone},
		{"X-WR-RELCALID:xyz", ics.PropXWRRelcalid},
		{"X-APPLE-DEFAULT-ALARM:TRUE", ics.PropXAppleDefaultAlarm},
		{"X-APPLE-OMIT-FROM-SYNC:TRUE", ics.PropXAppleOmitFromSync},
		{"X-MICROSOFT-CDO-BUSYSTATUS:BUSY", ics.PropXMicrosoftCdoBusystatus},
		{"X-MICROSOFT-CDO-ALLDAYEVENT:TRUE", ics.PropXMicrosoftCdoAlldayevent},
		{"X-GOOGLE-CONFERENCE:https://meet.google.com/abc", ics.PropXGoogleConference},
		{"X-GOOGLE-CALENDAR-ID:primary", ics.PropXGoogleCalendar},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.line[:strings.IndexByte(tc.line, ':')], func(t *testing.T) {
			pl, err := ics.ParseLine(tc.line)
			if err != nil {
				t.Fatalf("ParseLine error: %v", err)
			}
			if pl.Property != tc.want {
				t.Errorf("property: want %d (%s), got %d", tc.want, tc.want.String(), pl.Property)
			}
		})
	}
}
