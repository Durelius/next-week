// Package ics defines iota constants for every iCalendar (ICS) property,
// component, parameter, and value type defined in:
//   - RFC 5545  – Internet Calendaring and Scheduling Core Object Specification
//   - RFC 5546  – iCalendar Transport-Independent Interoperability Protocol (iTIP)
//   - RFC 6868  – Parameter Value Encoding in iCalendar and vCard
//   - RFC 7529  – Non-Gregorian Recurrence Rules
//   - RFC 7953  – Calendar Availability (VAVAILABILITY)
//   - RFC 7986  – New Properties for iCalendar
//   - RFC 9073  – Event Publishing Extensions to iCalendar
//   - RFC 9074  – VALARM Extensions for iCalendar
//   - Common vendor extensions (Apple, Google, Microsoft)
package ics

// ─────────────────────────────────────────────────────────────────────────────
// Component names  (VBEGIN / VEND values)
// ─────────────────────────────────────────────────────────────────────────────

type Component int

const (
	ComponentVCalendar     Component = iota // VCALENDAR
	ComponentVEvent                         // VEVENT
	ComponentVTodo                          // VTODO
	ComponentVJournal                       // VJOURNAL
	ComponentVFreeBusy                      // VFREEBUSY
	ComponentVTimezone                      // VTIMEZONE
	ComponentVAlarm                         // VALARM
	ComponentVAvailability                  // VAVAILABILITY   (RFC 7953)
	ComponentAvailable                      // AVAILABLE       (RFC 7953)
	ComponentStandard                       // STANDARD        (inside VTIMEZONE)
	ComponentDaylight                       // DAYLIGHT        (inside VTIMEZONE)
	ComponentVLocation                      // VLOCATION       (RFC 9073)
	ComponentVResource                      // VRESOURCE       (RFC 9073)
)

// String returns the ICS component name.
func (c Component) String() string {
	return [...]string{
		"VCALENDAR",
		"VEVENT",
		"VTODO",
		"VJOURNAL",
		"VFREEBUSY",
		"VTIMEZONE",
		"VALARM",
		"VAVAILABILITY",
		"AVAILABLE",
		"STANDARD",
		"DAYLIGHT",
		"VLOCATION",
		"VRESOURCE",
	}[c]
}

// ─────────────────────────────────────────────────────────────────────────────
// Property names
// ─────────────────────────────────────────────────────────────────────────────

type Property int

