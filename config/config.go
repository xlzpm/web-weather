package config

type ServerConfiguration struct {
	API         string `json:"apiUrl"`
	APIKeyParam string `json:"apiKeyParam"`
	ServerPort  int    `json:"serverPort"`
	APIKey      string
}
