package grab

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/atotto/clipboard"
)

// Max files to allow grabbing before warning the user
const maxFilesLimit = 200

// Files that prevent grabbing when detected
var protectedFiles = []string{".config", "ws_info.toml"}

// Grab auto-detects if the input is a file, directory, or just a filename
func Grab(input string) error {
	// Prevent grabbing home or root directory
	homeDir, _ := os.UserHomeDir()
	absPath, _ := filepath.Abs(input)

	if absPath == "/" || absPath == homeDir {
		return fmt.Errorf("error: refusing to grab the entire root or home directory")
	}

	// Check if input is a valid file or directory
	info, err := os.Stat(input)
	if err == nil {
		// Check if it's a **protected workspace** before grabbing
		if isProtectedWorkspace(input) {
			return fmt.Errorf("error: cannot grab workspace or protected directory")
		}

		if info.IsDir() {
			fmt.Println("Checking directory size before grabbing...")
			fileCount := countFiles(input)

			if fileCount == -1 {
				return fmt.Errorf("error: unable to count files in directory")
			}

			// If the directory has too many files, ask for confirmation
			if fileCount > maxFilesLimit {
				fmt.Printf("⚠️ Warning: The directory '%s' contains %d files. Proceed? (y/N): ", input, fileCount)
				if !confirmAction() {
					return fmt.Errorf("aborted: too many files to grab")
				}
			}

			fmt.Println("Grabbing all code files in directory:", input)
			return GrabCodesProject(input)
		}

		fmt.Println("Grabbing single file:", input)
		return GrabCode(input)
	}

	// If not a direct file or folder, assume it's a filename to search for
	fmt.Println("Searching for file:", input)
	filePath, err := findFileByName(".", input)
	if err != nil {
		return fmt.Errorf("error: %s not found", input)
	}

	fmt.Println("Grabbing file:", filePath)
	return GrabCode(filePath)
}

// isProtectedWorkspace checks if a directory contains a protected file (e.g., .config or ws_info.toml)
func isProtectedWorkspace(path string) bool {
	for _, protected := range protectedFiles {
		protectedPath := filepath.Join(path, protected)
		if _, err := os.Stat(protectedPath); err == nil {
			if protected == "ws_info.toml" {
				fmt.Println("❌ Cannot grab workspace: ws_info.toml detected.")
			} else if protected == ".config" {
				fmt.Println("❌ Cannot grab root directory: .config detected.")
			}
			return true
		}
	}
	return false
}

// findFileByName searches for a file by name within the given directory (recursively)
func findFileByName(root string, filename string) (string, error) {
	var foundPath string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Base(path) == filename {
			foundPath = path
			return filepath.SkipDir // Stop searching once found
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	if foundPath == "" {
		return "", fmt.Errorf("file %s not found", filename)
	}

	return foundPath, nil
}

// GrabCode copies the content of a single file to the clipboard
func GrabCode(filePath string) error {
	// Read file contents
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file %s: %v", filePath, err)
	}

	// Format content for clipboard
	clipboardContent := fmt.Sprintf(">>> %s\n%s\n", filePath, string(content))

	// Copy to clipboard
	if err := clipboard.WriteAll(clipboardContent); err != nil {
		return fmt.Errorf("failed to copy file content to clipboard: %v", err)
	}

	fmt.Println("Copied content of", filePath, "to clipboard.")
	return nil
}

// GrabCodesProject copies all contents of found code files in a project directory to the clipboard.
func GrabCodesProject(root string) error {
	content, err := getCodesProjectContent(root)
	if err != nil {
		return err
	}

	if err := clipboard.WriteAll(content); err != nil {
		return fmt.Errorf("failed to copy to clipboard: %v", err)
	}
	fmt.Println("Copied all code files' contents to clipboard.")
	return nil
}

