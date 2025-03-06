package implement

import (
	"agent/gorani/internal/tree"
	"fmt"
	"os"
	"os/exec"
)

// CreateGitBranch creates a new Git branch and switches to it.
func CreateGitBranch(branchName string) error {
	// Check if git is installed.
	if _, err := exec.LookPath("git"); err != nil {
		return fmt.Errorf("git is not installed or not found in PATH")
	}

	// Run the git switch -c command to create and switch to the new branch.
	cmd := exec.Command("git", "switch", "-c", branchName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create branch '%s': %v\n%s", branchName, err, string(output))
	}

	fmt.Printf("Successfully created and switched to branch '%s'.\n", branchName)
	return nil
}

// MergeBranch switches to the main branch after committing current changes,
// merges the given branch into main, and then deletes the branch.
func MergeBranch(branchName string) error {
	// Stage all changes.
	cmdAdd := exec.Command("git", "add", ".")
	outputAdd, err := cmdAdd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to add changes: %v\n%s", err, string(outputAdd))
	}
	fmt.Println("Staged changes with 'git add .'")

	// Commit the changes with branchName as the commit message.
	cmdCommit := exec.Command("git", "commit", "-m", branchName)
	outputCommit, err := cmdCommit.CombinedOutput()
	if err != nil {
		// If there is nothing to commit, git returns an error.
		// Log the output and continue.
		fmt.Printf("Warning: commit may have failed (possibly nothing to commit): %v\n%s", err, string(outputCommit))
	} else {
		fmt.Printf("Committed changes with message '%s'.\n", branchName)
	}

	// Switch to main branch.
	cmdSwitch := exec.Command("git", "switch", "main")
	outputSwitch, err := cmdSwitch.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to switch to main: %v\n%s", err, string(outputSwitch))
	}
	fmt.Println("Switched to main branch.")

	// Merge the branch.
	cmdMerge := exec.Command("git", "merge", branchName)
	outputMerge, err := cmdMerge.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to merge branch '%s': %v\n%s", branchName, err, string(outputMerge))
	}
	fmt.Printf("Merged branch '%s' successfully.\n", branchName)

	// Delete the branch after successful merge.
	cmdDelete := exec.Command("git", "branch", "-d", branchName)
	outputDelete, err := cmdDelete.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to delete branch '%s': %v\n%s", branchName, err, string(outputDelete))
	}
	fmt.Printf("Deleted branch '%s' successfully.\n", branchName)
	return nil
}

// PrepareImplementPrompt grabs the tree with functions output and writes a prompt to input.md.
func PrepareImplementPrompt() error {
	// Generate the plain-text tree including function details.
	treeOutput, err := tree.GenerateTreeWithFunctionsString(".", "")
	if err != nil {
		return fmt.Errorf("failed to generate tree with functions: %v", err)
	}

	// Prepare the prompt text.
	promptText := "The following code structure with functions is provided:\n\n" +
		treeOutput +
		"\n\nPlease implement any missing functions or suggest improvements as needed."

	// Write the prompt text to input.md.
	if err := os.WriteFile("input.md", []byte(promptText), 0644); err != nil {
		return fmt.Errorf("failed to write prompt to input.md: %v", err)
	}

	fmt.Println("Prompt prepared in input.md. You can now review/edit it or use your prompt command to send it to OpenAI.")
	return nil
}

// ImplementationManager defines an interface for branch management.
type ImplementationManager interface {
	// CreateBranch creates a new branch.
	CreateBranch(branchName string) error
	// Implement performs the implementation tasks.
	Implement() error
	// MergeBranch merges the given branch into main.
	MergeBranch(branchName string) error
}