const (
	// ── Calendar properties ──────────────────────────────────────────────────
	PropCalscale Property = iota // CALSCALE
	PropMethod                   // METHOD
	PropProdid                   // PRODID
	PropVersion                  // VERSION

	// ── Descriptive component properties ─────────────────────────────────────
	PropAttach          // ATTACH
	PropCategories      // CATEGORIES
	PropClass           // CLASS
	PropComment         // COMMENT
	PropDescription     // DESCRIPTION
	PropGeo             // GEO
	PropLocation        // LOCATION
	PropPercentComplete // PERCENT-COMPLETE
	PropPriority        // PRIORITY
	PropResources       // RESOURCES
	PropStatus          // STATUS
	PropSummary         // SUMMARY

	// ── Date and time component properties ───────────────────────────────────
	PropCompleted // COMPLETED
	PropDtend     // DTEND
	PropDue       // DUE
	PropDtstart   // DTSTART
	PropDuration  // DURATION
	PropFreebusy  // FREEBUSY
	PropTransp    // TRANSP

	// ── Timezone component properties ────────────────────────────────────────
	PropTzid         // TZID
	PropTzname       // TZNAME
	PropTzoffsetfrom // TZOFFSETFROM
	PropTzoffsetto   // TZOFFSETTO
	PropTzurl        // TZURL

	// ── Relationship component properties ────────────────────────────────────
	PropAttendee     // ATTENDEE
	PropContact      // CONTACT
	PropOrganizer    // ORGANIZER
	PropRecurrenceid // RECURRENCE-ID
	PropRelatedto    // RELATED-TO
	PropUrl          // URL
	PropUid          // UID

	// ── Recurrence component properties ──────────────────────────────────────
	PropExdate // EXDATE
	PropRdate  // RDATE
	PropRrule  // RRULE

	// ── Alarm component properties ────────────────────────────────────────────
	PropAction  // ACTION
	PropRepeat  // REPEAT
	PropTrigger // TRIGGER

	// ── Change management properties ─────────────────────────────────────────
	PropCreated      // CREATED
	PropDtstamp      // DTSTAMP
	PropLastModified // LAST-MODIFIED
	PropSequence     // SEQUENCE

	// ── Request status properties ─────────────────────────────────────────────
	PropRequestStatus // REQUEST-STATUS

	// ── RFC 7986: New properties ──────────────────────────────────────────────
	PropName            // NAME
	PropRefreshInterval // REFRESH-INTERVAL
	PropSource          // SOURCE
	PropColor           // COLOR
	PropImage           // IMAGE

	// ── RFC 9073: Event publishing extensions ─────────────────────────────────
	PropConference        // CONFERENCE
	PropParticipant       // (used inside VRESOURCE / participant sub-components)
	PropStructuredData    // STRUCTURED-DATA
	PropStyledDescription // STYLED-DESCRIPTION
	PropLocationType      // LOCATION-TYPE
	PropResourceType      // RESOURCE-TYPE

	// ── RFC 9074: VALARM extensions ───────────────────────────────────────────
	PropAcknowledged // ACKNOWLEDGED
	PropProximity    // PROXIMITY

	// ── RFC 7953: VAVAILABILITY ───────────────────────────────────────────────
	PropBusyType // BUSYTYPE

	// ── Commonly used vendor / de-facto extensions ────────────────────────────
	PropXWRCalname                   // X-WR-CALNAME       (Apple/Google)
	PropXWRCaldesc                   // X-WR-CALDESC       (Apple/Google)
	PropXWRTimezone                  // X-WR-TIMEZONE      (Apple/Google)
	PropXWRRelcalid                  // X-WR-RELCALID      (Apple)
	PropXAppleStructuredLocation     // X-APPLE-STRUCTURED-LOCATION
	PropXAppleTravelAdvisoryBehavior // X-APPLE-TRAVEL-ADVISORY-BEHAVIOR
	PropXAppleDefaultAlarm           // X-APPLE-DEFAULT-ALARM
	PropXAppleOmitFromSync           // X-APPLE-OMIT-FROM-SYNC
	PropXMicrosoftCdoBusystatus      // X-MICROSOFT-CDO-BUSYSTATUS
	PropXMicrosoftCdoIntendedstatus  // X-MICROSOFT-CDO-INTENDEDSTATUS
	PropXMicrosoftCdoAlldayevent     // X-MICROSOFT-CDO-ALLDAYEVENT
	PropXMicrosoftCdoImportance      // X-MICROSOFT-CDO-IMPORTANCE
	PropXMicrosoftCdoInsttype        // X-MICROSOFT-CDO-INSTTYPE
	PropXMicrosoftCdoOwnerapptid     // X-MICROSOFT-CDO-OWNER-APPT-ID
	PropXMicrosoftCdoApptsequence    // X-MICROSOFT-CDO-APPT-SEQUENCE
	PropXGoogleConference            // X-GOOGLE-CONFERENCE
	PropXGoogleStructuredLocation    // X-GOOGLE-STRUCTURED-LOCATION
	PropXGoogleCalendar              // X-GOOGLE-CALENDAR-ID
)

// String returns the ICS property name.
func (p Property) String() string {
	return [...]string{
		// Calendar
		"CALSCALE", "METHOD", "PRODID", "VERSION",
		// Descriptive
		"ATTACH", "CATEGORIES", "CLASS", "COMMENT", "DESCRIPTION",
		"GEO", "LOCATION", "PERCENT-COMPLETE", "PRIORITY", "RESOURCES",
		"STATUS", "SUMMARY",
		// Date/time
		"COMPLETED", "DTEND", "DUE", "DTSTART", "DURATION",
		"FREEBUSY", "TRANSP",
		// Timezone
		"TZID", "TZNAME", "TZOFFSETFROM", "TZOFFSETTO", "TZURL",
		// Relationship
		"ATTENDEE", "CONTACT", "ORGANIZER", "RECURRENCE-ID",
		"RELATED-TO", "URL", "UID",
		// Recurrence
		"EXDATE", "RDATE", "RRULE",
		// Alarm
		"ACTION", "REPEAT", "TRIGGER",
		// Change management
		"CREATED", "DTSTAMP", "LAST-MODIFIED", "SEQUENCE",
		// Request status
		"REQUEST-STATUS",
		// RFC 7986
		"NAME", "REFRESH-INTERVAL", "SOURCE", "COLOR", "IMAGE",
		// RFC 9073
		"CONFERENCE", "PARTICIPANT", "STRUCTURED-DATA", "STYLED-DESCRIPTION",

		"LOCATION-TYPE", "RESOURCE-TYPE",
		// RFC 9074
		"ACKNOWLEDGED", "PROXIMITY",
		// RFC 7953
		"BUSYTYPE",
		// Vendor extensions
		"X-WR-CALNAME", "X-WR-CALDESC", "X-WR-TIMEZONE", "X-WR-RELCALID",
		"X-APPLE-STRUCTURED-LOCATION",
		"X-APPLE-TRAVEL-ADVISORY-BEHAVIOR",
		"X-APPLE-DEFAULT-ALARM",
		"X-APPLE-OMIT-FROM-SYNC",
		"X-MICROSOFT-CDO-BUSYSTATUS",
		"X-MICROSOFT-CDO-INTENDEDSTATUS",
		"X-MICROSOFT-CDO-ALLDAYEVENT",
		"X-MICROSOFT-CDO-IMPORTANCE",
		"X-MICROSOFT-CDO-INSTTYPE",
		"X-MICROSOFT-CDO-OWNER-APPT-ID",
		"X-MICROSOFT-CDO-APPT-SEQUENCE",
		"X-GOOGLE-CONFERENCE",
		"X-GOOGLE-STRUCTURED-LOCATION",
		"X-GOOGLE-CALENDAR-ID",
	}[p]
}

