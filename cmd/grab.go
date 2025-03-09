package cmd

import (
	"agent/gorani/internal/grab"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var grabCmd = &cobra.Command{
	Use:   "grab [file or folder]",
	Short: "Grabs code files (file or folder auto-detected)",
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// No argument defaults to current directory
		if len(args) < 1 {
			return grab.Grab("./")
		}

		// One argument: let the grab package auto-detect file or folder.
		if len(args) == 1 {
			return grab.Grab(args[0])
		}

		// Multiple arguments: separate into directories and files.
		var dirs, files []string
		for _, arg := range args {
			info, err := os.Stat(arg)
			if err != nil {
				// If stat fails, assume it's a file (or adjust handling as needed)
				files = append(files, arg)
			} else if info.IsDir() {
				dirs = append(dirs, arg)
			} else {
				files = append(files, arg)
			}
		}

		// Use multifolder grab if all provided args are directories.
		if len(dirs) > 0 && len(files) == 0 {
			return grab.GrabMultipleFolders(dirs)
		}

		// Use multiple files grab if all provided args are files.
		if len(files) > 0 && len(dirs) == 0 {
			return grab.GrabFiles(files)
		}

		// If a mix of directories and files is provided, return an error.
		return fmt.Errorf("error: please provide either only files or only directories, not a mix")
	},
}

func init() {
	rootCmd.AddCommand(grabCmd)
}
