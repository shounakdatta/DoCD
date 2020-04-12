package cmd

import (
	"github.com/spf13/cobra"
	"time"
)

func printTimeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "curtime",
		Short: "Returns current time in ruby format",
		RunE: func(cmd *cobra.Command, args []string) error {
			now := time.Now()
			prettyTime := now.Format(time.RubyDate)
			cmd.Println("The current time is", prettyTime)
			return nil
		},
	}
}
