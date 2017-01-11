package jssquish

import (
	"io"
	"log"

	"github.com/robertkrimen/otto/ast"
	"github.com/robertkrimen/otto/parser"
)

func ParseRequires(r io.Reader, path string) ([]string, error) {
	program, err := parser.ParseFile(nil, path, r, parser.IgnoreRegExpErrors)
	if err != nil {
		return nil, err
	}

	visitor := NewRequireVisitor()
	for _, stmt := range program.Body {
		if err := WalkNode(visitor, stmt); err != nil {
			return nil, err
		}
	}
	return visitor.Requires(), nil
}

type RequireVisitor struct {
	requires map[string]bool
}

func NewRequireVisitor() *RequireVisitor {
	return &RequireVisitor{
		requires: make(map[string]bool),
	}
}

func (rv *RequireVisitor) Visit(n ast.Node) bool {
	if ce, ok := n.(*ast.CallExpression); ok {
		return rv.visitCallExpression(ce)
	}

	return true
}

func (rv *RequireVisitor) visitCallExpression(ce *ast.CallExpression) bool {
	if callee, ok := ce.Callee.(*ast.Identifier); ok && callee.Name == "require" {
		args := ce.ArgumentList

		// If encountering a `require` call with more than one argument, log a
		// message and bail without descending.
		if len(args) != 1 {
			log.Printf("require statement found with >1 arguments.")
			return false
		}

		// When encountering a non-string argument, log a message and bail without
		// descending.
		if str, ok := args[0].(*ast.StringLiteral); !ok {
			log.Printf("require statement with non-string argument found")
			return false
		} else {
			rv.requires[str.Value] = true
		}
	}
	return true
}

func (rv *RequireVisitor) Requires() []string {
	requires := make([]string, 0, len(rv.requires))
	for k := range rv.requires {
		requires = append(requires, k)
	}
	return requires
}
