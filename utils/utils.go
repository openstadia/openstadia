package utils

import (
	"os"
	"path/filepath"
)

const AppName = "Openstadia"

func GetConfigDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	appConfigDir := filepath.Join(configDir, AppName)

	err = os.Mkdir(appConfigDir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return "", err
	}

	return appConfigDir, nil
}
