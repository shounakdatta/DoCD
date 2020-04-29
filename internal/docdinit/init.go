package docdinit

import (
	"fmt"
	"github.com/shounakdatta/DoCD/internal/docdtypes"
	"github.com/spf13/cobra"
	"gopkg.in/abiosoft/ishell.v2"
)

var (
	shell = ishell.New()
)

func init() {
	shell.AddCmd(generate())
	shell.AddCmd(addServiceCmd())
	shell.Interrupt(interruptHandler)
}

// InitCmd : Initializes DoCD iShell that generates DoCD-config.json
func InitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initializes working directory with DoCD",
		RunE: func(cmd *cobra.Command, args []string) error {
			// display welcome info.
			shell.Println("Initializing configuration generator...")
			shell.Println("Enter `help` for a list of commands")
			shell.Run()
			return nil
		},
	}
}

func interruptHandler(c *ishell.Context, count int, str string) {
	fmt.Println(count, str)
	shell.Stop()
}

func generate() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "generate",
		Help: fmt.Sprintf("Generates %s", docdtypes.ConfigFileName),
		Func: func(c *ishell.Context) {
			if docdtypes.CheckConfigExists() {
				fmt.Println("A configuration file already exists.")
				fmt.Println("Would you like to overwrite it? (y/N)")
				overwriteConfig := c.ReadLine()
				if overwriteConfig != "y" {
					fmt.Println("Cancelling generator...")
					fmt.Println("Enter `exit` or `Ctrl-c` to finish.")
					return
				}
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

			c.Print("Add a service? (Y/n)")
			addService := c.ReadLine()
			for addService != "n" {
				newService := addNewService(c)
				configFile.Services = append(configFile.Services, newService)

				c.Print("Add another service? (Y/n): ")
				addService = c.ReadLine()
			}

			outputText := "\n\nDoCD initialization complete.\n" +
				"Enter service installation and build " +
				fmt.Sprintf("commands in the generated %s file.", docdtypes.ConfigFileName) +
				"\nEnter `exit` or `Ctrl-c` to finish."
			c.Println(outputText)

			err := docdtypes.WriteConfig(configFile)
			if err != nil {
				panic(err)
			}
		},
	}
}
