package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Durelius/next-week/internal/avl"
	"github.com/Durelius/next-week/internal/ics"
)

func main() {
	calendars := []string{
		"https://cal.laget.se/FAIKHerr.ics",
		"https://cal.laget.se/difhandbollherr.ics",
		os.Getenv("DAISY_PERSONAL"),
		"https://www.officeholidays.com/ics-all/various",
		"https://www.vic.gov.au/sites/default/files/2026-01/Victorian-public-holiday-dates.ics",
		"https://www.vic.gov.au/school-term-dates-and-holidays-victoria",
		"https://www.vic.gov.au/sites/default/files/2025-04/2026-multifaith-calendar-Victoria.ics",
		"https://www.vic.gov.au/sites/default/files/2026-01/Parliament_sitting_dates.ics",
		"https://data.riksdagen.se/kalender/?org=kamm&akt=be&akt=ib&akt=re&akt=se&akt=tl&akt=ur&akt=ub&akt=vi&akt=oh&utformat=icalendar",
		"https://www.calendarlabs.com/ical-calendar/ics/76/US_Holidays.ics",
		"https://www.calendarlabs.com/ical-calendar/ics/38/UK_Holidays.ics",
		"https://raw.githubusercontent.com/faribe/maldives-academic-calendar/main/ical/maldives_academic_calendar.ics",
		"https://raw.githubusercontent.com/davkat1/FrenchRepublicaniCalendar/main/FrenchRepublicanCalnedar_01012025-31122025.ics",
	}
	t := avl.New[int64, map[ics.Property]*ics.ParsedLine]()
	for _, cal := range calendars {

		res, err := http.Get(cal)
		if err != nil {
			log.Fatalf("err req: %v", err)
		}
		body, err := io.ReadAll(res.Body)
		res.Body.Close()
		if res.StatusCode > 299 {
			log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
		}
		if err != nil {
			log.Fatalf("err body: %v", err)
		}
		icsStr := string(body)
		icsStr = strings.ReplaceAll(icsStr, "\r", "")
		rows := strings.SplitAfterSeq(icsStr, "\n")
		insideEvent := false
		eventMap := make(map[ics.Property]*ics.ParsedLine)
		var currentCal ics.ParsedLine
		for row := range rows {
			parsedRow, err := ics.ParseLine(row)
			if err != nil {
				continue
			}
			if parsedRow.Property == ics.PropXWRCalname {
				currentCal = parsedRow
				continue
			}
			if parsedRow.IsBegin && parsedRow.Component == ics.ComponentVEvent {
				insideEvent = true
				continue
			}
			if parsedRow.IsEnd && parsedRow.Component == ics.ComponentVEvent {
				propDtStart := eventMap[ics.PropDtstart]
				if propDtStart != nil {
					start, err := ics.DtstampToUnix(*eventMap[ics.PropDtstart])
					if err == nil {
						eventMap[ics.PropXWRCalname] = &currentCal
						t.Insert(start, eventMap)
					}
				}
				clear(eventMap)

				insideEvent = false
				continue
			}
			if !insideEvent {
				continue
			}

			eventMap[parsedRow.Property] = &parsedRow

		}
	}
	t.Print()
	log.Println(t.Size())
}
