package main

import (
	"agent/coder/internal/docbuilder"
	"agent/coder/internal/grab"
	"agent/coder/internal/prompt"
	"agent/coder/internal/tree"
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

type Command struct {
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Handler     func(args []string) `json:"-"`
}

var commands map[string]Command

func init() {
	commands = map[string]Command{
		"tree": {
			Name:        "tree",
			Description: "Prints the directory tree structure",
			Handler: func(args []string) {
				fmt.Println("Printing Directory Tree:")
				path := "."
				if len(args) > 0 {
					path = args[0]
				}
				if err := tree.PrintTree(path, ""); err != nil {
					fmt.Println("Error printing tree:", err)
				}
			},
		},
		"tree-func": {
			Name:        "tree-func",
			Description: "Prints the directory tree structure with functions",
			Handler: func(args []string) {
				fmt.Println("Printing Directory Tree with Functions:")
				path := "."
				if len(args) > 0 {
					path = args[0]
				}
				if err := tree.PrintTreeWithFunctions(path, ""); err != nil {
					fmt.Println("Error printing tree with functions:", err)
				}
			},
		},
		"prompt": {
			Name:        "prompt",
			Description: "Prompts OpenAI with user input",
			Handler: func(args []string) {
				fmt.Println("Enter your prompt:")
				prompt.PromptFromNeovim()
			},
		},
		"grab": {
			Name:        "grab",
			Description: "Grabs code files (file or folder auto-detected)",
			Handler: func(args []string) {
				if len(args) < 1 {
					grab.Grab("./")
					return
				}
				if err := grab.Grab(args[0]); err != nil {
					fmt.Println("Error:", err)
				}
			},
		},
		"grab-summary": {
			Name:        "grab-summary",
			Description: "Grabs and prints detailed symbol info (package, functions, structs, interfaces) from Go files",
			Handler: func(args []string) {
				path := "."
				if len(args) > 0 {
					path = args[0]
				}
				if err := grab.GrabSummary(path); err != nil {
					fmt.Println("Error grabbing summary:", err)
				}
			},
		},
		"docbuilder": {
			Name:        "docbuilder",
			Description: "Generates README.md using documentation builder",
			Handler: func(args []string) {
				docbuilder.BuildReadme()
			},
		},
		"commands": {
			Name:        "commands",
			Description: "Displays available commands",
			Handler: func(args []string) {
				printCommands()
			},
		},
	}
}

func printCommands() {
	fmt.Println("Usage: go run main.go <command> <arguments>")
	fmt.Println("Available commands:")

	keys := make([]string, 0, len(commands))
	for k := range commands {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Printf("  %-12s - %s\n", k, commands[k].Description)
	}
}

func main() {
	// Special flag: if "--list-commands" is provided, output the commands as JSON.
	if len(os.Args) > 1 && os.Args[1] == "--list-commands" {
		var cmds []Command
		for _, cmd := range commands {
			cmds = append(cmds, cmd)
		}
		sort.Slice(cmds, func(i, j int) bool {
			return cmds[i].Name < cmds[j].Name
		})
		out, err := json.Marshal(cmds)
		if err != nil {
			fmt.Println("Error marshaling commands:", err)
			os.Exit(1)
		}
		fmt.Println(string(out))
		return
	}

	if len(os.Args) < 2 {
		printCommands()
		os.Exit(1)
	}

	cmd, exists := commands[os.Args[1]]
	if !exists {
		fmt.Println("Unknown command. Use 'commands' for available commands.")
		printCommands()
		return
	}

	// Pass remaining args to the command handler.
	cmd.Handler(os.Args[2:])
}
