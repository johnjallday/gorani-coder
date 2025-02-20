package tree

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/fatih/color"
)

// Function regex to match Go function definitions and extract parameters
var funcRegex = regexp.MustCompile(`\bfunc\s+(\(\w+\s+\*?\w+\)\s+)?(\w+)\s*\(([^)]*)\)\s*([^{]*)`)

// Colored formatters
var (
	dirColor        = color.New(color.FgCyan, color.Bold)    // Cyan for directories
	fileColor       = color.New(color.FgYellow)              // Yellow for files
	funcColor       = color.New(color.FgGreen, color.Bold)   // Green for function names
	publicFuncColor = color.New(color.FgHiGreen, color.Bold) // Light green for public function names
	paramColor      = color.New(color.FgMagenta)             // Magenta for parameters
	returnColor     = color.New(color.FgBlue)                // Blue for return types
	connectorColor  = color.New(color.FgWhite)               // White for tree connectors
)

// shouldIgnoreFile checks if the file should be ignored (e.g., .DS_Store, hidden files)
func shouldIgnoreFile(name string) bool {
	return strings.HasPrefix(name, ".") // Ignore hidden files (e.g., .DS_Store, .git, etc.)
}

// PrintTree recursively prints the directory tree structure with colors (ignoring .DS_Store)
func PrintTree(root string, indent string) error {
	entries, err := os.ReadDir(root)
	if err != nil {
		return err
	}

	for i, entry := range entries {
		if shouldIgnoreFile(entry.Name()) {
			continue // Skip ignored files
		}

		connector := "├──"
		if i == len(entries)-1 {
			connector = "└──"
		}

		connectorColor.Print(indent + connector + " ")

		// Print directories in cyan, files in yellow
		if entry.IsDir() {
			dirColor.Println(entry.Name()) // Print directory
			subIndent := indent + "│   "
			if i == len(entries)-1 {
				subIndent = indent + "    "
			}
			PrintTree(filepath.Join(root, entry.Name()), subIndent)
		} else {
			fileColor.Println(entry.Name()) // Print file
		}
	}

	return nil
}

// PrintTreeWithFunctions prints the directory tree structure and functions inside Go files with parameters (ignoring .DS_Store)
func PrintTreeWithFunctions(root string, indent string) error {
	entries, err := os.ReadDir(root)
	if err != nil {
		return err
	}

	for i, entry := range entries {
		if shouldIgnoreFile(entry.Name()) {
			continue // Skip ignored files
		}

		connector := "├──"
		if i == len(entries)-1 {
			connector = "└──"
		}

		connectorColor.Print(indent + connector + " ")
		if entry.IsDir() {
			dirColor.Println(entry.Name()) // Print directory in cyan
			subIndent := indent + "│   "
			if i == len(entries)-1 {
				subIndent = indent + "    "
			}
			PrintTreeWithFunctions(filepath.Join(root, entry.Name()), subIndent)
		} else if strings.HasSuffix(entry.Name(), ".go") {
			fileColor.Println(entry.Name()) // Print file in yellow

			// If it's a Go file, extract and print function names with parameters
			funcs, err := extractFunctions(filepath.Join(root, entry.Name()))
			if err == nil && len(funcs) > 0 {
				for _, f := range funcs {
					connectorColor.Print(indent + "    ├── ")
					fmt.Println(f) // Already formatted inside extractFunctions()
				}
			}
		} else {
			fileColor.Println(entry.Name()) // Print other files in yellow
		}
	}

	return nil
}

// extractFunctions extracts function names, parameters, and return types from a Go source file
func extractFunctions(filePath string) ([]string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	matches := funcRegex.FindAllStringSubmatch(string(content), -1)
	var functions []string
	for _, match := range matches {
		funcName := match[2]                   // Extract function name
		params := match[3]                     // Extract function parameters
		returns := strings.TrimSpace(match[4]) // Extract return types

		// Determine the color based on whether the function is public or private
		funcColorToUse := funcColor
		if unicode.IsUpper(rune(funcName[0])) {
			funcColorToUse = publicFuncColor
		}

		// Format parameters with color
		paramList := []string{}
		if params != "" {
			for _, param := range strings.Split(params, ",") {
				paramList = append(paramList, paramColor.Sprint(strings.TrimSpace(param)))
			}
		}

		// Format return types with color
		returnStr := ""
		if returns != "" {
			returnStr = returnColor.Sprintf(" -> %s", strings.TrimSpace(returns))
		}

		// Store formatted function output
		functions = append(functions, fmt.Sprintf("%s(%s)%s", funcColorToUse.Sprint(funcName), strings.Join(paramList, ", "), returnStr))
	}

	return functions, nil
}
