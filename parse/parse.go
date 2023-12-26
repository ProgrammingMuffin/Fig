package _parse

import (
	"errors"
	"fmt"
	"os"

	_ast "github.com/ProgrammingMuffin/Fig/ast"
	_lex "github.com/ProgrammingMuffin/Fig/lex"
)

var Scan int = 0
var Tokens []_lex.Token

func ParseTokens(tokens []_lex.Token) []_ast.Node {
	nodes := []_ast.Node{}
	Tokens = tokens
	for {
		switch tokenType := Tokens[Scan].(type) {
		case _lex.Keyword:
			if tokenType.Value == "func" {
				node, err := ParseFuncDecl()
				if err != nil {
					goto FINISH
				}
				nodes = append(nodes, node)
				fmt.Println("successfully parsed function")
			}
		}
		if SafeInc() != nil {
			goto FINISH
		} else {
			continue
		}
	FINISH:
		fmt.Println("finished parsing")
		break
	}
	return nodes
}

func ParseFuncDecl() (*_ast.FuncDecl, error) {
	funcDecl := _ast.FuncDecl{}
	err := SafeInc()
	if err != nil {
		return nil, err
	}
	switch ident := Tokens[Scan].(type) {
	case _lex.Ident:
		funcDecl.FuncName = _ast.Ident{
			Value: ident.Value,
			Pos:   ident.Pos,
			End:   ident.End,
		}
	}
	err = SafeInc()
	if err != nil {
		return nil, err
	}
	switch Tokens[Scan].(type) {
	case _lex.LBrace:
		block, err := ParseBlock()
		if err != nil {
			return nil, err
		}
		funcDecl.Block = *block
		return &funcDecl, nil
	}
	fmt.Println("error parsing function")
	os.Exit(0)
	return nil, nil
}

func ParseBlock() (*_ast.Block, error) {
	pos, _ := Tokens[Scan].GetPos()
	block := _ast.Block{LBrace: pos}
	err := SafeInc()
	if err != nil {
		return nil, err
	}
	switch x := Tokens[Scan].(type) {
	case _lex.Ident:
		stmts, err := ParseStatements()
		if err != nil {
			return nil, err
		}
		block.Stmts = stmts
	case _lex.RBrace:
		block.RBrace = x.Pos
		return &block, nil
	}
	switch x := Tokens[Scan].(type) {
	case _lex.RBrace:
		block.RBrace = x.Pos
		return &block, nil
	}
	fmt.Println("Error parsing block")
	os.Exit(0)
	return nil, nil
}

func ParseStatements() ([]_ast.Stmt, error) {
	stmts := []_ast.Stmt{}
	for {
		switch Tokens[Scan].(type) {
		case _lex.RBrace:
			goto OUTSIDE
		}
		stmt, err := ParseStatement()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			stmts = append(stmts, stmt)
		} else {
			if err != nil {
				stmt, err = ParseStatement()
				stmts = append(stmts, stmt)
			}
			goto OUTSIDE
		}
	}
OUTSIDE:
	return stmts, nil
}

func ParseStatement() (_ast.Stmt, error) {
	switch Tokens[Scan].(type) {
	case _lex.LParen:
		term, err := ParseTerm()
		if err != nil {
			return nil, err
		}
		expr, err := ParseStatement()
		switch x := expr.(type) {
		case *_ast.BinaryExpr:
			x.Lhs = term
		}
		return expr, nil
	}
	if (Scan + 1) < len(Tokens) {
		switch x := Tokens[Scan+1].(type) {
		case _lex.Operator:
			if x.Kind == "=" {
				return ParseAssignStatement()
			} else if x.Kind == "+" || x.Kind == "-" || x.Kind == "*" || x.Kind == "/" {
				stmt, err := ParseBinaryExpr(nil, nil, 1)
				return stmt, err
			}
		default:
			return nil, nil
		}
	} else {
		return nil, nil
	}
	fmt.Println("error parsing statement")
	os.Exit(0)
	return nil, nil
}

