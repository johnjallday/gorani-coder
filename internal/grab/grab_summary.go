package grab

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
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
// collects detailed symbol info (including file, package, functions, structs, and interfaces),
func GrabSummary(root string) error {
	fmt.Println("Grabbing summary from:", root)
	return nil
}
