package grab

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/atotto/clipboard"
)

// GrabPublicFuncs walks through Go files under the provided root directory,
// collects detailed information about public functions (exported functions),
// and copies the accumulated output to the clipboard.
func GrabPublicFuncs(root string) error {
	var outputBuffer bytes.Buffer

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Only process Go files (skip directories and non-Go files)
		if info.IsDir() || filepath.Ext(path) != ".go" {
			return nil
		}

		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return fmt.Errorf("error parsing file %s: %v", path, err)
		}

		// Buffer to collect public functions from this file.
		var fileOutput bytes.Buffer

		// Iterate over declarations to find function declarations.
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			// Check if the function is public (exported) using ast.IsExported.
			if !ast.IsExported(fn.Name.Name) {
				continue
			}

			// Build parameter list.
			var params []string
			if fn.Type.Params != nil {
				for _, field := range fn.Type.Params.List {
					typeStr := exprToString(field.Type)
					// If no parameter name is provided, only the type is recorded.
					if len(field.Names) == 0 {
						params = append(params, typeStr)
					} else {
						for _, name := range field.Names {
							params = append(params, fmt.Sprintf("%s %s", name.Name, typeStr))
						}
					}
				}
			}
			paramStr := strings.Join(params, ", ")

			// Build result list.
			var results []string
			if fn.Type.Results != nil {
				for _, field := range fn.Type.Results.List {
					results = append(results, exprToString(field.Type))
				}
			}
			resultStr := ""
			if len(results) > 0 {
				resultStr = fmt.Sprintf(" -> %s", strings.Join(results, ", "))
			}

			fileOutput.WriteString(fmt.Sprintf("  Public Function: %s(%s)%s\n", fn.Name.Name, paramStr, resultStr))
		}

		// Only add output for the file if at least one public function was found.
		if fileOutput.Len() > 0 {
			outputBuffer.WriteString(fmt.Sprintf("File: %s\n", path))
			outputBuffer.Write(fileOutput.Bytes())
			outputBuffer.WriteString("\n")
		}
		return nil
	})
	if err != nil {
		return err
	}

	// Copy the collected output to the clipboard or inform the user if none was found.
	if outputBuffer.Len() == 0 {
		fmt.Println("No public functions found.")
	} else {
		if err := clipboard.WriteAll(outputBuffer.String()); err != nil {
			return fmt.Errorf("failed to copy to clipboard: %v", err)
		}
		fmt.Println("Copied public functions to clipboard.")
	}
	return nil
}
