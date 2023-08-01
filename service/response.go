package service

type WeatherReports struct {
	Reports map[string]*CityReport `json:"reports"`
}

type CityReport struct {
	Desсription string `json:"description"`
	Temperature string `json:"temperature"`
}
