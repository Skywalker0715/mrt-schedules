package station

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Skywalker0715/mrt-schedules/common/client"
)

type Service interface {
	GetAllStation() (response []StationResponse, err error)
	CheckSchedulesByStation(id string) (response []ScheduleResponse, err error)
}

type service struct {
	client *http.Client
}

func NewService() Service {
	return &service{
		client: &http.Client{},
	}
}

func (s *service) GetAllStation() (response []StationResponse, err error) {
	url := "https://www.jakartamrt.co.id/id/val/stasiuns"

	byteResponse, err := client.DoRequest(s.client, url)
	if err != nil {
		return
	}

	log.Printf("Fetching all stations")
	var stations []Station
	err = json.Unmarshal(byteResponse, &stations)
	if err != nil {
		log.Printf("Error unmarshaling stations: %v", err)
		return
	}

	log.Printf("Found %d stations", len(stations))

	for _, item := range stations {
		log.Printf("Processing station: ID=%s, Name=%s", item.ID, item.Name)
		stationResponse := ConvertStationToResponse(item)
		response = append(response, stationResponse)
	}

	return
}

func (s *service) CheckSchedulesByStation(id string) (response []ScheduleResponse, err error) {
	// 🔧 Clean the ID from any whitespace
	cleanId := strings.TrimSpace(id)
	log.Printf("=== DEBUG: Raw ID: '%s' (len=%d) ===", id, len(id))
	log.Printf("=== DEBUG: Clean ID: '%s' (len=%d) ===", cleanId, len(cleanId))

	// Use the cleaned ID for the rest of the function
	id = cleanId

	url := "https://www.jakartamrt.co.id/id/val/stasiuns"
	byteResponse, err := client.DoRequest(s.client, url)
	if err != nil {
		log.Printf("ERROR: Failed to fetch data from API: %v", err)
		return
	}

	log.Printf("DEBUG: Raw API response length: %d bytes", len(byteResponse))

	// Debug: Print first 500 characters of response
	responseStr := string(byteResponse)
	if len(responseStr) > 500 {
		log.Printf("DEBUG: Raw API response preview: %s...", responseStr[:500])
	} else {
		log.Printf("DEBUG: Full API response: %s", responseStr)
	}

	// Try to unmarshal as array first
	var schedules []Schedule
	err = json.Unmarshal(byteResponse, &schedules)
	if err != nil {
		log.Printf("ERROR: Failed to unmarshal as array: %v", err)

		// Try to unmarshal as single object
		var singleSchedule Schedule
		err2 := json.Unmarshal(byteResponse, &singleSchedule)
		if err2 != nil {
			log.Printf("ERROR: Failed to unmarshal as single object too: %v", err2)
			return nil, fmt.Errorf("failed to parse JSON response: %v", err)
		}
		schedules = append(schedules, singleSchedule)
	}

	log.Printf("DEBUG: Successfully parsed %d schedules", len(schedules))

	// Debug: Print all available stations with their IDs
	log.Printf("DEBUG: Available stations:")
	for i, item := range schedules {
		log.Printf("  [%d] ID='%s', Name='%s'", i, item.ID, item.Name)
	}

	// Find schedule by station ID
	var scheduleSelected Schedule
	found := false
	log.Printf("DEBUG: Searching for station with ID='%s'", id)

	for i, item := range schedules {
		log.Printf("DEBUG: [%d] Comparing '%s' == '%s' -> %t", i, item.ID, id, item.ID == id)
		if item.ID == id {
			scheduleSelected = item
			found = true
			log.Printf("DEBUG: Found matching station: ID='%s', Name='%s'", item.ID, item.Name)
			log.Printf("DEBUG: Lebak Bulus schedule: '%s'", item.ScheduleLebakBulus)
			log.Printf("DEBUG: Bundaran HI schedule: '%s'", item.ScheduleBundaranHI)
			break
		}
	}

	if !found {
		log.Printf("ERROR: Station with ID '%s' not found", id)
		log.Printf("DEBUG: Available IDs are: %v", func() []string {
			var ids []string
			for _, s := range schedules {
				ids = append(ids, s.ID)
			}
			return ids
		}())
		err = errors.New("station not found")
		return
	}

	scheduleResponse, err := ConvertDataToResponse(scheduleSelected)
	if err != nil {
		log.Printf("ERROR: Failed to convert schedule data: %v", err)
		return
	}

	// Convert to final response format
	for _, item := range scheduleResponse {
		response = append(response, ConvertScheduleToResponse(item))
	}

	log.Printf("DEBUG: Returning %d schedule entries", len(response))
	return
}

