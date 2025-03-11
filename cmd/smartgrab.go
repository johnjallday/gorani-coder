package cmd

import (
	"agent/gorani/internal/grab"
	"fmt"

	"github.com/spf13/cobra"
)

var smartGrabCmd = &cobra.Command{
	Use:   "smartgrab [folder]",
	Short: "Generates a summary of Go symbols and sends it to OpenAI",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		folder := "./"
		if len(args) == 1 {
			folder = args[0]
		}
		err := grab.SmartGrab(folder)
		if err != nil {
			return fmt.Errorf("smart grab failed: %v", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(smartGrabCmd)
}