func ParseAssignStatement() (_ast.Stmt, error) {
	stmt := _ast.AssignStmt{}
	var lhs _ast.Ident
	switch x := Tokens[Scan].(type) {
	case _lex.Ident:
		lhs.Value = x.Value
		lhs.Pos = x.Pos
		lhs.End = x.End
	}
	stmt.Lhs = lhs
	SafeInc()
	err := SafeInc()
	if err != nil {
		return nil, err
	}
	switch Tokens[Scan].(type) {
	case _lex.Ident:
		nextStmt, err := ParseStatement()
		if err != nil {
			return nil, err
		}
		stmt.Rhs = nextStmt
	case _lex.LParen:
		term, err := ParseTerm()
		if err != nil {
			return nil, err
		}
		_, err = ParseBinaryExpr(nil, term, 0)
		if err != nil {
			return nil, err
		}
		stmt.Rhs = term
	case _lex.Number:
		nextStmt, err := ParseStatement()
		if err != nil {
			return nil, err
		}
		stmt.Rhs = nextStmt
	}
	SafeInc()
	switch Tokens[Scan].(type) {
	case _lex.Semicolon:
		SafeInc()
		return &stmt, nil
	}
	fmt.Println("error parsing assignment statement")
	return nil, nil
}

func ParseTerm() (*_ast.Term, error) {
	term := _ast.Term{}
	err := SafeInc()
	if err != nil {
		return nil, err
	}
	stmt, err := ParseStatement()
	if err != nil {
		return nil, err
	}
	SafeInc()
	switch Tokens[Scan].(type) {
	case _lex.RParen:
		term.ExprStmt = stmt
		return &term, nil
	default:
		return nil, nil
	}
}

func ParseBinaryExpr(prev *_ast.BinaryExpr, prevTerm _ast.ExprStmt, prec int) (*_ast.BinaryExpr, error) {
	binaryExpr := _ast.BinaryExpr{}
	if Scan+2 < len(Tokens) {
		switch x := Tokens[Scan+1].(type) {
		case _lex.Semicolon:
			return nil, nil
		case _lex.RBrace:
			return nil, nil
		case _lex.RParen:
			return nil, nil
		case _lex.Operator:
			fmt.Println("operator is: ", x)
			prec := 1
			if x.Kind == "*" || x.Kind == "-" {
				prec = 2
			}
			switch x := Tokens[Scan].(type) {
			case _lex.Ident:
				binaryExpr.Lhs = &_ast.Ident{Value: x.Value, Pos: x.Pos, End: x.End}
			case _lex.Number:
				binaryExpr.Lhs = &_ast.BasicLit{
					Value: x.Value,
					Kind:  x.Kind,
				}
			}
			if prev != nil {
				prev.Rhs = &binaryExpr
			} else if prevTerm != nil {
				binaryExpr.Lhs = prevTerm
			}
			binaryExpr.Op = _ast.Operator{Type: x.Kind}
			SafeInc()
			SafeInc()
			switch Tokens[Scan].(type) {
			case _lex.LParen:
				term, err := ParseTerm()
				if err != nil {
					return nil, err
				}
				binaryExpr.Rhs = term
			}
			expr, err := ParseBinaryExpr(&binaryExpr, nil, prec)
			if err != nil {
				return nil, err
			}
			if expr != nil {
				binaryExpr.Rhs = expr
			} else {
				switch x := Tokens[Scan].(type) {
				case _lex.Ident:
					binaryExpr.Rhs = &_ast.Ident{Value: x.Value, Pos: x.Pos, End: x.End}
				case _lex.Number:
					binaryExpr.Rhs = &_ast.BasicLit{Value: x.Value, Kind: x.Kind}
				}
			}
			return &binaryExpr, nil
		}
	}
	fmt.Println("Error parsing Binary expression")
	os.Exit(0)
	return nil, nil
}

func SafeInc() error {
	Scan++
	if Scan >= len(Tokens) {
		return errors.New("EOF")
	}
	return nil
}
