package commandbuilder

import (
	"agent/gorani/internal/grab"
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// RegisterActions scans for public functions, asks the user which ones to register,
// and generates a Go file that creates a Cobra command for each selected function.
func RegisterActions() {
	rootDir := "./internal"
	funcs, err := grab.GrabPublicFuncsWithDescriptions(rootDir)
	if err != nil {
		fmt.Printf("Error grabbing public functions: %v\n", err)
		return
	}

	if len(funcs) == 0 {
		fmt.Println("No public functions found.")
		return
	}

	// Create a sorted list of function names for predictable ordering.
	var keys []string
	for k := range funcs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Display the list of available functions.
	fmt.Println("Available public functions:")
	for i, fname := range keys {
		desc := funcs[fname]
		fmt.Printf("%d. %s: %s\n", i+1, fname, desc)
	}

	// Prompt the user for a selection.
	fmt.Println("Enter the numbers of the functions you want to register, separated by commas:")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading input: %v\n", err)
		return
	}
	input = strings.TrimSpace(input)
	if input == "" {
		fmt.Println("No functions selected. Aborting.")
		return
	}

	selections := strings.Split(input, ",")
	var selectedFunctions []string
	for _, s := range selections {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		idx, err := strconv.Atoi(s)
		if err != nil {
			fmt.Printf("Invalid number: %s\n", s)
			continue
		}
		if idx < 1 || idx > len(keys) {
			fmt.Printf("Selection %d out of range.\n", idx)
			continue
		}
		selectedFunctions = append(selectedFunctions, keys[idx-1])
	}

	if len(selectedFunctions) == 0 {
		fmt.Println("No valid functions selected. Aborting.")
		return
	}

	// Determine the output file. Here we create a file under the cmd directory.
	outputDir := "./cmd"
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			fmt.Printf("Failed to create directory %s: %v\n", outputDir, err)
			return
		}
	}
	outputFile := filepath.Join(outputDir, "registered_actions.go")

	// Generate the file content.
	// It registers a new Cobra command for each selected function.
	code := "package cmd\n\n"
	code += "import (\n\t\"fmt\"\n\t\"github.com/spf13/cobra\"\n)\n\n"
	code += "func init() {\n"
	for _, fn := range selectedFunctions {
		// Use the function name as the command name and use its description.
		code += fmt.Sprintf("\trootCmd.AddCommand(&cobra.Command{\n")
		code += fmt.Sprintf("\t\tUse: \"%s\",\n", fn)
		code += fmt.Sprintf("\t\tShort: \"%s\",\n", escapeString(funcs[fn]))
		code += "\t\tRun: func(cmd *cobra.Command, args []string) {\n"
		code += fmt.Sprintf("\t\t\tfmt.Println(\"Action for %s is executed\")\n", fn)
		code += "\t\t},\n"
		code += "\t})\n\n"
	}
	code += "}\n"

	// Write the generated code to the file.
	if err := os.WriteFile(outputFile, []byte(code), 0644); err != nil {
		fmt.Printf("Error writing file %s: %v\n", outputFile, err)
		return
	}

	fmt.Printf("Registered actions file created: %s\n", outputFile)
}

// escapeString escapes double quotes for safe inclusion in a Go string literal.
func escapeString(s string) string {
	return strings.ReplaceAll(s, "\"", "\\\"")
}
