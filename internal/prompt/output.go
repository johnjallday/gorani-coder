package prompt

import (
	"encoding/json"
	"fmt"
	"os"
)

func ProcessScriptsFromOutputFile() error {
	filePath := "output.md"

	// 1. Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read %s: %v", filePath, err)
	}

	// 2. Unmarshal into CodeResponse
	var codeResp CodeResponse
	if err := json.Unmarshal(data, &codeResp); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	// 3. Print the filename
	fmt.Println("Filename from JSON:", codeResp.Filename)

	// 4. Overwrite (or create) a new file with the content of the *first* script
	if len(codeResp.Scripts) > 0 {
		firstScript := codeResp.Scripts[0]

		// Overwrite a file named by "Filename" with the script’s content
		if err := os.WriteFile(codeResp.Filename, []byte(firstScript), 0644); err != nil {
			return fmt.Errorf("failed to write file %q: %v", codeResp.Filename, err)
		}
		fmt.Printf("✅ Successfully wrote the first script to %q\n", codeResp.Filename)
	} else {
		fmt.Println("No scripts found in the JSON response.")
	}

	return nil
}
