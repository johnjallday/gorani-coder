package cmd

import (
	"agent/gorani/internal/grab"

	"github.com/spf13/cobra"
)

var grabCmd = &cobra.Command{
	Use:   "grab [file or folder]",
	Short: "Grabs code files (file or folder auto-detected)",
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return grab.Grab("./")
		} else if len(args) == 1 {
			return grab.Grab(args[0])
		} else {
			return grab.GrabFiles(args)
		}
	},
}

func init() {
	rootCmd.AddCommand(grabCmd)
}
