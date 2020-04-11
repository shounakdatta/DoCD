package cmd

import (
	"github.com/spf13/cobra"
	"gopkg.in/abiosoft/ishell.v2"
	"strings"
)

var (
	shell = ishell.New()
)

func init() {
	// register a function for "greet" command.
	shell.AddCmd(&ishell.Cmd{
		Name: "greet",
		Help: "greet user",
		Func: func(c *ishell.Context) {
			c.Println("Hello", strings.Join(c.Args, " "))
		},
	})
}

// TestiShell : Generates DoCD-config.json file
func TestiShell() *cobra.Command {
	return &cobra.Command{
		Use:   "ishell",
		Short: "Launch a sample interactive shell",
		RunE: func(cmd *cobra.Command, args []string) error {
			// display welcome info.
			shell.Println("Sample Interactive Shell")

			shell.Run()
			return nil
		},
	}
}