// ─────────────────────────────────────────────────────────────────────────────
// Parameter names
// ─────────────────────────────────────────────────────────────────────────────

type Parameter int

const (
	// ── RFC 5545 parameters ───────────────────────────────────────────────────
	ParamAltrep        Parameter = iota // ALTREP
	ParamCn                             // CN
	ParamCutype                         // CUTYPE
	ParamDelegatedFrom                  // DELEGATED-FROM
	ParamDelegatedTo                    // DELEGATED-TO
	ParamDir                            // DIR
	ParamEncoding                       // ENCODING
	ParamFmttype                        // FMTTYPE
	ParamFbtype                         // FBTYPE
	ParamLanguage                       // LANGUAGE
	ParamMember                         // MEMBER
	ParamPartstat                       // PARTSTAT
	ParamRange                          // RANGE
	ParamRelated                        // RELATED
	ParamReltype                        // RELTYPE
	ParamRole                           // ROLE
	ParamRsvp                           // RSVP
	ParamSentby                         // SENT-BY
	ParamTzid                           // TZID  (as parameter, e.g. on DTSTART)
	ParamValue                          // VALUE

	// ── RFC 7986 parameters ───────────────────────────────────────────────────
	ParamDisplay // DISPLAY
	ParamEmail   // EMAIL
	ParamFeature // FEATURE
	ParamLabel   // LABEL

	// ── RFC 9073 parameters ───────────────────────────────────────────────────
	ParamDerived // DERIVED
	ParamGap     // GAP
	ParamOrder   // ORDER
	ParamSchema  // SCHEMA
	ParamLinkrel // LINKREL

	// ── RFC 9074 parameters ───────────────────────────────────────────────────
	ParamSize      // SIZE
	ParamFilename  // FILENAME
	ParamManagedid // MANAGED-ID

	// ── Vendor extension parameters ───────────────────────────────────────────
	ParamXAppleDefaultAlarm // X-APPLE-DEFAULT-ALARM
	ParamXNumguests         // X-NUM-GUESTS  (Google)
)

// String returns the ICS parameter name.
func (p Parameter) String() string {
	return [...]string{
		"ALTREP", "CN", "CUTYPE", "DELEGATED-FROM", "DELEGATED-TO",
		"DIR", "ENCODING", "FMTTYPE", "FBTYPE", "LANGUAGE",
		"MEMBER", "PARTSTAT", "RANGE", "RELATED", "RELTYPE",
		"ROLE", "RSVP", "SENT-BY", "TZID", "VALUE",
		// RFC 7986
		"DISPLAY", "EMAIL", "FEATURE", "LABEL",
		// RFC 9073
		"DERIVED", "GAP", "ORDER", "SCHEMA", "LINKREL",
		// RFC 9074
		"SIZE", "FILENAME", "MANAGED-ID",
		// Vendor
		"X-APPLE-DEFAULT-ALARM", "X-NUM-GUESTS",
	}[p]
}

// ─────────────────────────────────────────────────────────────────────────────
// Property value types  (VALUE= parameter values)
// ─────────────────────────────────────────────────────────────────────────────

type ValueType int

