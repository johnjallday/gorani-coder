package cmd

import (
	"agent/gorani/internal/grab"

	"github.com/spf13/cobra"
)

var summaryCmd = &cobra.Command{
	Use:   "summary [folder]",
	Short: "Grabs summary of Go symbols (functions, structs, interfaces) from the specified folder",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Default to current directory if no folder argument is provided.
		folder := "./"
		if len(args) == 1 {
			folder = args[0]
		}
		return grab.GrabSummary(folder)
	},
}

func init() {
	rootCmd.AddCommand(summaryCmd)
}
