package grab

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"

	"github.com/atotto/clipboard"
)

// exprToString converts an AST expression into its string representation.
func exprToString(expr ast.Expr) string {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, token.NewFileSet(), expr); err != nil {
		return ""
	}
	return buf.String()
}

// buildSummary walks through Go files under the provided root directory,
// and returns a summary of functions, structs, and interfaces along with their package info.
func buildSummary(root string) (string, error) {
	fset := token.NewFileSet()
	var summaryBuffer bytes.Buffer

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err // Propagate errors.
		}
		// Only process .go files.
		if info.IsDir() || filepath.Ext(path) != ".go" {
			return nil
		}

		// Parse the file.
		file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			fmt.Fprintf(&summaryBuffer, "Error parsing file %s: %v\n", path, err)
			return nil // Skip files with errors.
		}

		// Get package name.
		packageName := file.Name.Name
		summaryBuffer.WriteString(fmt.Sprintf("\nFile: %s (package %s)\n", path, packageName))

		// Walk the AST.
		ast.Inspect(file, func(n ast.Node) bool {
			switch node := n.(type) {
			case *ast.FuncDecl:
				recv := ""
				if node.Recv != nil && len(node.Recv.List) > 0 {
					recv = fmt.Sprintf("(%s) ", exprToString(node.Recv.List[0].Type))
				}
				summaryBuffer.WriteString(fmt.Sprintf("  [Package: %s] Function: %s%s\n", packageName, recv, node.Name.Name))
			case *ast.GenDecl:
				if node.Tok == token.TYPE {
					for _, spec := range node.Specs {
						typeSpec, ok := spec.(*ast.TypeSpec)
						if !ok {
							continue
						}
						switch typeSpec.Type.(type) {
						case *ast.StructType:
							summaryBuffer.WriteString(fmt.Sprintf("  [Package: %s] Struct: %s\n", packageName, typeSpec.Name.Name))
						case *ast.InterfaceType:
							summaryBuffer.WriteString(fmt.Sprintf("  [Package: %s] Interface: %s\n", packageName, typeSpec.Name.Name))
						}
					}
				}
			}
			return true
		})

		return nil
	})
	if err != nil {
		return "", fmt.Errorf("error walking the path %s: %v", root, err)
	}

	summary := summaryBuffer.String()
	if summary == "" {
		summary = "No Go symbols found."
	}
	return summary, nil
}

// GrabSummary generates a summary of Go symbols from the provided root,
// copies the summary to the clipboard, and prints a confirmation message.
func GrabSummary(root string) error {
	summary, err := buildSummary(root)
	if err != nil {
		return err
	}

	// Copy the summary to the clipboard.
	if err := clipboard.WriteAll(summary); err != nil {
		return fmt.Errorf("failed to copy summary to clipboard: %v", err)
	}
	fmt.Println("Summary copied to clipboard.")
	return nil
}
