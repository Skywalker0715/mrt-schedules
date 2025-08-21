package station

type Station struct {
	ID   string `json:"nid"`
	Name string `json:"title"`
}

type StationResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Schedule struct {
	ID                 string `json:"nid"`             // ID stasiun
	Name               string `json:"title"`           // Nama stasiun
	ScheduleLebakBulus string `json:"jadwal_lb_biasa"` // Jadwal ke Lebak Bulus (hari biasa)
	ScheduleBundaranHI string `json:"jadwal_hi_biasa"` // Jadwal ke Bundaran HI (hari biasa)

	// Optional: Tambahan field jadwal lainnya jika diperlukan
	ScheduleLebakBulusLibur string `json:"jadwal_lb_libur"` // Jadwal ke Lebak Bulus (hari libur)
	ScheduleBundaranHILibur string `json:"jadwal_hi_libur"` // Jadwal ke Bundaran HI (hari libur)
}

type ScheduleResponse struct {
	StationName string `json:"station_name"`
	Time        string `json:"time"`
}
