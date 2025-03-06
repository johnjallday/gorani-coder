package version

import (
	"fmt"
	"os"
)

// WriteReadme generates a README.md file with basic information.
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
