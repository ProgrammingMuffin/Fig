package _lex

import (
	"fmt"
	"os"
	"strings"
	"unicode"
)

type Token interface {
	GetPos() (int, int)
}

type LBrace struct {
	Token
	Pos int
	End int
}

type RBrace struct {
	Token
	Pos int
	End int
}

type Keyword struct {
	Token,
	Value string
	Pos int
	End int
}

type Ident struct {
	Token
	Pos   int
	End   int
	Value string
}

type Number struct {
	Token
	Pos int
	End int
}

func (x LBrace) GetPos() (int, int) {
	return x.Pos, x.End
}

func (x RBrace) GetPos() (int, int) {
	return x.Pos, x.End
}

func (x Keyword) GetPos() (int, int) {
	return x.Pos, x.End
}

func (x Ident) GetPos() (int, int) {
	return x.Pos, x.End
}

func (x Number) GetPos() (int, int) {
	return x.Pos, x.End
}

func LexSourceFile(sourceFile string) []Token {
	tokens := []Token{}
	data, err := os.ReadFile(sourceFile)
	if err != nil {
		fmt.Println("Error reading file")
	}
	var value strings.Builder
	var numberValue strings.Builder
	var symbol strings.Builder
	for index, char := range data {
		if unicode.IsLetter(rune(char)) || charIn(char, '_') {
			value.Write([]byte{char})
		} else if unicode.IsNumber(rune(char)) {
			if index == 0 {
				numberValue.Write([]byte{char})
			} else {
				if unicode.IsNumber(rune(data[index-1])) || (!unicode.IsLetter(rune(data[index-1])) && !unicode.IsNumber(rune(data[index-1])) && data[index-1] != '_') {
					numberValue.Write([]byte{char})
				} else {
					value.Write([]byte{char})
				}
			}
		} else {
			if isSymbol(char) {
				symbol.WriteByte(char)
			}
		}
		if index < (len(data) - 1) {
			if !isSymbol(char) && !unicode.IsLetter(rune(data[index+1])) && !unicode.IsNumber(rune(data[index+1])) && data[index+1] != '_' {
				if value.Len() != 0 {
					pos := index - len(value.String()) + 1
					end := index
					switch value.String() {
					case "func":
						tokens = append(tokens, Keyword{Pos: pos, End: end, Value: "func"})
					default:
						tokens = append(tokens, Ident{Pos: pos, End: end, Value: value.String()})
					}
					value.Reset()
				}
			} else if isSymbol(char) && (unicode.IsLetter(rune(data[index+1])) || unicode.IsNumber(rune(data[index+1])) || data[index+1] == ' ' || data[index+1] == '\t' || data[index+1] == '\n') && !isSymbol(data[index+1]) {
				if symbol.Len() != 0 {
					pos := index - len(symbol.String())
					end := index
					switch symbol.String() {
					case "{":
						tokens = append(tokens, LBrace{Pos: pos, End: end})
					case "}":
						tokens = append(tokens, RBrace{Pos: pos, End: end})
					}
					symbol.Reset()
				}
			}
		}
	}
	fmt.Println("The tokens are: ", tokens)
	return tokens
}

func charIn(char byte, chars ...byte) bool {
	for _, character := range chars {
		if character == char {
			return true
		}
	}
	return false
}

func isSymbol(char byte) bool {
	return charIn(char, '{', '}', '(', ')', ';')
}
