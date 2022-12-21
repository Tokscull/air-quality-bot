package waqi

import (
	"air-quality-bot/pkg/http"
	"fmt"
)

const (
	getByGeoUrl       = "https://api.waqi.info/feed/geo:%f;%f/?token=%s"
	getByCityUrl      = "https://api.waqi.info/feed/%s/?token=%s"
	getByStationIdUrl = "https://api.waqi.info/feed/%d/?token=%s"
)

type Service struct {
	token string
}

func NewService(token string) *Service {
	s := &Service{
		token: token,
	}
	return s
}

func (s *Service) GetByGeo(latitude float64, longitude float64) (*AirQuality, error) {
	url := fmt.Sprintf(getByGeoUrl, latitude, longitude, s.token)
	return get(url)
}

func (s *Service) GetByCity(city string) (*AirQuality, error) {
	url := fmt.Sprintf(getByCityUrl, city, s.token)
	return get(url)
}

func (s *Service) GetByStationId(stationId int) (*AirQuality, error) {
	url := fmt.Sprintf(getByStationIdUrl, stationId, s.token)
	return get(url)
}

func get(url string) (*AirQuality, error) {
	resp, err := http.Get[responseJSON](url)

	if err != nil {
		return nil, err
	}

	if resp.Status != "ok" {
		return nil, fmt.Errorf("waqi server error: %s", resp.Message)
	}

	return resp.ToAirQuality(), nil
}
