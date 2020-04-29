package docdtypes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// ReadConfig : Reads the DoCD configuration file
func ReadConfig() Config {
	var config Config
	configFile, err := os.Open(ConfigFileName)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer configFile.Close()

	byteValue, _ := ioutil.ReadAll(configFile)
	json.Unmarshal([]byte(byteValue), &config)
	return config
}

// WriteConfig : Writes to the DoCD configuration file
func WriteConfig(configFile Config) error {
	configFilePath, errA := GetConfigFilePath()
	file, errB := json.MarshalIndent(configFile, "", "	")
	if errA != nil {
		return errA
	}
	if errB != nil {
		return errB
	}
	return ioutil.WriteFile(configFilePath, file, 0644)
}

// GetConfigFilePath : Returns the filepath of the DoCD configuration file
func GetConfigFilePath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	return fmt.Sprintf(wd+"/%s", ConfigFileName), nil
}

// CheckConfigExists : Checks if DoCD configuration file already exists in directory
func CheckConfigExists() bool {
	configFilePath, _ := GetConfigFilePath()
	if _, err := os.Stat(configFilePath); err == nil {
		return true
	}
	return false
}
