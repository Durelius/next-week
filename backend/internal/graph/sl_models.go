package graph

const (
	PATH_AGENCY    = "/data/sl_agency.csv"
	PATH_ROUTES    = "/data/sl_routes.csv"
	PATH_STOPTIMES = "/data/sl_stop_times.csv"
	PATH_STOPS     = "/data/sl_stops.csv"
	PATH_TRIPS     = "/data/sl_trips.csv"
)

type Agency struct {
	AgencyID       string `csv:"agency_id" json:"agencyId"`
	AgencyName     string `csv:"agency_name" json:"agencyName"`
	AgencyURL      string `csv:"agency_url" json:"agencyUrl"`
	AgencyTimezone string `csv:"agency_timezone" json:"agencyTimezone"`
	AgencyLanguage string `csv:"agency_lang" json:"agencyLang"`
}

type Routes struct {
	RouteID        string `csv:"route_id" json:"routeId"`
	AgencyID       string `csv:"agency_id" json:"agencyId"`
	RouteShortName string `csv:"route_short_name" json:"routeShortName"`
	RouteLongName  string `csv:"route_long_name" json:"routeLongName"`
	RouteType      string `csv:"route_type" json:"routeType"`
	RouteURL       string `csv:"route_url" json:"routeUrl"`
}

type StopTimes struct {
	TripID        string `csv:"trip_id" json:"tripId"`
	ArrivalTime   string `csv:"arrival_time" json:"arrivalTime"`
	DepartureTime string `csv:"departure_time" json:"departureTime"`
	StopID        string `csv:"stop_id" json:"stopId"`
	StopSequence  int    `csv:"stop_sequence" json:"stopSequence"`
	PickupType    string `csv:"pickup_type" json:"pickupType"`
	DropOffType   string `csv:"drop_off_type" json:"dropOffType"`
}

type Stop struct {
	StopID        string `csv:"stop_id" json:"stopId"`
	StopName      string `csv:"stop_name" json:"stopName"`
	StopNameLower string `csv:"-" json:"-"`
	StopLatitude  string `csv:"stop_lat" json:"stopLat"`
	StopLongitude string `csv:"stop_lon" json:"stopLon"`
	LocationType  string `csv:"location_type" json:"locationType"`
}

type Trips struct {
	RouteID       string `csv:"route_id" json:"routeId"`
	ServiceID     string `csv:"service_id" json:"serviceId"`
	TripID        string `csv:"trip_id" json:"tripId"`
	TripHeadsign  string `csv:"trip_headsign" json:"tripHeadsign"`
	TripShortName string `csv:"trip_short_name" json:"tripShortName"`
}
