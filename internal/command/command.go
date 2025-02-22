package command

import (
	"agent/gorani/internal/docbuilder"
	"agent/gorani/internal/grab"
	"agent/gorani/internal/prompt"
	"agent/gorani/internal/tree"
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

// Command defines the interface that all commands must implement.
type Command interface {
	Name() string
	Description() string
	Execute(args []string) error
}

// BaseCommand provides a simple implementation of Command using a handler function.
type BaseCommand struct {
	name        string
	description string
	handler     func(args []string) error
}

func (c BaseCommand) Name() string {
	return c.name
}

func (c BaseCommand) Description() string {
	return c.description
}

func (c BaseCommand) Execute(args []string) error {
	return c.handler(args)
}

// NewCommand is a helper to create a new Command.
func NewCommand(name, description string, handler func(args []string) error) Command {
	return BaseCommand{name: name, description: description, handler: handler}
}

// commands is the registry of all available commands.
var commands = map[string]Command{}

func init() {
	// Register the "tree" command.
	commands["tree"] = NewCommand("tree", "Prints the directory tree structure", func(args []string) error {
		fmt.Println("Printing Directory Tree:")
		path := "."
		if len(args) > 0 {
			path = args[0]
		}
		return tree.PrintTree(path, "")
	})

	// Register the "tree-func" command.
	commands["tree-func"] = NewCommand("tree-func", "Prints the directory tree structure with functions", func(args []string) error {
		fmt.Println("Printing Directory Tree with Functions:")
		path := "."
		if len(args) > 0 {
			path = args[0]
		}
		return tree.PrintTreeWithFunctions(path, "")
	})

	// Register the "prompt" command.
	commands["prompt"] = NewCommand("prompt", "Prompts OpenAI with user input", func(args []string) error {
		fmt.Println("Enter your prompt:")
		prompt.PromptFromNeovim()
		return nil
	})

	// Register the "grab" command.
	commands["grab"] = NewCommand("grab", "Grabs code files (file or folder auto-detected)", func(args []string) error {
		if len(args) < 1 {
			return grab.Grab("./")
		}
		return grab.Grab(args[0])
	})

	// Register the "docbuilder" command.
	commands["docbuilder"] = NewCommand("docbuilder", "Generates README.md using documentation builder", func(args []string) error {
		docbuilder.BuildReadme()
		return nil
	})

	// Register the "commands" command.
	commands["commands"] = NewCommand("commands", "Displays available commands", func(args []string) error {
		PrintCommands()
		return nil
	})
}

// PrintCommands prints the list of available commands.
func PrintCommands() {
	fmt.Println("Usage: go run main.go <command> <arguments>")
	fmt.Println("Available commands:")
	keys := make([]string, 0, len(commands))
	for k := range commands {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		cmd := commands[k]
		fmt.Printf("  %-12s - %s\n", cmd.Name(), cmd.Description())
	}
}

// ListCommandsJSON prints the list of commands in JSON format.
func ListCommandsJSON() {
	var cmds []Command
	for _, cmd := range commands {
		cmds = append(cmds, cmd)
	}
	sort.Slice(cmds, func(i, j int) bool {
		return cmds[i].Name() < cmds[j].Name()
	})
	out, err := json.Marshal(cmds)
	if err != nil {
		fmt.Println("Error marshaling commands:", err)
		os.Exit(1)
	}
	fmt.Println(string(out))
}

// Execute parses the command‑line arguments and executes the corresponding command.
func Execute(args []string) {
	// Special flag: if "--list-commands" is provided, output the commands as JSON.
	if len(args) > 0 && args[0] == "--list-commands" {
		ListCommandsJSON()
		return
	}

	if len(args) < 1 {
		PrintCommands()
		os.Exit(1)
	}

	cmdName := args[0]
	cmd, exists := commands[cmdName]
	if !exists {
		fmt.Println("Unknown command. Use 'commands' for available commands.")
		PrintCommands()
		os.Exit(1)
	}

	if err := cmd.Execute(args[1:]); err != nil {
		fmt.Printf("Error executing command %s: %v\n", cmdName, err)
		os.Exit(1)
	}
}
