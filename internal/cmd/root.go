package cmd

import (
	"fmt"
	"github.com/shounakdatta/DoCD/internal/docdbuild"
	"github.com/shounakdatta/DoCD/internal/docdinit"
	"github.com/spf13/cobra"
	"os"
)

var (
	rootCmd = &cobra.Command{
		Use:   "docd",
		Short: "The CLI for docking and managing DoCD applications",
	}
)

func init() {
	// Add cmd package commands
	rootCmd.AddCommand(printTimeCmd())
	rootCmd.AddCommand(TestiShell())

	// Add docdinit package commands
	rootCmd.AddCommand(docdinit.InitCmd())

	// Add docdbuld package commands
	rootCmd.AddCommand(docdbuild.BuildCmd())
}

// Execute : Runs the root command
func Execute() error {
	err := rootCmd.Execute()
	if err != nil {
		er(err)
	}
	return err
}

func er(msg error) {
	fmt.Println("Error:", msg.Error())
	os.Exit(1)
}
