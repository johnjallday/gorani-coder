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
)

// exprToString converts an AST expression into its string representation.
func exprToString(expr ast.Expr) string {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, token.NewFileSet(), expr); err != nil {
		return ""
	}
	return buf.String()
}

// GrabSummary walks through Go files under the provided root directory,
// collects and prints detailed symbol info including functions, structs, and interfaces,
// along with the package each symbol is declared in.
func GrabSummary(root string) error {
	fmt.Println("Grabbing summary from:", root)
	fset := token.NewFileSet()

	// Walk through all files in the directory recursively.
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err // Propagate error from Walk
		}
		// Only process .go files.
		if info.IsDir() || filepath.Ext(path) != ".go" {
			return nil
		}

		// Parse the Go file.
		file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			fmt.Printf("Error parsing file %s: %v\n", path, err)
			return nil // Skip files with parsing errors.
		}

		// Retrieve package name from the file.
		packageName := file.Name.Name
		fmt.Printf("\nFile: %s (package %s)\n", path, packageName)

		// Walk the AST to find functions, structs, and interfaces.
		ast.Inspect(file, func(n ast.Node) bool {
			switch node := n.(type) {
			// Function declarations (including methods)
			case *ast.FuncDecl:
				recv := ""
				if node.Recv != nil && len(node.Recv.List) > 0 {
					recv = fmt.Sprintf("(%s) ", exprToString(node.Recv.List[0].Type))
				}
				fmt.Printf("  [Package: %s] Function: %s%s\n", packageName, recv, node.Name.Name)

			// Type declarations (could be struct or interface)
			case *ast.GenDecl:
				if node.Tok == token.TYPE {
					for _, spec := range node.Specs {
						typeSpec, ok := spec.(*ast.TypeSpec)
						if !ok {
							continue
						}
						switch typeSpec.Type.(type) {
						case *ast.StructType:
							fmt.Printf("  [Package: %s] Struct: %s\n", packageName, typeSpec.Name.Name)
						case *ast.InterfaceType:
							fmt.Printf("  [Package: %s] Interface: %s\n", packageName, typeSpec.Name.Name)
						}
					}
				}
			}
			return true
		})

		return nil
	})
	if err != nil {
		return fmt.Errorf("error walking the path %s: %v", root, err)
	}
	return nil
}
