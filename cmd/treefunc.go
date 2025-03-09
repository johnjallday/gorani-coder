package cmd

import (
	"agent/gorani/internal/tree"
	"fmt"

	"github.com/spf13/cobra"
)

var treeFuncCmd = &cobra.Command{
	Use:   "tree-func [path]",
	Short: "Prints the directory tree structure with functions",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := "."
		if len(args) > 0 {
			path = args[0]
		}
		fmt.Println("Printing Directory Tree with Functions:")
		return tree.CopyTreeWithFunctionsToClipboard(path)
	},
}

func init() {
	rootCmd.AddCommand(treeFuncCmd)
}
