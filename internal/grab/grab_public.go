package grab

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

// funcRegex matches Go function definitions.
var funcRegex = regexp.MustCompile(`(?m)^\s*func\s+(\(\w+\s+\*?\w+\)\s+)?([A-Z]\w*)\s*\(([^)]*)\)\s*([^{]*)`)

// GrabPublicFuncsWithDescriptions extracts public function names and their descriptions.
func GrabPublicFuncsWithDescriptions(root string) (map[string]string, error) {
	functions := make(map[string]string)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			fileFuncs, err := extractPublicFuncsWithDescriptions(path)
			if err != nil {
				return err
			}
			for name, desc := range fileFuncs {
				functions[name] = desc
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return functions, nil
}

// extractPublicFuncsWithDescriptions reads a Go file and extracts public function names with their descriptions.
func extractPublicFuncsWithDescriptions(filePath string) (map[string]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	functions := make(map[string]string)
	scanner := bufio.NewScanner(file)
	var lastComment string

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		// Capture comment lines
		if strings.HasPrefix(trimmed, "//") {
			lastComment = strings.TrimPrefix(trimmed, "// ")
			continue
		}

		// Match public function definition
		matches := funcRegex.FindStringSubmatch(line)
		if len(matches) > 2 {
			funcName := matches[2]

			// Ensure it's a public function (starts with an uppercase letter)
			if unicode.IsUpper(rune(funcName[0])) {
				functions[funcName] = lastComment
			}

			// Reset lastComment after function is found
			lastComment = ""
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return functions, nil
}

// PrintPublicFunctions prints the extracted public functions and their descriptions.
func PrintPublicFunctions(root string) error {
	functions, err := GrabPublicFuncsWithDescriptions(root)
	if err != nil {
		return err
	}

	fmt.Println("Public Functions and Descriptions:")
	for name, desc := range functions {
		fmt.Printf("- %s: %s\n", name, desc)
	}

	return nil
}
