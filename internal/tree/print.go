package tree

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/atotto/clipboard"
	"github.com/fatih/color"
)

// funcRegex matches Go function definitions.
var funcRegex = regexp.MustCompile(`\bfunc\s+(\(\w+\s+\*?\w+\)\s+)?(\w+)\s*\(([^)]*)\)\s*([^{]*)`)

// Colored formatters for terminal output.
var (
	dirColor        = color.New(color.FgCyan, color.Bold)    // Directories in cyan.
	fileColor       = color.New(color.FgYellow)              // Files in yellow.
	funcColor       = color.New(color.FgGreen, color.Bold)   // Functions in green.
	publicFuncColor = color.New(color.FgHiGreen, color.Bold) // Public functions in light green.
	paramColor      = color.New(color.FgMagenta)             // Parameters in magenta.
	returnColor     = color.New(color.FgBlue)                // Return types in blue.
	connectorColor  = color.New(color.FgWhite)               // Tree connectors in white.
)

// shouldIgnoreFile skips hidden files like .DS_Store or .git.
func shouldIgnoreFile(name string) bool {
	return strings.HasPrefix(name, ".")
}

// ---------------------
// Colorful Tree Printers
// ---------------------

// PrintTree prints a colored directory tree.
func PrintTree(root, indent string) error {
	entries, err := os.ReadDir(root)
	if err != nil {
		return err
	}

	// Filter out ignored entries.
	var validEntries []os.DirEntry
	for _, entry := range entries {
		if !shouldIgnoreFile(entry.Name()) {
			validEntries = append(validEntries, entry)
		}
	}

	for i, entry := range validEntries {
		connector := "├──"
		if i == len(validEntries)-1 {
			connector = "└──"
		}

		connectorColor.Print(indent + connector + " ")
		if entry.IsDir() {
			dirColor.Println(entry.Name())
			subIndent := indent + "│   "
			if i == len(validEntries)-1 {
				subIndent = indent + "    "
			}
			if err := PrintTree(filepath.Join(root, entry.Name()), subIndent); err != nil {
				return err
			}
		} else {
			fileColor.Println(entry.Name())
		}
	}

	return nil
}

// PrintTreeWithFunctions prints a colored tree and, for Go files, extracts and prints functions.
func PrintTreeWithFunctions(root, indent string) error {
	entries, err := os.ReadDir(root)
	if err != nil {
		return err
	}

	var validEntries []os.DirEntry
	for _, entry := range entries {
		if !shouldIgnoreFile(entry.Name()) {
			validEntries = append(validEntries, entry)
		}
	}

	for i, entry := range validEntries {
		connector := "├──"
		if i == len(validEntries)-1 {
			connector = "└──"
		}

		connectorColor.Print(indent + connector + " ")
		if entry.IsDir() {
			dirColor.Println(entry.Name())
			subIndent := indent + "│   "
			if i == len(validEntries)-1 {
				subIndent = indent + "    "
			}
			if err := PrintTreeWithFunctions(filepath.Join(root, entry.Name()), subIndent); err != nil {
				return err
			}
		} else if strings.HasSuffix(entry.Name(), ".go") {
			fileColor.Println(entry.Name())
			funcs, err := extractFunctions(filepath.Join(root, entry.Name()))
			if err == nil && len(funcs) > 0 {
				for _, f := range funcs {
					connectorColor.Print(indent + "    ├── ")
					// f already includes ANSI color codes.
					fmt.Println(f)
				}
			}
		} else {
			fileColor.Println(entry.Name())
		}
	}

	return nil
}

// extractFunctions extracts function details from a Go source file.
func extractFunctions(filePath string) ([]string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	matches := funcRegex.FindAllStringSubmatch(string(content), -1)
	var functions []string
	for _, match := range matches {
		funcName := match[2]
		params := match[3]
		returns := strings.TrimSpace(match[4])

		// Use a different color for public functions.
		colorToUse := funcColor
		if unicode.IsUpper(rune(funcName[0])) {
			colorToUse = publicFuncColor
		}

		// Colorize parameters.
		var paramList []string
		if params != "" {
			for _, param := range strings.Split(params, ",") {
				paramList = append(paramList, paramColor.Sprint(strings.TrimSpace(param)))
			}
		}

		returnStr := ""
		if returns != "" {
			returnStr = returnColor.Sprintf(" -> %s", strings.TrimSpace(returns))
		}

		functions = append(functions, fmt.Sprintf("%s(%s)%s", colorToUse.Sprint(funcName), strings.Join(paramList, ", "), returnStr))
	}
	return functions, nil
}

