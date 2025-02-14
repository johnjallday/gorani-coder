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
	"strings"

	"github.com/atotto/clipboard"
)

// exprToString converts an AST expression into its string representation.
func exprToString(expr ast.Expr) string {
	var buf bytes.Buffer
	err := printer.Fprint(&buf, token.NewFileSet(), expr)
	if err != nil {
		return ""
	}
	return buf.String()
}

// GrabSummary walks through Go files under the provided root directory,
// collects detailed symbol info (including file, package, functions, structs, and interfaces),
// and copies the accumulated output to the clipboard.
func GrabSummary(root string) error {
	var outputBuffer bytes.Buffer

	// Walk the directory tree
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Only process .go files (skip directories and non-Go files)
		if info.IsDir() || filepath.Ext(path) != ".go" {
			return nil
		}

		fset := token.NewFileSet()
		// Parse the file
		file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return fmt.Errorf("error parsing file %s: %v", path, err)
		}

		// Append file and package info to the buffer
		outputBuffer.WriteString(fmt.Sprintf("File: %s\n", path))
		outputBuffer.WriteString(fmt.Sprintf("Package: %s\n", file.Name.Name))

		// Iterate over declarations to find symbols
		for _, decl := range file.Decls {
			switch d := decl.(type) {
			case *ast.FuncDecl:
				// Function declaration: include parameters and results if available
				var params []string
				if d.Type.Params != nil {
					for _, field := range d.Type.Params.List {
						// Each field might have multiple names
						typeStr := exprToString(field.Type)
						if len(field.Names) == 0 {
							// Anonymous parameter (e.g., context.Context)
							params = append(params, typeStr)
						} else {
							for _, name := range field.Names {
								params = append(params, fmt.Sprintf("%s %s", name.Name, typeStr))
							}
						}
					}
				}
				paramStr := strings.Join(params, ", ")

				var results []string
				if d.Type.Results != nil {
					for _, field := range d.Type.Results.List {
						results = append(results, exprToString(field.Type))
					}
				}
				resultStr := ""
				if len(results) > 0 {
					resultStr = fmt.Sprintf(" -> %s", strings.Join(results, ", "))
				}

				outputBuffer.WriteString(fmt.Sprintf("  Function: %s(%s)%s\n", d.Name.Name, paramStr, resultStr))

			case *ast.GenDecl:
				// Look for type declarations (structs and interfaces)
				if d.Tok == token.TYPE {
					for _, spec := range d.Specs {
						ts, ok := spec.(*ast.TypeSpec)
						if !ok {
							continue
						}
						switch ts.Type.(type) {
						case *ast.StructType:
							outputBuffer.WriteString(fmt.Sprintf("  Struct: %s\n", ts.Name.Name))
						case *ast.InterfaceType:
							outputBuffer.WriteString(fmt.Sprintf("  Interface: %s\n", ts.Name.Name))
						}
					}
				}
			}
		}
		outputBuffer.WriteString("\n")
		return nil
	})
	if err != nil {
		return err
	}

	// Copy the collected output to the clipboard
	if err := clipboard.WriteAll(outputBuffer.String()); err != nil {
		return fmt.Errorf("failed to copy to clipboard: %v", err)
	}

	fmt.Println("GrabSummary copied to clipboard.")
	return nil
}
