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
