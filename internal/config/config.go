package config

import (
	"encoding/json"
	"os"
)

const configFileName = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", nil
	}
	return dir + "/" + configFileName, nil
}

func write(conf Config) error {
	filePath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	updatedJSON, err := json.Marshal(conf)
	if err != nil {
		return err
	}
	if err := os.WriteFile(filePath, updatedJSON, 0644); err != nil {
		return err
	}
	return nil
}
