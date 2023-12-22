package _ast

type Node interface {
	Symbol()
}

func (*FuncDecl) Symbol() {}
func (*Ident) Symbol()    {}
func (*Block) Symbol()    {}

type FuncDecl struct {
	Block    Block
	FuncName Ident
}

type Ident struct {
	Value string
	Pos   int
	End   int
}

type Block struct {
	LBrace int
	RBrace int
}
