package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/pkg/errors"

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
		report, err := wc.getSingleCityWeather(city)
		if err != nil {
			return nil, err
		}

		result.Reports[city] = report
	}

	return &result, nil
}

func (wc *SyncWeatherService) getSingleCityWeather(city string) (*CityReport, error) {
	log.Debugf("Fetching weather for city %s", city)

	fullUrl := wc.buildUrl(city)

	log.Info(fullUrl)

	req, reqErr := http.NewRequest("GET", fullUrl, nil)
	if reqErr != nil {
		return nil, reqErr
	}

	resp, callErr := wc.Client.Do(req)
	if callErr != nil {
		return nil, callErr
	}

	if resp.StatusCode == http.StatusOK {
		return createReport(resp)
	} else if resp.StatusCode == http.StatusNotFound {
		log.Warnf("City '%s' not found. No weather data returned.", city)
		return &CityReport{
			Desсription: "not found",
		}, nil
	}

	return nil, errors.Errorf("Weather API error: %d", resp.StatusCode)
}

func (wc *SyncWeatherService) buildUrl(city string) string {
	return fmt.Sprintf("%s%s&%s=%s", wc.config.API, city, wc.config.APIKeyParam, wc.config.APIKey)
}

func createReport(resp *http.Response) (*CityReport, error) {
	bytes, ioErr := io.ReadAll(resp.Body)
	if ioErr != nil {
		return nil, ioErr
	}

	var result map[string]interface{}
	unmErr := json.Unmarshal(bytes, &result)
	if unmErr != nil {
		return nil, unmErr
	}

	weatherData := result["weather"].([]interface{})[0].(map[string]interface{})
	mainData := result["main"].(map[string]interface{})

	report := &CityReport{
		Desсription: weatherData["description"].(string),
		Temperature: formatTemperature(mainData["temp"]),
	}

	return report, nil
}

func formatTemperature(temp interface{}) float64 {
	num := temp.(float64)

	val, err := strconv.ParseFloat(fmt.Sprintf("%.0f", num-273.0), 64)
	if err != nil {
		log.Errorf("Error converting temperature string [%f]: %s", num, err)
		return -1.0
	}

	return val
}
