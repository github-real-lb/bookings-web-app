package config

import (
	"encoding/json"
	"os"
)

type DBConfig struct {
	// DBConnectionString is the connection string for the database
	DBConnectionString string `json:"db_connection_string"`

	// ServerAddress is http.Server listening address
	ServerAddress string `json:"server_address"`

	// TestDBConnectionString is the connection string for the test database
	TestDBConnectionString string `json:"test_db_connection_string"`
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
