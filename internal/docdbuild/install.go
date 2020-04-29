package docdbuild

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/shounakdatta/DoCD/internal/docdtypes"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"strings"
)

// InstallCmd : Installs services and service dependencies
func InstallCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "install",
		Short: "Installs services and service dependencies",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get config file
			config := docdtypes.ReadConfig()

			// Check if Admin
			isAdmin, _ := checkAdmin()
			if !isAdmin {
				color.Yellow(
					"Warning: You are not running DoCD as an administrator.\n" +
						"Service installations will be skipped.")
				installServices = false
			}

			installServicesAndDependencies(config)

			return nil
		},
	}
}

func installServicesAndDependencies(config docdtypes.Config) {
	// Get working directory
	dir, _ := os.Getwd()
	for _, service := range config.Services {
		installService(service, config.BasePackageManager)
		installServiceDependencies(service, dir)
	}
}

func installService(service docdtypes.Service, bpm string) {
	if installServices {
		serviceCmd := exec.Command(bpm, "install", service.ServiceName, "--confirm")
		serviceCmd.Stdout = os.Stdout
		serviceErr := serviceCmd.Run()
		if serviceErr != nil {
			fmt.Println(serviceErr.Error())
			os.Exit(1)
		}
	}
}

func installServiceDependencies(service docdtypes.Service, dir string) {
	fmt.Println("Installing", service.ServiceName, "dependencies...")
	for _, commandObj := range service.InstallationCommands {
		command := strings.Split(commandObj.Command, " ")
		cmd := exec.Command(command[0], command[1:]...)
		path := dir + commandObj.Directory
		cmd.Dir = path
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}
	fmt.Println("All", service.ServiceName, "dependencies installed in", service.Path)
}
