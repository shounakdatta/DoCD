package docdbuild

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/shounakdatta/DoCD/internal/docdtypes"
	"github.com/spf13/cobra"
	"gopkg.in/abiosoft/ishell.v2"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var (
	startshell = ishell.New()
	ciLogFile  *os.File
	exitChan   chan int
)

func init() {
	startshell.DeleteCmd("exit")
	startshell.AddCmd(terminateServicesCmd())
	startshell.AddCmd(enableCICmd())
	startshell.AddCmd(disableCICmd())
	startshell.Interrupt(interruptHandler)
}

// StartCmd : Launches all services using the build commands in DoCD-config.json
func StartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Launches all services",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get config file
			config := docdtypes.ReadConfig()

			// Get working directory
			dir, _ := os.Getwd()

			// Make log directory
			os.MkdirAll(dir+"\\logs", os.ModePerm)

			// Start services
			startAllServices(config)

			// Initialize webhook
			http.HandleFunc("/github-push-master", autoDeploy)
			go http.ListenAndServe(":6000", nil)

			startshell.Run()
			return nil
		},
	}
}

func startAllServices(config docdtypes.Config) {
	// Get working directory
	dir, _ := os.Getwd()
	ciLogFile, _ = os.Create("../logs/ci-service.log")

	for _, service := range config.Services {
		// Create log file
		logFile, err := os.Create(service.LogFilePath)
		if err != nil {
			panic(err)
		}
		startService(service, dir, logFile)
	}
	color.Cyan("All services started")
	color.Cyan("To terminate session, enter `terminate` or press CTRL+C")
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
		cmdMap[service.ServiceName] = append(cmdMap[service.ServiceName], len(cmdSlice)-1)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}
}

func terminateServicesCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name:    "terminate",
		Aliases: []string{"exit", "stop"},
		Help:    "Terminates all services",
		Func: func(c *ishell.Context) {
			startshell.Stop()
			terminateAllServices()
		},
	}
}

func terminateAllServices() {
	color.Cyan("Terminating services...")
	fmt.Println(cmdMap)
	for _, cmdRef := range cmdSlice {
		terminateService(cmdRef)
	}
	color.Cyan("All services terminated")
	os.Exit(0)
}

func terminateService(cmdRef cmdReference) {
	if cmdRef.LogFile != nil {
		cmdRef.LogFile.Close()
	}

	kill := exec.Command("taskkill", "/T", "/F", "/PID", strconv.Itoa(cmdRef.Cmd.Process.Pid))
	err := kill.Run()
	if err != nil {
		fmt.Println("Error killing process")
	}
	fmt.Println("A process was killed")
}

func enableCICmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "enable-ci",
		Help: "Enables continuous deployment",
		Func: func(c *ishell.Context) {
			command := strings.Split("ngrok http 6000", " ")
			cmd := exec.Command(command[0], command[1:]...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = ciLogFile
			err := cmd.Start()
			cmdSlice = append(cmdSlice, cmdReference{Cmd: cmd})
			cmdMap["ci"] = append(cmdMap["ci"], len(cmdSlice)-1)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		},
	}
}

func disableCICmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "disable-ci",
		Help: "Disables continuous deployment",
		Func: func(c *ishell.Context) {
			terminateService(cmdSlice[cmdMap["ci"][0]])
			cmdMap["ci"] = cmdMap["ci"][1:]
		},
	}
}

func interruptHandler(c *ishell.Context, count int, str string) {
	startshell.Stop()
	terminateAllServices()
}
