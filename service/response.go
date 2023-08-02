package service

type WeatherReports struct {
	Reports map[string]*CityReport `json:"reports"`
}

type CityReport struct {
	Desсription string  `json:"description"`
	Temperature float64 `json:"temperature"`
}
