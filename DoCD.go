package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type Command struct {
	Directory string
	Command   string
}

type Service struct {
	Type                 string
	Path                 string
	PackageManager       string
	InstallationCommands []Command
}

type Config struct {
	BasePackageManager string
	BuildFile          string
	Services           []Service
}

func main() {
	// Get config file
	configFile, err := os.Open("DoCD-config.json")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer configFile.Close()

	// Read config file
	byteValue, _ := ioutil.ReadAll(configFile)
	var config Config
	json.Unmarshal([]byte(byteValue), &config)

	// Get working directory
	dir, _ := os.Getwd()

	// Install service dependencies
	for _, service := range config.Services {
		for _, commandObj := range service.InstallationCommands {
			command := strings.Split(commandObj.Command, " ")
			cmd := exec.Command(command[0], command[1:]...)
			path := dir + commandObj.Directory
			cmd.Dir = path
			stdout, err := cmd.Output()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Println(path, err, string(stdout))
		}
	}

	fmt.Println("Installing dependencies...")
}
