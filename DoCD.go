package main

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"gopkg.in/go-playground/webhooks.v5/github"
	"io/ioutil"
	"net/http"
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

func main() {
	// Register signal handlers
	signalChan := make(chan os.Signal, 1)
	exitChan := make(chan int)
	SignalHandler(signalChan, exitChan)

	// Get config file
	config := ReadConfig()

	// Get working directory
	dir, _ := os.Getwd()

	// Make log directory
	os.MkdirAll(dir+"\\logs", os.ModePerm)

	// Initialize services
	InitializeServices(config)
	color.Cyan("To terminate session, press CTRL+C")

	// Initialize webhook
	http.HandleFunc("/github-push-master", DeployMaster)
	go http.ListenAndServe(":6000", nil)

	// Wait for exit signal
	code := <-exitChan

	// Kill all services in their respective terminals
	fmt.Println("Terminating services...")
	for _, cmd := range cmdSlice {
		cmd.Process.Kill()
	}

	os.Exit(code)
}

// SignalHandler : Handles all signals sent to DoCD
func SignalHandler(signalChan chan os.Signal, exitChan chan int) {
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

// ReadConfig : Reads the DoCD configuration file
func ReadConfig() Config {
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

// InitializeServices : Installs service dependecies and launches services
func InitializeServices(config Config) {
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

// DeployMaster : Pulls latest commit from remote master and deploys
func DeployMaster(res http.ResponseWriter, req *http.Request) {
	hook, _ := github.New(github.Options.Secret(""))
	// dir, _ := os.Getwd()

	payload, err := hook.Parse(req, github.PushEvent)
	if err != nil {
		if err == github.ErrEventNotFound {
			fmt.Println("Unknown event")
		}
	}

	switch payload.(type) {

	case github.PushPayload:
		push := payload.(github.PushPayload)
		fmt.Println("Change detected on", push.Ref, "- deploying...")
		cmd := exec.Command("git", "pull")
		err := cmd.Run()
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println("Deployment complete")
	}
	fmt.Fprintf(res, "Hello, %s!", req.URL.Path[1:])
}
