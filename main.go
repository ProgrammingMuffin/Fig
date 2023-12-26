package main

import (
	"fmt"

	_gen "github.com/ProgrammingMuffin/Fig/gen"
	_lex "github.com/ProgrammingMuffin/Fig/lex"
	_parse "github.com/ProgrammingMuffin/Fig/parse"
	_semantic "github.com/ProgrammingMuffin/Fig/semantic"
	"github.com/k0kubun/pp/v3"
)

func main() {
	tokens := _lex.LexSourceFile("main.fig")
	ast := _parse.ParseTokens(tokens)
	for _, astVal := range ast {
		_semantic.VisitNode(astVal)
	}
	pp.Println(ast)
	for _, astVal := range ast {
		_gen.VisitTree(astVal)
	}
	fmt.Println(_gen.Module.String())
}