const (
	ValueBinary     ValueType = iota // BINARY
	ValueBoolean                     // BOOLEAN
	ValueCalAddress                  // CAL-ADDRESS
	ValueDate                        // DATE
	ValueDatetime                    // DATE-TIME
	ValueDuration                    // DURATION
	ValueFloat                       // FLOAT
	ValueInteger                     // INTEGER
	ValuePeriod                      // PERIOD
	ValueRecur                       // RECUR
	ValueText                        // TEXT
	ValueTime                        // TIME
	ValueUri                         // URI
	ValueUtcOffset                   // UTC-OFFSET
	ValueXName                       // X-NAME  (custom/experimental)
)

// String returns the ICS VALUE type name.
func (v ValueType) String() string {
	return [...]string{
		"BINARY", "BOOLEAN", "CAL-ADDRESS", "DATE", "DATE-TIME",
		"DURATION", "FLOAT", "INTEGER", "PERIOD", "RECUR",
		"TEXT", "TIME", "URI", "UTC-OFFSET", "X-NAME",
	}[v]
}

// ─────────────────────────────────────────────────────────────────────────────
// CALSCALE values
// ─────────────────────────────────────────────────────────────────────────────

type Calscale int

const (
	CalscaleGregorian Calscale = iota // GREGORIAN
	// RFC 7529 non-Gregorian calendars
	CalscaleChinese      // CHINESE
	CalscaleEthiopic     // ETHIOPIC
	CalscaleHebrew       // HEBREW
	CalscaleIslamic      // ISLAMIC
	CalscaleIslamicCivil // ISLAMIC-CIVIL
	CalscaleIslamicTbla  // ISLAMIC-TBLA
	CalscaleJapanese     // JAPANESE
	CalscalePersian      // PERSIAN
)

func (c Calscale) String() string {
	return [...]string{
		"GREGORIAN", "CHINESE", "ETHIOPIC", "HEBREW",
		"ISLAMIC", "ISLAMIC-CIVIL", "ISLAMIC-TBLA", "JAPANESE", "PERSIAN",
	}[c]
}

// ─────────────────────────────────────────────────────────────────────────────
// METHOD values  (RFC 5546 iTIP)
// ─────────────────────────────────────────────────────────────────────────────

type Method int

const (
	MethodPublish        Method = iota // PUBLISH
	MethodRequest                      // REQUEST
	MethodReply                        // REPLY
	MethodAdd                          // ADD
	MethodCancel                       // CANCEL
	MethodRefresh                      // REFRESH
	MethodCounter                      // COUNTER
	MethodDeclinecounter               // DECLINECOUNTER
)

func (m Method) String() string {
	return [...]string{
		"PUBLISH", "REQUEST", "REPLY", "ADD",
		"CANCEL", "REFRESH", "COUNTER", "DECLINECOUNTER",
	}[m]
}

// ─────────────────────────────────────────────────────────────────────────────
// CLASS values
// ─────────────────────────────────────────────────────────────────────────────

type Class int

const (
	ClassPublic       Class = iota // PUBLIC
	ClassPrivate                   // PRIVATE
	ClassConfidential              // CONFIDENTIAL
)

func (c Class) String() string {
	return [...]string{"PUBLIC", "PRIVATE", "CONFIDENTIAL"}[c]
}

// ─────────────────────────────────────────────────────────────────────────────
// STATUS values
// ─────────────────────────────────────────────────────────────────────────────

type Status int

const (
	// VEVENT statuses
	StatusTentative Status = iota // TENTATIVE
	StatusConfirmed               // CONFIRMED
	StatusCancelled               // CANCELLED
	// VTODO statuses
	StatusNeedsAction // NEEDS-ACTION
	StatusCompleted   // COMPLETED
	StatusInProcess   // IN-PROCESS
	// VJOURNAL statuses
	StatusDraft // DRAFT
	StatusFinal // FINAL
)

func (s Status) String() string {
	return [...]string{
		"TENTATIVE", "CONFIRMED", "CANCELLED",
		"NEEDS-ACTION", "COMPLETED", "IN-PROCESS",
		"DRAFT", "FINAL",
	}[s]
}

// ─────────────────────────────────────────────────────────────────────────────
// TRANSP values
// ─────────────────────────────────────────────────────────────────────────────

type Transp int

const (
	TranspOpaque      Transp = iota // OPAQUE
	TranspTransparent               // TRANSPARENT
)

