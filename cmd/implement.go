package cmd

import (
	"agent/gorani/internal/implement"
	"fmt"

	"github.com/spf13/cobra"
)

var implementCmd = &cobra.Command{
	Use:   "implement <create|merge|prepare> [branchName]",
	Short: "Manages Git branches (create, merge) or prepares implementation prompt",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		action := args[0]
		switch action {
		case "create":
			if len(args) < 2 {
				return fmt.Errorf("Usage: implement create <branchName>")
			}
			branchName := args[1]
			return implement.CreateGitBranch(branchName)
		case "merge":
			if len(args) < 2 {
				return fmt.Errorf("Usage: implement merge <branchName>")
			}
			branchName := args[1]
			return implement.MergeBranch(branchName)
		case "prepare":
			return implement.PrepareImplementPrompt()
		case "prompt":
			return implement.Implement()
		default:
			return fmt.Errorf("Unknown action: %s", action)
		}
	},
}

func init() {
	rootCmd.AddCommand(implementCmd)
}
