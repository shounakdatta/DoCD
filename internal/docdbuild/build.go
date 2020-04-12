package docdbuild

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/shounakdatta/DoCD/internal/docdtypes"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

// Global variables
var (
	cmdSlice        []cmdReference
	installServices bool = true
)

// BuildCmd : Installs dependencies and builds all services
func BuildCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "build",
		Short: "Installs dependencies and builds all services",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Register signal handlers
			signalChan := make(chan os.Signal, 1)
			exitChan := make(chan int)
			signalHandler(signalChan, exitChan)

			// Get config file
			config := ReadConfig()

			// Get working directory
			dir, _ := os.Getwd()

			// Make log directory
			os.MkdirAll(dir+"\\logs", os.ModePerm)

			// Check if Admin
			isAdmin, _ := checkAdmin()
			if !isAdmin {
				color.Yellow(
					"Warning: You are not running DoCD as an administrator.\n" +
						"Service installations will be skipped.")
				installServices = false
			}

			// Initialize services
			InitializeServices(config)
			color.Cyan("To terminate session, press CTRL+C")

			// Initialize webhook
			http.HandleFunc("/github-push-master", autoDeploy)
			go http.ListenAndServe(":6000", nil)

			// Wait for exit signal
			_ = <-exitChan

			// Kill all services in their respective terminals
			fmt.Println("Terminating services...")
			for _, cmdRef := range cmdSlice {
				cmdRef.Cmd.Process.Kill()
				cmdRef.Cmd.Process.Wait()
				cmdRef.LogFile.Close()
			}

			return nil
		},
	}
}

// ReadConfig : Reads the DoCD configuration file
func ReadConfig() docdtypes.Config {
	var config docdtypes.Config
	configFile, err := os.Open(docdtypes.ConfigFileName)

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
func InitializeServices(config docdtypes.Config) {
	// Get working directory
	dir, _ := os.Getwd()
	for _, service := range config.Services {
		// Create log file
		logFile, err := os.Create(service.LogFilePath)
		if err != nil {
			panic(err)
		}

		// Install service
		if installServices {
			serviceCmd := exec.Command(config.BasePackageManager, "install", service.ServiceName)
			serviceCmd.Stdout = os.Stdout
			serviceErr := serviceCmd.Run()
			if serviceErr != nil {
				fmt.Println(serviceErr.Error())
				os.Exit(1)
			}
		}

		// Install service dependencies
		fmt.Println("Installing", service.ServiceName, "dependencies...")
		for _, commandObj := range service.InstallationCommands {
			command := strings.Split(commandObj.Command, " ")
			cmd := exec.Command(command[0], command[1:]...)
			path := dir + commandObj.Directory
			cmd.Dir = path
			cmd.Stdout = logFile
			err := cmd.Run()
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		}
		fmt.Println("All", service.ServiceName, "dependencies installed in", service.Path)

		// Build service
		for _, commandObj := range service.BuildCommands {
			command := strings.Split(commandObj.Command, " ")
			cmd := exec.Command(command[0], command[1:]...)
			cmd.Env = os.Environ()
			cmd.Env = append(cmd.Env, commandObj.Environment...)
			path := dir + commandObj.Directory
			cmd.Dir = path
			cmd.Stdout = logFile
			cmd.Stderr = logFile
			err := cmd.Start()
			cmdSlice = append(cmdSlice, cmdReference{cmd, logFile})
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		}
	}
	color.Cyan("All services started")
}
