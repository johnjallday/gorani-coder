package cmd

import (
	"agent/gorani/internal/prompt"
	"fmt"

	"github.com/spf13/cobra"
)

var promptCmd = &cobra.Command{
	Use:   "prompt",
	Short: "Prompts OpenAI with user input",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Enter your prompt:")
		prompt.PromptFromNeovim()
	},
}

func init() {
	rootCmd.AddCommand(promptCmd)
}
