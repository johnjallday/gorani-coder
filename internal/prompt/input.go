package prompt

import (
	"fmt"
	"os"
	"os/exec"
)

// OpenInputInNeovim opens "input.md" in Neovim and returns the edited content.
func OpenInputInNeovim() (string, error) {
	filePath := "input.md"

	// Open input.md in Neovim
	cmd := exec.Command("nvim", filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to open input.md in Neovim: %v", err)
	}

	// Read the content after editing
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read input.md: %v", err)
	}

	return string(content), nil
}
