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
	switch rbrace := Tokens[Scan].(type) {
	case _lex.RBrace:
		block.RBrace = rbrace.Pos
		return &block, nil
	}
	fmt.Println("Error parsing block")
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
