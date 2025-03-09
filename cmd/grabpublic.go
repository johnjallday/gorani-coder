package cmd

import (
	"agent/gorani/internal/grab"
	"fmt"

	"github.com/spf13/cobra"
)

var grabPublicCmd = &cobra.Command{
	Use:   "grab-public",
	Short: "Prints public functions and their descriptions",
	RunE: func(cmd *cobra.Command, args []string) error {
		root := "./internal" // Adjust path as needed
		functions, err := grab.GrabPublicFuncsWithDescriptions(root)
		if err != nil {
			return err
		}

		fmt.Println("Public Functions and Descriptions:")
		for name, desc := range functions {
			fmt.Printf("- %s: %s\n", name, desc)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(grabPublicCmd)
}
