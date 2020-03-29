package main

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

// Command : Structure of service installation and build commands
type Command struct {
	Directory   string
	Command     string
	Environment []string
}

// Service : Structure of configuration services
type Service struct {
	Type                 string
	Path                 string
	PackageManager       string
	LogFile              string
	InstallationCommands []Command
	BuildCommands        []Command
}

// Config : Structure of DOCD-config.json
type Config struct {
	BasePackageManager string
	Services           []Service
}

// Global variables
var cmdSlice []*exec.Cmd

const (
	// ConfigFile : Configuration file name
	configFile = "DoCD-config.json"
)

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
			case syscall.SIGHUP:
				exitChan <- 0

			case syscall.SIGINT:
				exitChan <- 0

			case syscall.SIGTERM:
				exitChan <- 0

			case syscall.SIGQUIT:
				exitChan <- 0

			default:
				fmt.Println("Unknown signal.")
				exitChan <- 1
			}
		}
	}()
}

func readConfig() Config {
	var config Config
	configFile, err := os.Open(configFile)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer configFile.Close()

	byteValue, _ := ioutil.ReadAll(configFile)
	json.Unmarshal([]byte(byteValue), &config)
	return config
}

func initializeServices(config Config) {
	// Get working directory
	dir, _ := os.Getwd()
	for _, service := range config.Services {
		// Create log file
		logFile, err := os.Create(service.LogFile)
		if err != nil {
			panic(err)
		}

		// Install service dependencies
		fmt.Println("Installing", service.Type, "dependencies...")
		for _, commandObj := range service.InstallationCommands {
			command := strings.Split(commandObj.Command, " ")
			cmd := exec.Command(command[0], command[1:]...)
			path := dir + commandObj.Directory
			cmd.Dir = path
			cmd.Stdout = logFile
			err := cmd.Start()
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			cmd.Wait()
		}
		fmt.Println("All", service.Type, "dependencies installed in", service.Path)

		// Build service
		for _, commandObj := range service.BuildCommands {
			command := strings.Split(commandObj.Command, " ")
			fmt.Println(command)
			cmd := exec.Command(command[0], command[1:]...)
			cmd.Env = os.Environ()
			cmd.Env = append(cmd.Env, commandObj.Environment...)
			path := dir + commandObj.Directory
			cmd.Dir = path
			cmd.Stdout = logFile
			err := cmd.Start()
			cmdSlice = append(cmdSlice, cmd)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		}

		logFile.Close()
	}
	color.Cyan("All services started")
}

func main() {
	// Register signal handlers
	signalChan := make(chan os.Signal, 1)
	exitChan := make(chan int)
	signalHandler(signalChan, exitChan)

	// Get config file
	config := readConfig()

	// Initialize services
	initializeServices(config)
	color.Cyan("To terminate session, press CTRL+C")

	// Wait for exit signal
	code := <-exitChan

	// Kill all services in their respective terminals
	fmt.Println("Terminating services...")
	for _, cmd := range cmdSlice {
		cmd.Process.Kill()
	}

	os.Exit(code)
}