func (t Transp) String() string {
	return [...]string{"OPAQUE", "TRANSPARENT"}[t]
}

// ─────────────────────────────────────────────────────────────────────────────
// ACTION values  (VALARM)
// ─────────────────────────────────────────────────────────────────────────────

type Action int

const (
	ActionAudio     Action = iota // AUDIO
	ActionDisplay                 // DISPLAY
	ActionEmail                   // EMAIL
	ActionProcedure               // PROCEDURE  (deprecated but still seen)
)

func (a Action) String() string {
	return [...]string{"AUDIO", "DISPLAY", "EMAIL", "PROCEDURE"}[a]
}

// ─────────────────────────────────────────────────────────────────────────────
// CUTYPE values  (calendar user type, ATTENDEE/ORGANIZER)
// ─────────────────────────────────────────────────────────────────────────────

type CUType int

const (
	CUTypeIndividual CUType = iota // INDIVIDUAL
	CUTypeGroup                    // GROUP
	CUTypeResource                 // RESOURCE
	CUTypeRoom                     // ROOM
	CUTypeUnknown                  // UNKNOWN
)

func (c CUType) String() string {
	return [...]string{"INDIVIDUAL", "GROUP", "RESOURCE", "ROOM", "UNKNOWN"}[c]
}

// ─────────────────────────────────────────────────────────────────────────────
// ROLE values
// ─────────────────────────────────────────────────────────────────────────────

type Role int

const (
	RoleChair          Role = iota // CHAIR
	RoleReqParticipant             // REQ-PARTICIPANT
	RoleOptParticipant             // OPT-PARTICIPANT
	RoleNonParticipant             // NON-PARTICIPANT
)

func (r Role) String() string {
	return [...]string{"CHAIR", "REQ-PARTICIPANT", "OPT-PARTICIPANT", "NON-PARTICIPANT"}[r]
}

// ─────────────────────────────────────────────────────────────────────────────
// PARTSTAT values
// ─────────────────────────────────────────────────────────────────────────────

type Partstat int

const (
	PartstatNeedsAction Partstat = iota // NEEDS-ACTION
	PartstatAccepted                    // ACCEPTED
	PartstatDeclined                    // DECLINED
	PartstatTentative                   // TENTATIVE
	PartstatDelegated                   // DELEGATED
	PartstatCompleted                   // COMPLETED
	PartstatInProcess                   // IN-PROCESS
)

func (p Partstat) String() string {
	return [...]string{
		"NEEDS-ACTION", "ACCEPTED", "DECLINED",
		"TENTATIVE", "DELEGATED", "COMPLETED", "IN-PROCESS",
	}[p]
}

// ─────────────────────────────────────────────────────────────────────────────
// FBTYPE values  (free/busy type)
// ─────────────────────────────────────────────────────────────────────────────

type FBType int

const (
	FBTypeFree            FBType = iota // FREE
	FBTypeBusy                          // BUSY
	FBTypeBusyUnavailable               // BUSY-UNAVAILABLE
	FBTypeBusyTentative                 // BUSY-TENTATIVE
)

func (f FBType) String() string {
	return [...]string{"FREE", "BUSY", "BUSY-UNAVAILABLE", "BUSY-TENTATIVE"}[f]
}

// ─────────────────────────────────────────────────────────────────────────────
// RELTYPE values  (relationship type, RELATED-TO)
// ─────────────────────────────────────────────────────────────────────────────

type RelType int

const (
	RelTypeParent  RelType = iota // PARENT
	RelTypeChild                  // CHILD
	RelTypeSibling                // SIBLING
	// RFC 9073 extensions
	RelTypeFinishToStart  // FINISHTOSTART
	RelTypeFinishToFinish // FINISHTOFINISH
	RelTypeStartToFinish  // STARTTOFINISH
	RelTypeStartToStart   // STARTTOSTART
	RelTypeFirst          // FIRST
	RelTypeNext           // NEXT
	RelTypeDepends        // DEPENDS-ON
	RelTypeRefid          // REFID
	RelTypeConceptual     // CONCEPT
	RelTypeLinked         // LINK
)

func (r RelType) String() string {
	return [...]string{
		"PARENT", "CHILD", "SIBLING",
		"FINISHTOSTART", "FINISHTOFINISH", "STARTTOFINISH", "STARTTOSTART",
		"FIRST", "NEXT", "DEPENDS-ON", "REFID", "CONCEPT", "LINK",
	}[r]
}

