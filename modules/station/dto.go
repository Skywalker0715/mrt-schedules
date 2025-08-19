package station

type Station struct {
	ID   string `json:"nid"`
	Name string `json:"Title"`
}

type StationResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Schedule struct {
	StationId          string `json:"nid"`
	StationName        string `json:"Title"`
	ScheduleBundaranHI string `json:"jadwal_hi_biasa"`
	ScheduleLebakBulus string `json:"jadwal_lb_biasa"`
}

type ScheduleResponse struct {
	StationName string `json:"station_name"`
	Time        string `json:"time"`
}

type schedules struct {
	StationID           string `json:"nid"`
	StationName         string `json:"Title"`
	SchedulesBundaranHI string `json:"jadwal_hi_biasa"`
	SchedulesLebakBulus string `json:"jadwal_lb_biasa"`
}

type SchedulesResponse struct {
	StationName string `json:"station_name"`
	Time        string `json:"time"`
}
