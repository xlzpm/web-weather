package service

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/xlzpm/web-weather/config"
)

type WeatherService interface {
	GetWeather(cities []string) (*WeatherReports, error)
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type SyncWeatherService struct {
	config *config.ServerConfiguration
	Client HttpClient
}

func CreateWeatherService(config *config.ServerConfiguration, client HttpClient) WeatherService {
	ws := &SyncWeatherService{config: config, Client: client}
	return ws
}

func (wc *SyncWeatherService) GetWeather(cities []string) (*WeatherReports, error) {
	result := WeatherReports{
		Reports: make(map[string]*CityReport),
	}

	for _, city := range cities {
		report, err := ws.GetSingleCityWeather(city)

	}
}

func (wc *SyncWeatherService) GetSingleCityWeather(city string) (*CityReport, error) {
	log.Debugf("Fetching weather for city %s", city)

	fullUrl := wc.buildUrl(city)

	log.Info(fullUrl)
}

func (wc *SyncWeatherService) buildUrl(city string) string {
	return fmt.Sprintf("%s%s&%s=%s", wc.config.API, city, wc.config.APIKeyParam, wc.config.APIKey)
}
