package _gen

import (
	"strings"

	_ast "github.com/ProgrammingMuffin/Fig/ast"
)

var Module strings.Builder

func VisitTree(node _ast.Node) string {
	switch x := node.(type) {
	case *_ast.FuncDecl:
		Module.WriteString("void " + x.FuncName.Value + "()")
		VisitBlock(x.Block)
	}
	return Module.String()
}

func VisitBlock(node _ast.Block) {
	Module.WriteString("{")
	for _, stmt := range node.Stmts {
		VisitStatement(stmt)
		Module.WriteString(";")
	}
	Module.WriteString("}\n")
}

func VisitStatement(node _ast.Stmt) {
	switch x := node.(type) {
	case *_ast.AssignStmt:
		Module.WriteString("int " + x.Lhs.Value + " = ")
		VisitStatement(x.Rhs)
	case *_ast.BinaryExpr:
		switch y := x.Lhs.(type) {
		case *_ast.Ident:
			Module.WriteString(y.Value + " " + x.Op.Type + " ")
		case *_ast.BasicLit:
			Module.WriteString(y.Value + " " + x.Op.Type + " ")
		case *_ast.BinaryExpr:
			VisitStatement(y)
		case *_ast.Term:
			Module.WriteString("( ")
			switch z := y.ExprStmt.(type) {
			case *_ast.BinaryExpr:
				VisitStatement(z)
			}
			Module.WriteString(" ) " + x.Op.Type + " ")
		}
		switch y := x.Rhs.(type) {
		case *_ast.Ident:
			Module.WriteString(y.Value)
		case *_ast.BasicLit:
			Module.WriteString(y.Value)
		case *_ast.BinaryExpr:
			VisitStatement(y)
		case *_ast.Term:
			Module.WriteString("( ")
			switch z := y.ExprStmt.(type) {
			case *_ast.BinaryExpr:
				VisitStatement(z)
			}
			Module.WriteString(" )")
		}
	case *_ast.Term:
		Module.WriteString("( ")
		switch z := x.ExprStmt.(type) {
		case *_ast.BinaryExpr:
			VisitStatement(z)
		}
	}
}