// getCodesProjectContent collects and formats code files' contents from a directory without writing to the clipboard.
func getCodesProjectContent(root string) (string, error) {
	var fileContents []string
	supportedExtensions := []string{
		".go", ".py", ".js", ".java", ".cpp", ".c", ".cs", ".rb", ".php", ".html",
		".css", ".sh",
	}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			ext := filepath.Ext(path)
			for _, supportedExt := range supportedExtensions {
				if ext == supportedExt {
					fmt.Println("Found code file:", path)

					// Read file contents
					content, readErr := os.ReadFile(path)
					if readErr != nil {
						fmt.Println("Error reading file:", path, readErr)
					} else {
						// Store formatted content
						fileContents = append(fileContents, fmt.Sprintf(">>> %s\n%s\n", path, string(content)))
					}
					// Once matched, no need to check further extensions
					break
				}
			}
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	if len(fileContents) == 0 {
		return "", fmt.Errorf("no code files found in %s", root)
	}

	return strings.Join(fileContents, "\n---\n"), nil
}

// countFiles counts the number of files in a directory (to prevent excessive grabbing)
func countFiles(root string) int {
	count := 0
	err := filepath.Walk(root, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			count++
		}
		return nil
	})
	if err != nil {
		return -1 // Error while counting files
	}
	return count
}

// confirmAction prompts the user for confirmation before proceeding
func confirmAction() bool {
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}

// GrabFiles accepts multiple file paths, reads their contents, and copies the combined content to the clipboard.
func GrabFiles(filePaths []string) error {
	var allContents []string

	for _, filePath := range filePaths {
		// Verify the file exists and is not a directory
		info, err := os.Stat(filePath)
		if err != nil {
			return fmt.Errorf("error: file %s not found", filePath)
		}
		if info.IsDir() {
			return fmt.Errorf("error: %s is a directory, not a file", filePath)
		}

		// Read file contents
		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("error reading file %s: %v", filePath, err)
		}

		// Format the content
		formatted := fmt.Sprintf(">>> %s\n%s\n", filePath, string(content))
		allContents = append(allContents, formatted)
	}

	// Join all file contents with a separator
	combinedContent := strings.Join(allContents, "\n---\n")

	// Copy to clipboard
	if err := clipboard.WriteAll(combinedContent); err != nil {
		return fmt.Errorf("failed to copy combined content to clipboard: %v", err)
	}

	fmt.Println("Copied multiple files' contents to clipboard.")
	return nil
}

// GrabMultipleFolders accepts multiple folder paths, gathers code files from each, and writes the combined content to the clipboard.
func GrabMultipleFolders(folders []string) error {
	var allContents []string

	for _, folder := range folders {
		// Check if the folder exists and is a directory
		info, err := os.Stat(folder)
		if err != nil || !info.IsDir() {
			fmt.Printf("Skipping %s: not a valid folder\n", folder)
			continue
		}

		// Skip if it's a protected workspace
		if isProtectedWorkspace(folder) {
			fmt.Printf("Skipping %s: protected workspace\n", folder)
			continue
		}

		// Count files and confirm if too many files are present
		fileCount := countFiles(folder)
		if fileCount == -1 {
			fmt.Printf("Skipping %s: unable to count files\n", folder)
			continue
		}
		if fileCount > maxFilesLimit {
			fmt.Printf("⚠️ Warning: The directory '%s' contains %d files. Proceed? (y/N): ", folder, fileCount)
			if !confirmAction() {
				fmt.Printf("Skipping %s: too many files to grab\n", folder)
				continue
			}
		}

		// Get the code files' content from the folder
		folderContent, err := getCodesProjectContent(folder)
		if err != nil {
			fmt.Printf("Error grabbing folder %s: %v\n", folder, err)
			continue
		}
		allContents = append(allContents, folderContent)
	}

	if len(allContents) == 0 {
		return fmt.Errorf("no folder content was grabbed")
	}

	// Combine all folder contents with a clear separator
	combinedContent := strings.Join(allContents, "\n===\n")
	if err := clipboard.WriteAll(combinedContent); err != nil {
		return fmt.Errorf("failed to copy combined content to clipboard: %v", err)
	}
	fmt.Println("Copied content of multiple folders to clipboard.")
	return nil
}