func ConvertDataToResponse(schedule Schedule) (response []ScheduleResponse, err error) {
	var (
		LebakBulusTripName = "Stasiun Lebak Bulus Grab"
		BundaranHITripName = "Stasiun Bundaran HI Bank DKI"
	)

	log.Printf("DEBUG: Converting schedule data for station: %s", schedule.Name)
	log.Printf("DEBUG: Lebak Bulus schedule length: %d chars", len(schedule.ScheduleLebakBulus))
	log.Printf("DEBUG: Bundaran HI schedule length: %d chars", len(schedule.ScheduleBundaranHI))

	// Handle empty schedules
	if schedule.ScheduleLebakBulus == "" && schedule.ScheduleBundaranHI == "" {
		log.Printf("WARNING: Both schedules are empty for station %s", schedule.Name)
		return response, nil
	}

	currentTime := time.Now().Format("15:04")
	log.Printf("DEBUG: Current time: %s", currentTime)

	// Process Lebak Bulus schedule
	if schedule.ScheduleLebakBulus != "" {
		log.Printf("DEBUG: Processing Lebak Bulus schedule...")
		scheduleLebakBulusParsed, err := ConvertScheduleToTimeFormat(schedule.ScheduleLebakBulus)
		if err != nil {
			log.Printf("ERROR: Failed to parse Lebak Bulus schedule: %v", err)
		} else {
			log.Printf("DEBUG: Parsed %d Lebak Bulus time entries", len(scheduleLebakBulusParsed))
			count := 0
			for _, item := range scheduleLebakBulusParsed {
				timeStr := item.Format("15:04")
				if timeStr > currentTime {
					response = append(response, ScheduleResponse{
						StationName: LebakBulusTripName,
						Time:        timeStr,
					})
					count++
				}
			}
			log.Printf("DEBUG: Added %d future Lebak Bulus schedules", count)
		}
	}

	// Process Bundaran HI schedule
	if schedule.ScheduleBundaranHI != "" {
		log.Printf("DEBUG: Processing Bundaran HI schedule...")
		scheduleBundaranHIParsed, err := ConvertScheduleToTimeFormat(schedule.ScheduleBundaranHI)
		if err != nil {
			log.Printf("ERROR: Failed to parse Bundaran HI schedule: %v", err)
		} else {
			log.Printf("DEBUG: Parsed %d Bundaran HI time entries", len(scheduleBundaranHIParsed))
			count := 0
			for _, item := range scheduleBundaranHIParsed {
				timeStr := item.Format("15:04")
				if timeStr > currentTime {
					response = append(response, ScheduleResponse{
						StationName: BundaranHITripName,
						Time:        timeStr,
					})
					count++
				}
			}
			log.Printf("DEBUG: Added %d future Bundaran HI schedules", count)
		}
	}

	log.Printf("DEBUG: Total schedule entries created: %d", len(response))
	return response, nil
}

func ConvertScheduleToTimeFormat(schedule string) (response []time.Time, err error) {
	if schedule == "" {
		return response, nil
	}

	schedules := strings.Split(schedule, ",")
	log.Printf("DEBUG: Splitting schedule into %d parts", len(schedules))

	validCount := 0
	for i, item := range schedules {
		trimmedTime := strings.TrimSpace(item)
		if trimmedTime == "" {
			continue
		}

		parsedTime, parseErr := time.Parse("15:04", trimmedTime)
		if parseErr != nil {
			log.Printf("WARNING: Invalid time format '%s' at index %d: %v", trimmedTime, i, parseErr)
			continue // Skip invalid time entries
		}

		response = append(response, parsedTime)
		validCount++
	}

	log.Printf("DEBUG: Successfully parsed %d/%d valid time entries", validCount, len(schedules))
	return response, nil
}

func ConvertStationToResponse(station Station) StationResponse {
	return StationResponse{
		ID:   station.ID,
		Name: station.Name,
	}
}

func ConvertScheduleToResponse(schedule ScheduleResponse) ScheduleResponse {
	return ScheduleResponse{
		StationName: schedule.StationName,
		Time:        schedule.Time,
	}
}
