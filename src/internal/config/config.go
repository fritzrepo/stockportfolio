package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	TransactionFilePath string `json:"transactionFilePath"`
	DatabaseFilePath    string `json:"databaseFilePath"`
}

func LoadConfigFromJSON(filename string) (*Config, error) {
	// Datei öffnen
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// JSON dekodieren
	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
