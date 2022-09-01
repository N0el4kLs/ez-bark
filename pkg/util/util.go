package util

import (
	"log"
	"os"
	"path/filepath"
)

// get global config directory
func configPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Println("Can not get use home directory")
	}

	config := filepath.Join(home, ".config", "ezbark")
	return config
}

// get global config.yml path
func globalConfigFile() string {
	path := configPath()
	file := filepath.Join(path, "config.yml")
	return file
}

// check the config file is exist not.
// return nil when exists
// return err when not exists
func checkGlobalConf() error {
	_, err := os.Stat(globalConfigFile())
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		return err
	}
	return err
}

// create config file at ~/.config/ezbark/config.yml
func createGlobalConf() error {
	err := os.MkdirAll(configPath(), os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func getCurrentPath() string {
	pwd, _ := os.Getwd()
	return pwd
}
