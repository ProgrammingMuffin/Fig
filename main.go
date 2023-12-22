package main

import (
	_lex "github.com/ProgrammingMuffin/Fig/lex"
	_parse "github.com/ProgrammingMuffin/Fig/parse"
)

func main() {
	tokens := _lex.LexSourceFile("main.fig")
	_ = _parse.ParseTokens(tokens)
}
