package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (cfg Config) SetUser(userName string) error {
	cfg.CurrentUserName = userName

	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	err = write(data)
	if err != nil {
		return err
	}

	return nil
}

const configFileName = ".gatorconfig.json"

func Read() (Config, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var configStruct Config
	err = json.Unmarshal(data, &configStruct)
	if err != nil {
		return Config{}, err
	}

	return configStruct, nil
}

func getConfigFilePath() (string, error) {
	path, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	path += "/" + configFileName
	return path, err
}

func write(data []byte) error {
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