// ---------------------
// Plain-Text Generators & Clipboard Functions
// ---------------------

// GenerateTreeString builds a plain-text tree representation (without ANSI colors).
func GenerateTreeString(root, indent string) (string, error) {
	var sb strings.Builder
	entries, err := os.ReadDir(root)
	if err != nil {
		return "", err
	}

	var validEntries []os.DirEntry
	for _, entry := range entries {
		if !shouldIgnoreFile(entry.Name()) {
			validEntries = append(validEntries, entry)
		}
	}

	for i, entry := range validEntries {
		connector := "├──"
		if i == len(validEntries)-1 {
			connector = "└──"
		}
		sb.WriteString(indent + connector + " " + entry.Name() + "\n")

		if entry.IsDir() {
			subIndent := indent + "│   "
			if i == len(validEntries)-1 {
				subIndent = indent + "    "
			}
			subTree, err := GenerateTreeString(filepath.Join(root, entry.Name()), subIndent)
			if err != nil {
				return "", err
			}
			sb.WriteString(subTree)
		}
	}

	return sb.String(), nil
}

// GenerateTreeWithFunctionsString builds a plain-text tree including Go function details.
func GenerateTreeWithFunctionsString(root, indent string) (string, error) {
	var sb strings.Builder
	entries, err := os.ReadDir(root)
	if err != nil {
		return "", err
	}

	var validEntries []os.DirEntry
	for _, entry := range entries {
		if !shouldIgnoreFile(entry.Name()) {
			validEntries = append(validEntries, entry)
		}
	}

	for i, entry := range validEntries {
		connector := "├──"
		if i == len(validEntries)-1 {
			connector = "└──"
		}
		sb.WriteString(indent + connector + " " + entry.Name() + "\n")

		if entry.IsDir() {
			subIndent := indent + "│   "
			if i == len(validEntries)-1 {
				subIndent = indent + "    "
			}
			subTree, err := GenerateTreeWithFunctionsString(filepath.Join(root, entry.Name()), subIndent)
			if err != nil {
				return "", err
			}
			sb.WriteString(subTree)
		} else if strings.HasSuffix(entry.Name(), ".go") {
			funcs, err := extractFunctions(filepath.Join(root, entry.Name()))
			if err == nil && len(funcs) > 0 {
				for _, f := range funcs {
					sb.WriteString(indent + "    ├── " + f + "\n")
				}
			}
		}
	}

	return sb.String(), nil
}

// CopyTreeToClipboard generates the plain-text tree, prints it, and copies it to the clipboard.
func CopyTreeToClipboard(root string) error {
	treeStr, err := GenerateTreeString(root, "")
	if err != nil {
		return err
	}

	// Print the plain-text tree.
	fmt.Println(treeStr)

	// Copy to clipboard.
	return clipboard.WriteAll(treeStr)
}

// CopyTreeWithFunctionsToClipboard generates the plain-text tree with functions, prints it, and copies it.
func CopyTreeWithFunctionsToClipboard(root string) error {
	treeStr, err := GenerateTreeWithFunctionsString(root, "")
	if err != nil {
		return err
	}

	// Print the plain-text tree.
	fmt.Println(treeStr)

	// Copy to clipboard.
	return clipboard.WriteAll(treeStr)
}

// ---------------------
// Main: Demonstration
// ---------------------

func main() {
	rootDir := "./" // Change this to your target directory

	// Print the colored tree.
	fmt.Println("Colored Tree:")
	if err := PrintTree(rootDir, ""); err != nil {
		fmt.Println("Error printing colored tree:", err)
	}

	// Print the colored tree with function details.
	fmt.Println("\nColored Tree with Functions:")
	if err := PrintTreeWithFunctions(rootDir, ""); err != nil {
		fmt.Println("Error printing colored tree with functions:", err)
	}

	// Print and copy the plain-text tree.
	fmt.Println("\nPlain-Text Tree (copied to clipboard):")
	if err := CopyTreeToClipboard(rootDir); err != nil {
		fmt.Println("Error copying plain-text tree to clipboard:", err)
	} else {
		fmt.Println("Plain-text tree copied to clipboard!")
	}

	// Print and copy the plain-text tree with functions.
	fmt.Println("\nPlain-Text Tree with Functions (copied to clipboard):")
	if err := CopyTreeWithFunctionsToClipboard(rootDir); err != nil {
		fmt.Println("Error copying plain-text tree with functions to clipboard:", err)
	} else {
		fmt.Println("Plain-text tree with functions copied to clipboard!")
	}
}
