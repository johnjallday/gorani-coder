package version

import (
	"fmt"
	"os"
	"os/exec"
)

// WriteReadme generates a README.md file with basic information
func WriteReadme() error {
	readmeContent := `# Project Title

## Description
A brief description of your project.

## Installation
Instructions for installation.

## Usage
How to use this project.

## License
Specify the license details here.
`

	file, err := os.Create("README.md")
	if err != nil {
		return fmt.Errorf("failed to create README.md: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(readmeContent)
	if err != nil {
		return fmt.Errorf("failed to write to README.md: %v", err)
	}

	fmt.Println("README.md created successfully.")
	return nil
}

// CreateGitBranch creates a new Git branch and switches to it.
func CreateGitBranch(branchName string) error {
	// Check if git is installed
	if _, err := exec.LookPath("git"); err != nil {
		return fmt.Errorf("git is not installed or not found in PATH")
	}

	// Run the git checkout -b command to create and switch to the new branch
	cmd := exec.Command("git", "checkout", "-b", branchName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create branch '%s': %v\n%s", branchName, err, string(output))
	}

	fmt.Printf("Successfully created and switched to branch '%s'.\n", branchName)
	return nil
}
