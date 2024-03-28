package config

import (
	"encoding/json"
	"os"
)

type DBConfig struct {
	// ConnectionString is the connection string for the postgreSql server
	ConnectionString string `json:"connection_string"`

	// ServerAddress is http.Server listening address
	ServerAddress string `json:"server_address"`
}

func LoadDBConfig(filename string) (*DBConfig, error) {
	dbConfig := DBConfig{}

	file, err := os.Open(filename)
	if err != nil {
		return &dbConfig, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&dbConfig)
	return &dbConfig, err
}