// ─────────────────────────────────────────────────────────────────────────────
// RANGE values
// ─────────────────────────────────────────────────────────────────────────────

type RangeValue int

const (
	RangeThisandfuture RangeValue = iota // THISANDFUTURE
)

func (r RangeValue) String() string { return "THISANDFUTURE" }

// ─────────────────────────────────────────────────────────────────────────────
// RELATED values  (TRIGGER parameter)
// ─────────────────────────────────────────────────────────────────────────────

type Related int

const (
	RelatedStart Related = iota // START
	RelatedEnd                  // END
)

func (r Related) String() string { return [...]string{"START", "END"}[r] }

// ─────────────────────────────────────────────────────────────────────────────
// ENCODING values
// ─────────────────────────────────────────────────────────────────────────────

type Encoding int

const (
	Encoding8bit   Encoding = iota // 8BIT
	EncodingBase64                 // BASE64
)

func (e Encoding) String() string { return [...]string{"8BIT", "BASE64"}[e] }

// ─────────────────────────────────────────────────────────────────────────────
// DISPLAY values  (RFC 7986, IMAGE property)
// ─────────────────────────────────────────────────────────────────────────────

type Display int

const (
	DisplayBadge     Display = iota // BADGE
	DisplayGraphic                  // GRAPHIC
	DisplayFullsize                 // FULLSIZE
	DisplayThumbnail                // THUMBNAIL
)

func (d Display) String() string {
	return [...]string{"BADGE", "GRAPHIC", "FULLSIZE", "THUMBNAIL"}[d]
}

// ─────────────────────────────────────────────────────────────────────────────
// FEATURE values  (RFC 7986, CONFERENCE property)
// ─────────────────────────────────────────────────────────────────────────────

type Feature int

const (
	FeatureAudio     Feature = iota // AUDIO
	FeatureChat                     // CHAT
	FeatureFeed                     // FEED
	FeatureModerator                // MODERATOR
	FeaturePhone                    // PHONE
	FeatureScreen                   // SCREEN
	FeatureVideo                    // VIDEO
)

func (f Feature) String() string {
	return [...]string{
		"AUDIO", "CHAT", "FEED", "MODERATOR", "PHONE", "SCREEN", "VIDEO",
	}[f]
}

// ─────────────────────────────────────────────────────────────────────────────
// BUSYTYPE values  (RFC 7953, VAVAILABILITY)
// ─────────────────────────────────────────────────────────────────────────────

type BusyType int

const (
	BusyTypeBusy            BusyType = iota // BUSY
	BusyTypeBusyUnavailable                 // BUSY-UNAVAILABLE
	BusyTypeBusyTentative                   // BUSY-TENTATIVE
)

func (b BusyType) String() string {
	return [...]string{"BUSY", "BUSY-UNAVAILABLE", "BUSY-TENTATIVE"}[b]
}

// ─────────────────────────────────────────────────────────────────────────────
// PROXIMITY values  (RFC 9074, VALARM extensions)
// ─────────────────────────────────────────────────────────────────────────────

type Proximity int

const (
	ProximityArrive     Proximity = iota // ARRIVE
	ProximityDepart                      // DEPART
	ProximityConnect                     // CONNECT
	ProximityDisconnect                  // DISCONNECT
)

func (p Proximity) String() string {
	return [...]string{"ARRIVE", "DEPART", "CONNECT", "DISCONNECT"}[p]
}

// ─────────────────────────────────────────────────────────────────────────────
// RRULE / RDATE frequency and weekday constants
// ─────────────────────────────────────────────────────────────────────────────

type Frequency int

const (
	FreqSecondly Frequency = iota // SECONDLY
	FreqMinutely                  // MINUTELY
	FreqHourly                    // HOURLY
	FreqDaily                     // DAILY
	FreqWeekly                    // WEEKLY
	FreqMonthly                   // MONTHLY
	FreqYearly                    // YEARLY
)

func (f Frequency) String() string {
	return [...]string{
		"SECONDLY", "MINUTELY", "HOURLY", "DAILY", "WEEKLY", "MONTHLY", "YEARLY",
	}[f]
}

type Weekday int

const (
	WeekdaySu Weekday = iota // SU
	WeekdayMo                // MO
	WeekdayTu                // TU
	WeekdayWe                // WE
	WeekdayTh                // TH
	WeekdayFr                // FR
	WeekdaySa                // SA
)

func (w Weekday) String() string {
	return [...]string{"SU", "MO", "TU", "WE", "TH", "FR", "SA"}[w]
}
