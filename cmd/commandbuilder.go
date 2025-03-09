package cmd

import (
	"agent/gorani/internal/commandbuilder"

	"github.com/spf13/cobra"
)

// commandBuilderCmd defines a Cobra command that calls the internal commandbuilder.
var commandBuilderCmd = &cobra.Command{
	Use:   "commandbuilder",
	Short: "Interactively register public functions as commands",
	Long:  "Scans for public functions in the internal package and lets you select which functions to register as Cobra subcommands.",
	Run: func(cmd *cobra.Command, args []string) {
		commandbuilder.RegisterActions()
	},
}

func init() {
	rootCmd.AddCommand(commandBuilderCmd)
}
