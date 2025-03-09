package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use: "BuildReadme",
		Short: "and the logo image from docs/logo.png and combines them into README.md.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Action for BuildReadme is executed")
		},
	})

}
