package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "gorani",
	Short: "Gorani is a CLI tool for code analysis and git branch management",
	Long: `Gorani provides multiple functionalities such as printing the directory tree,
grabbing code files, generating documentation, and managing Git branches.`,
	// If no subcommand is provided, show help.
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
