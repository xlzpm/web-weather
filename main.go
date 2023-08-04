package main

import (
	"encoding/json"
	"fmt"

	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/xlzpm/web-weather/api"
	"github.com/xlzpm/web-weather/config"
	"github.com/xlzpm/web-weather/service"
)

const (
	environment = "ENVIRONMENT"
	apiKey      = "API_KEY"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.JSONFormatter{})

	serverConfig, confErr := loadConfiguration()

	if confErr != nil {
		log.Error(errors.Wrap(confErr, "Cannot start server. Configuration error"))
		return
	}

	log.Debug("Starting weather server...")
	router := mux.NewRouter().StrictSlash(true)

	weatherApi := api.CreateWeatherEndpoint(service.CreateWeatherService(serverConfig, http.DefaultClient))

	router.HandleFunc("/weather", weatherApi.WeatherEndpoint).Methods(http.MethodPost)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", serverConfig.ServerPort), router))
}

func loadConfiguration() (*config.ServerConfiguration, error) {
	env := os.Getenv(environment)

	var configFile = ""

	switch env {
	case "dev":
		configFile = "configFile/config.dev.json"
	case "test":
		configFile = "configFile/config.test.json"
	case "prod":
		configFile = "configFile/config.prod.json"
	default:
		return nil, errors.Errorf("Unknown environment set - '%s'", env)
	}

	log.Infof("loading configuration from '%s'", configFile)
	file, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}

	conf := config.ServerConfiguration{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&conf)
	if err != nil {
		return nil, err
	}

	key := os.Getenv(apiKey)
	if len(key) == 0 {
		return nil, errors.Errorf("Bad API key set - '%s'", key)
	}

	conf.APIKey = key
	return &conf, nil
}
