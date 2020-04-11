package docdinit

import (
	"encoding/json"
	"fmt"
	"github.com/shounakdatta/DoCD/internal/docdtypes"
	"github.com/spf13/cobra"
	"gopkg.in/abiosoft/ishell.v2"
	"io/ioutil"
	"os"
)

var (
	shell = ishell.New()
)

func init() {
	shell.AddCmd(generate())
}

// InitCmd : Initializes DoCD iShell that generates DoCD-config.json
func InitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initializes working directory with DoCD",
		RunE: func(cmd *cobra.Command, args []string) error {
			// display welcome info.
			shell.Println("Initializing DoCD...")

			shell.Run()
			return nil
		},
	}
}

func generate() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "generate",
		Help: "Generates DoCD-config.json",
		Func: func(c *ishell.Context) {
			wd, err := os.Getwd()
			if err != nil {
				c.Println(err.Error())
				return
			}

			c.ShowPrompt(false)
			defer c.ShowPrompt(true)

			var configFile docdtypes.Config

			c.Print("Package Name: ")
			configFile.ProjectName = c.ReadLine()

			choices := []string{"choco", "brew", "apt-get"}
			choice := c.MultiChoice([]string{
				"choco (Windows)",
				"brew (macOS)",
				"apt-get (Linux)",
			}, "Base Package Manager: ")
			configFile.BasePackageManager = choices[choice]
			configFile.InstallServices = true

			c.Print("Add a service? (Y/n)")
			addService := c.ReadLine()
			for addService != "n" {
				var newService docdtypes.Service

				// Get service name
				serviceName := "python"
				c.Print(fmt.Sprintf("Service Name (%s): ", serviceName))
				newServiceName := c.ReadLine()
				if newServiceName != "" {
					serviceName = newServiceName
				}
				newService.ServiceName = serviceName

				// Get package manager
				packageManager := "pip"
				c.Print(fmt.Sprintf("Package Manager (%s): ", packageManager))
				newPackageManager := c.ReadLine()
				if newPackageManager != "" {
					packageManager = newPackageManager
				}
				newService.PackageManager = packageManager

				// Get service path
				servicePath := "./server"
				c.Print(fmt.Sprintf("Path (%s): ", servicePath))
				newServicePath := c.ReadLine()
				if newServicePath != "" {
					servicePath = newServicePath
				}
				newService.Path = servicePath

				logFile := fmt.Sprintf("./logs/%s-service.log", serviceName)
				c.Print(fmt.Sprintf("Log File Path (%s): ", logFile))
				newLogFile := c.ReadLine()
				if newLogFile != "" {
					logFile = newLogFile
				}
				newService.LogFilePath = logFile

				newService.InstallationCommands = []docdtypes.Command{}
				newService.BuildCommands = []docdtypes.Command{}

				configFile.Services = append(configFile.Services, newService)

				c.Print("Add another service? (Y/n): ")
				addService = c.ReadLine()
			}

			outputText := "\n\nDoCD initialization complete.\n" +
				"Enter service installation and build " +
				"commands in the generated DoCD-config.json file."
			c.Println(outputText)

			file, _ := json.MarshalIndent(configFile, "", "	")
			_ = ioutil.WriteFile(wd+"/DoCD-config.json", file, 0644)

		},
	}
}
