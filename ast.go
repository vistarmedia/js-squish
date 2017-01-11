package jssquish

import (
	"fmt"

	"github.com/robertkrimen/otto/ast"
)

type Visitor interface {
	Visit(n ast.Node) bool
}

// Walks the AST of a parsed Javascript file. At each node, this will call the
// `Vist` method of a given `Visitor`. If it returns true, this will walk
// further down the AST. If false, this will continue to the next sibling.
func WalkNode(v Visitor, node ast.Node) error {
	return walkNode(v, node, 0)
}

// This is a noisy function. Basically, it's a lot of metadata saying which
// parts of an AST node can be descended into described as a giant switch
// statement.
func walkNode(v Visitor, node ast.Node, depth int) error {
	if !v.Visit(node) {
		return nil
	}

	next := func(n ast.Node) error {
		return walkNode(v, n, depth+1)
	}

	nextIfNotNil := func(n ast.Node) error {
		if n != nil {
			return next(n)
		}
		return nil
	}

	switch t := node.(type) {

	default:
		return fmt.Errorf("ast.WalkNode can't handle %T: %#v", node, node)

	case *ast.ArrayLiteral:
		for _, expr := range t.Value {
			if err := next(expr); err != nil {
				return nil
			}
		}

	case *ast.AssignExpression:
		if err := next(t.Left); err != nil {
			return err
		}
		return next(t.Right)

	// TODO: BadExpression
	// TODO: BadStatement

	case *ast.BinaryExpression:
		if err := next(t.Left); err != nil {
			return nil
		}
		return next(t.Right)

	case *ast.BlockStatement:
		for _, stmt := range t.List {
			if err := next(stmt); err != nil {
				return err
			}
		}

	case *ast.BooleanLiteral:

	case *ast.BracketExpression:
		if err := next(t.Left); err != nil {
			return err
		}
		return next(t.Member)

	case *ast.BranchStatement:
		return nextIfNotNil(t.Label)

	case *ast.CallExpression:
		if err := next(t.Callee); err != nil {
			return err
		}
		for _, expr := range t.ArgumentList {
			if err := next(expr); err != nil {
				return err
			}
		}

	case *ast.CaseStatement:
		if err := next(t.Test); err != nil {
			return nil
		}
		for _, stmt := range t.Consequent {
			if err := next(stmt); err != nil {
				return err
			}
		}

	case *ast.CatchStatement:
		if err := nextIfNotNil(t.Parameter); err != nil {
			return err
		}
		return next(t.Body)

	// TODO: Comments?

	case *ast.ConditionalExpression:
		if err := next(t.Test); err != nil {
			return err
		}
		if err := next(t.Consequent); err != nil {
			return err
		}
		return next(t.Alternate)

	case *ast.DebuggerStatement:

	case *ast.DoWhileStatement:
		if err := next(t.Test); err != nil {
			return err
		}
		return next(t.Body)

	case *ast.DotExpression:
		if err := next(t.Left); err != nil {
			return err
		}
		return nextIfNotNil(t.Identifier)

	case *ast.EmptyExpression:

	case *ast.EmptyStatement:

	case *ast.ExpressionStatement:
		return next(t.Expression)

	case *ast.ForInStatement:
		if err := next(t.Into); err != nil {
			return nil
		}
		if err := next(t.Source); err != nil {
			return nil
		}
		return next(t.Body)

	case *ast.ForStatement:
		if err := next(t.Initializer); err != nil {
			return nil
		}
		if err := next(t.Update); err != nil {
			return nil
		}
		if err := next(t.Test); err != nil {
			return nil
		}
		return next(t.Body)

	case *ast.FunctionLiteral:
		if err := nextIfNotNil(t.Name); err != nil {
			return err
		}
		return next(t.Body)

	case *ast.FunctionStatement:
		return next(t.Function)

	case *ast.Identifier:

	case *ast.IfStatement:
		if err := next(t.Test); err != nil {
			return err
		}
		if err := next(t.Consequent); err != nil {
			return err
		}
		return nextIfNotNil(t.Alternate)

	case *ast.LabelledStatement:
		return next(t.Statement)

	case *ast.NewExpression:
		if err := next(t.Callee); err != nil {
			return err
		}
		for _, expr := range t.ArgumentList {
			if err := next(expr); err != nil {
				return err
			}
		}

	case *ast.NullLiteral:

	case *ast.NumberLiteral:

	case *ast.ObjectLiteral:
		for _, prop := range t.Value {
			if err := next(prop.Value); err != nil {
				return err
			}
		}

	case *ast.RegExpLiteral:

	case *ast.ReturnStatement:
		return nextIfNotNil(t.Argument)

	case *ast.SequenceExpression:
		for _, expr := range t.Sequence {
			if err := next(expr); err != nil {
				return err
			}
		}

	case *ast.StringLiteral:

	case *ast.SwitchStatement:
		if err := next(t.Discriminant); err != nil {
			return err
		}
		for _, stmt := range t.Body {
			if err := next(stmt); err != nil {
				return err
			}
		}

	case *ast.ThisExpression:

	case *ast.ThrowStatement:
		return next(t.Argument)

	case *ast.TryStatement:
		if err := next(t.Body); err != nil {
			return err
		}
		if t.Catch != nil {
			if err := next(t.Catch); err != nil {
				return err
			}
		}
		return nextIfNotNil(t.Finally)

	case *ast.UnaryExpression:
		return next(t.Operand)

	case *ast.VariableExpression:
		return nextIfNotNil(t.Initializer)

	case *ast.VariableStatement:
		for _, child := range t.List {
			if err := next(child); err != nil {
				return err
			}
		}

	case *ast.WhileStatement:
		if err := next(t.Test); err != nil {
			return err
		}
		return next(t.Body)

	case *ast.WithStatement:
		if err := next(t.Object); err != nil {
			return err
		}
		return next(t.Body)

	}

	return nil
}
