package grab

import (
	"agent/gorani/internal/prompt"
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// SmartGrab generates a summary of Go symbols from the provided root,
// determines the current git branch (if it's not "main" or "origin"),
// builds a feature note, saves the combined prompt into input.md, and prints a confirmation.

func SmartGrab(root string) error {
	featureBranch, err := getFeatureBranch()
	if err != nil {
		return err
	}
	if featureBranch == "" {
		// No valid branch was found (active branch is "main" or "origin", or none exists).
		fmt.Println("No valid feature branch found. Aborting SmartGrab.")
		return nil
	}

	grabPrompt, err := buildPrompt(featureBranch, root)
	if err != nil {
		return err
	}

	// Save the prompt to input.md.
	if err := os.WriteFile("input.md", []byte(grabPrompt), 0644); err != nil {
		return fmt.Errorf("failed to save prompt to input.md: %v", err)
	}
	fmt.Println("Prompt saved to input.md.")

	input, err := os.ReadFile("input.md")
	if err != nil {
		return fmt.Errorf("error reading input.md file: %w", err)
	}

	prompt.PromptOpenaiFiles(string(input)) // no return value
	fmt.Println("Prompt sent to OpenAI.")

	// handle output
	// read output.md
	output, err := os.ReadFile("output.md")
	if err != nil {
		return fmt.Errorf("error reading output.md file: %w", err)
	}

	// Define a struct that matches the output.md JSON structure
	type Output struct {
		Files []string `json:"files"`
	}

	var parsedOutput Output
	if err := json.Unmarshal(output, &parsedOutput); err != nil {
		return fmt.Errorf("error parsing output.md JSON: %w", err)
	}

	fmt.Println("Files extracted from output.md:", parsedOutput.Files)

	// Pass the list of files to GrabFiles
	if err := GrabFiles(parsedOutput.Files); err != nil {
		return fmt.Errorf("error grabbing files: %w", err)
	}

	return nil
}

// getFeatureBranch retrieves the active git branch and returns it if it is not "main" or "origin".
// If no active branch is found or if the active branch is "main" or "origin", it returns an empty string.
func getFeatureBranch() (string, error) {
	// Retrieve the list of git branches.
	cmd := exec.Command("git", "branch")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to execute 'git branch': %v", err)
	}

	// Find the active branch (the line starting with '*').
	var currentBranch string
	branches := strings.Split(string(output), "\n")
	for _, line := range branches {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "*") {
			// Remove the '*' and trim spaces.
			currentBranch = strings.TrimSpace(strings.TrimPrefix(line, "*"))
			break
		}
	}

	// If no active branch is found, abort.
	if currentBranch == "" {
		fmt.Println("No active branch found. Aborting SmartGrab.")
		return "", nil
	}

	// Abort if the active branch is "main" or "origin".
	if currentBranch == "main" || currentBranch == "origin" {
		fmt.Printf("Active branch is '%s'. Aborting SmartGrab.\n", currentBranch)
		return "", nil
	}

	return currentBranch, nil
}

// buildPrompt builds a prompt string using the feature branch name, a user-provided feature description,
// and the code summary from the given root. It returns the combined prompt string or an error if the summary cannot be generated.
func buildPrompt(featureBranch, root string) (string, error) {
	// Print the feature branch.
	fmt.Printf("Feature Branch: %s\n", featureBranch)

	// Ask the user for a detailed feature description.
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Please enter a detailed feature description: ")
	description, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read feature description: %v", err)
	}
	description = strings.TrimSpace(description)

	// Generate the code summary.
	summary, err := buildSummary(root)
	if err != nil {
		return "", err
	}

	// Build the combined prompt.
	promptText := fmt.Sprintf("I want to build a feature called %s.\n", featureBranch)
	promptText += "Feature Description:\n" + description + "\n\n"
	promptText += "Here is the summary of the code:\n" + summary + "\n\n"
	promptText += "Give me a list of files needed in order to build this feature."
	return promptText, nil
}
