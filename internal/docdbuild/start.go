package docdbuild

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/shounakdatta/DoCD/internal/docdtypes"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

// StartCmd : Launches all services using the build commands in DoCD-config.json
func StartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Launches all services",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Register signal handlers
			signalChan := make(chan os.Signal, 1)
			exitChan := make(chan int)
			signalHandler(signalChan, exitChan)

			// Get config file
			config := docdtypes.ReadConfig()

			// Get working directory
			dir, _ := os.Getwd()

			// Make log directory
			os.MkdirAll(dir+"\\logs", os.ModePerm)

			// Start services
			startServices(config)

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

func startServices(config docdtypes.Config) {
	// Get working directory
	dir, _ := os.Getwd()

	for _, service := range config.Services {
		// Create log file
		logFile, err := os.Create(service.LogFilePath)
		if err != nil {
			panic(err)
		}
		startService(service, dir, logFile)
	}
	color.Cyan("All services started")
	color.Cyan("To terminate session, press CTRL+C")
}

func startService(service docdtypes.Service, dir string, logFile *os.File) {
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
