package cmd

import (
	"agent/gorani/internal/docbuilder"

	"github.com/spf13/cobra"
)

var docbuilderCmd = &cobra.Command{
	Use:   "docbuilder",
	Short: "Generates README.md using documentation builder",
	Run: func(cmd *cobra.Command, args []string) {
		docbuilder.BuildReadme()
	},
}

func init() {
	rootCmd.AddCommand(docbuilderCmd)
}
