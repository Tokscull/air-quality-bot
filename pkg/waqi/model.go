package waqi

import (
	"time"
)

type PollutionLvl int

const (
	PollutionLvlGood PollutionLvl = iota
	PollutionLvlModerate
	PollutionLvlSensUnhealthy
	PollutionLvlUnhealthy
	PollutionLvlVeryUnhealthy
	PollutionLvlHazardous
)

type AirQuality struct {
	Station Station
	Time    time.Time
	AQI     int16
	Level   PollutionLvl
}

type Station struct {
	ID        int
	Name      string
	URL       string
	Longitude float64
	Latitude  float64
}

func getPollutionLevel(index int16) PollutionLvl {
	if index < 51 {
		return PollutionLvlGood
	}
	if index < 101 {
		return PollutionLvlModerate
	}
	if index < 151 {
		return PollutionLvlSensUnhealthy
	}
	if index < 201 {
		return PollutionLvlUnhealthy
	}
	if index < 301 {
		return PollutionLvlVeryUnhealthy
	}
	return PollutionLvlHazardous
}
