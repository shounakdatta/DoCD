package docdinit

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/shounakdatta/DoCD/internal/docdtypes"
	"gopkg.in/abiosoft/ishell.v2"
	"os"
	"path/filepath"
	"strings"
)

var (
	cyan = color.New(color.FgCyan).SprintFunc()
	red  = color.New(color.FgRed).SprintFunc()
)

func addServiceCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "add-service",
		Help: "Add a new service to existing configuration file",
		Func: func(c *ishell.Context) {
			configFile := docdtypes.ReadConfig()
			addService := "Y"
			for addService != "n" {
				newService := addNewService(c)
				configFile.Services = append(configFile.Services, newService)
				c.Print("Add another service? (Y/n): ")
				addService, _ = checkInterrupt(c.ReadLine(), true)
			}
			docdtypes.WriteConfig(configFile)
			c.ShowPrompt(true)
		},
	}
}

func addNewService(c *ishell.Context) docdtypes.Service {
	var newService docdtypes.Service
	c.Println(cyan("Enter `exit` to stop utility"))

	// Get service name
	serviceName := "python"
	c.Print(fmt.Sprintf("Service Name* (%s): ", serviceName))
	newServiceName, _ := checkInterrupt(c.ReadLine(), true)
	if newServiceName != "" {
		serviceName = newServiceName
	}
	newService.ServiceName = serviceName

	// Get package manager
	packageManager := ""
	c.Print("Package Manager (i.e. %s): ")
	newPackageManager, _ := checkInterrupt(c.ReadLine(), true)
	if newPackageManager != "" {
		packageManager = newPackageManager
	}
	newService.PackageManager = packageManager

	// Get service path
	servicePath := ""
	c.Print("Path (i.e. ./server): ")
	newServicePath, _ := checkInterrupt(c.ReadLine(), true)
	if newServicePath != "" {
		servicePath = newServicePath
	}
	newService.Path = servicePath

	logFile := fmt.Sprintf("./logs/%s-service.log", serviceName)
	c.Print(fmt.Sprintf("Log File Path* (%s): ", logFile))
	newLogFile, _ := checkInterrupt(c.ReadLine(), true)
	if newLogFile != "" {
		logFile = newLogFile
	}
	newService.LogFilePath = logFile

	c.Print("Add installation commands? (Y/n)")
	addInstCommand, _ := checkInterrupt(c.ReadLine(), true)
	newService.InstallationCommands = []docdtypes.Command{}
	for addInstCommand != "n" {
		instCmd, err := addNewCommand(c)
		if err == nil {
			newService.InstallationCommands = append(newService.InstallationCommands, instCmd)
		}
		c.Print("Add another command? (Y/n): ")
		addInstCommand, _ = checkInterrupt(c.ReadLine(), true)
	}

	c.Print("Add build commands? (Y/n)")
	addBuildCommand, _ := checkInterrupt(c.ReadLine(), true)
	newService.BuildCommands = []docdtypes.Command{}
	for addBuildCommand != "n" {
		buildCmd, err := addNewCommand(c)
		if err == nil {
			newService.BuildCommands = append(newService.BuildCommands, buildCmd)
		}
		c.Print("Add another command? (Y/n): ")
		addBuildCommand, _ = checkInterrupt(c.ReadLine(), true)
	}

	return newService
}

func addNewCommand(c *ishell.Context) (docdtypes.Command, error) {
	newCmd := docdtypes.Command{
		Directory:   "\\",
		Command:     "",
		Environment: []string{},
	}
	cwd, _ := os.Getwd()
	c.SetPrompt(fmt.Sprintf("%s>", cwd))
	c.ShowPrompt(true)
	c.Println(cyan("Enter `exit` to cancel command entry"))

	for newCmd.Command == "" {
		command, exit := checkInterrupt(c.ReadLine(), false)
		if exit {
			c.SetPrompt(">>>")
			c.ShowPrompt(false)
			return newCmd, errors.New("Cancel Readline")
		}
		cmdChunks := strings.Split(command, " ")

		if cmdChunks[0] == "cd" {
			dirChange := strings.ReplaceAll(cmdChunks[1], "/", "\\")
			dirChange = strings.TrimRight(dirChange, "\\") + "\\"

			_, DNE := os.Stat(cwd + "\\" + dirChange)
			if DNE != nil {
				c.Println(red(fmt.Sprintf("Directory %s not found", cmdChunks[1])))
				continue
			}

			newCmd.Directory += dirChange
			cwd = filepath.Dir(cwd + "\\" + dirChange)
			c.SetPrompt(fmt.Sprintf("%s>", cwd))
		} else if cmdChunks[0] == "export" {
			newCmd.Environment = append(newCmd.Environment, cmdChunks[1:]...)
		} else {
			newCmd.Command = command
		}
	}

	c.SetPrompt(">>>")
	c.ShowPrompt(false)
	return newCmd, nil
}
