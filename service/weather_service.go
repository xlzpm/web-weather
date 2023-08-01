package service

import (
	"net/http"

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
