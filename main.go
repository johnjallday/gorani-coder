package main

import (
	"agent/gorani/internal/command"
	"os"
)

func main() {
	// Pass the command-line arguments (excluding the program name) to the commands executor.
	command.Execute(os.Args[1:])
}
