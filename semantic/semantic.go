package _semantic

import (
	"fmt"
	"os"

	_ast "github.com/ProgrammingMuffin/Fig/ast"
	"github.com/k0kubun/pp"
)

var Symtab map[string]string = make(map[string]string) //matches symbol with type names

func VisitNode(node _ast.Node) {
	switch x := node.(type) {
	case *_ast.FuncDecl:
		VisitBlock(x.Block)
	}
}

func VisitBlock(node _ast.Block) {
	for _, stmt := range node.Stmts {
		VisitStatement(stmt)
	}
}

func VisitStatement(node _ast.Stmt) string {
	switch x := node.(type) {
	case *_ast.AssignStmt:
		x.Lhs.Type = VisitStatement(x.Rhs)
		Symtab[x.Lhs.Value] = x.Lhs.Type
		return x.Lhs.Type
	case *_ast.Term:
		return VisitStatement(x.ExprStmt)
	case *_ast.BinaryExpr:
		return VisitBinaryExpression(*x)
	}
	fmt.Println("type mismatch in statement")
	return ""
}

func VisitBinaryExpression(node _ast.BinaryExpr) string {
	fmt.Println("symtab is : ")
	pp.Println(Symtab)
	type1 := ""
	type2 := ""
	switch x := node.Lhs.(type) {
	case *_ast.BasicLit:
		type1 = x.Kind
	case *_ast.BinaryExpr:
		type1 = VisitBinaryExpression(*x)
	case *_ast.Term:
		type1 = VisitStatement(x.ExprStmt)
	case *_ast.Ident:
		if val, ok := Symtab[x.Value]; ok {
			type1 = val
			x.Type = val
		}
	}
	switch x := node.Rhs.(type) {
	case *_ast.BasicLit:
		type2 = x.Kind
	case *_ast.BinaryExpr:
		type2 = VisitBinaryExpression(*x)
	case *_ast.Term:
		type2 = VisitStatement(x.ExprStmt)
	case *_ast.Ident:
		if val, ok := Symtab[x.Value]; ok {
			type2 = val
			x.Type = val
		}
	}
	if type1 == type2 && type1 != "" {
		return type1
	}
	fmt.Println("Type mismatch in binary expression")
	os.Exit(0)
	return type1
}
