package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/xlzpm/web-weather/service"
)

type WeatherRequest struct {
	Cities []string
}

type SyncWeatherEndpoint struct {
	Service service.WeatherService
}

func CreateWeatherEndpoint(svc service.WeatherService) *SyncWeatherEndpoint {
	return &SyncWeatherEndpoint{Service: svc}
}

func (h *SyncWeatherEndpoint) WeatherEndpoint(w http.ResponseWriter, r *http.Request) {
	wr := WeatherRequest{}

	reqBody, readErr := io.ReadAll(r.Body)
	if readErr != nil {
		log.Error("Error reading client HTTP request", readErr)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "Invalid Request format. Could not parse request body")
		return
	}

	jsonErr := json.Unmarshal(reqBody, &wr.Cities)
	if jsonErr != nil {
		log.Error("Error converting client weather request into json", jsonErr)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Invalid Request format. Check json format is an array of strings."))
		return
	}

	result, callErr := h.Service.GetWeather(wr.Cities)
	if callErr != nil {
		log.Error("Error invoking Weather API", callErr)
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte("Weather API error"))
		return
	}

	data, marshalErr := json.Marshal(result)
	if marshalErr != nil {
		log.Error("Error converting Weather API response into json", callErr)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Error marshalling JSON response"))
		return
	}

	log.Info("Successfully retreived weather request for client")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	_, _ = w.Write(data)
}
