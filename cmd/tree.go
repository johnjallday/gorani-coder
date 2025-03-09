package cmd

import (
	"agent/gorani/internal/tree"
	"fmt"

	"github.com/spf13/cobra"
)

var treeCmd = &cobra.Command{
	Use:   "tree [path]",
	Short: "Prints the directory tree structure",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := "."
		if len(args) > 0 {
			path = args[0]
		}
		fmt.Println("Printing Directory Tree:")
		return tree.CopyTreeToClipboard(path)
	},
}

func init() {
	rootCmd.AddCommand(treeCmd)
}
