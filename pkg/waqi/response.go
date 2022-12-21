package waqi

import "time"

type valueJSON struct {
	Value float32 `json:"v"`
}

type responseJSON struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		ID   int   `json:"idx"`
		AQI  int16 `json:"aqi"`
		Time struct {
			ISO *time.Time `json:"iso"`
		} `json:"time"`
		City struct {
			Name string    `json:"name"`
			URL  string    `json:"url"`
			Geo  []float64 `json:"geo"`
		} `json:"city"`
		IAQI struct {
			CO   valueJSON `json:"co"`
			H    valueJSON `json:"h"`
			NO2  valueJSON `json:"no2"`
			P    valueJSON `json:"p"`
			PM25 valueJSON `json:"pm25"`
			T    valueJSON `json:"t"`
			W    valueJSON `json:"w"`
			WG   valueJSON `json:"wg"`
		} `json:"iaqi"`
	} `json:"data"`
}

func (r responseJSON) ToAirQuality() *AirQuality {
	q := &AirQuality{
		Station: Station{
			ID:        r.Data.ID,
			Name:      r.Data.City.Name,
			URL:       r.Data.City.URL,
			Latitude:  r.Data.City.Geo[0],
			Longitude: r.Data.City.Geo[1],
		},
		Time:  r.Data.Time.ISO.UTC(),
		AQI:   r.Data.AQI,
		Level: getPollutionLevel(r.Data.AQI),
	}
	return q
}
