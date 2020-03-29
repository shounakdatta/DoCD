package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

type Command struct {
	Directory   string
	Command     string
	Environment []string
}

type Service struct {
	Type                 string
	Path                 string
	PackageManager       string
	InstallationCommands []Command
	BuildCommands        []Command
}

type Config struct {
	BasePackageManager string
	Services           []Service
}

// Global variables
var cmdSlice []*exec.Cmd

// signalChan := make(chan os.Signal, 1)
// exitChan := make(chan int)

func signalHandler(signalChan chan os.Signal, exitChan chan int) {
	signal.Notify(signalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		for {
			s := <-signalChan
			switch s {
			// kill -SIGHUP XXXX
			case syscall.SIGHUP:
				exitChan <- 0

			// kill -SIGINT XXXX or Ctrl+c
			case syscall.SIGINT:
				exitChan <- 1

			// kill -SIGTERM XXXX
			case syscall.SIGTERM:
				exitChan <- 0

			// kill -SIGQUIT XXXX
			case syscall.SIGQUIT:
				exitChan <- 0

			default:
				fmt.Println("Unknown signal.")
				exitChan <- 1
			}
		}
	}()
}

func main() {
	// Register signal handlers
	signalChan := make(chan os.Signal, 1)
	exitChan := make(chan int)
	signalHandler(signalChan, exitChan)

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

	// Initialize services
	for _, service := range config.Services {

		// Install service dependencies
		fmt.Println("Installing", service.Type, "dependencies...")
		for _, commandObj := range service.InstallationCommands {
			command := strings.Split(commandObj.Command, " ")
			cmd := exec.Command(command[0], command[1:]...)
			path := dir + commandObj.Directory
			cmd.Dir = path
			_, err := cmd.Output()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}
		fmt.Println("All", service.Type, "dependencies installed in", service.Path)

		// Build service
		for _, commandObj := range service.BuildCommands {
			command := strings.Split(commandObj.Command, " ")
			cmd := exec.Command(command[0], command[1:]...)
			cmd.Env = os.Environ()
			cmd.Env = append(cmd.Env, commandObj.Environment...)
			path := dir + commandObj.Directory
			cmd.Dir = path
			err := cmd.Start()
			cmdSlice = append(cmdSlice, cmd)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}
	}
	fmt.Println("All services started")

	code := <-exitChan
	fmt.Println("Terminating services...")
	for _, cmd := range cmdSlice {
		cmd.Process.Kill()
	}
	os.Exit(code)
	return
}
