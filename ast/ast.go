package _ast

type Node interface {
	Symbol()
}

func (*FuncDecl) Symbol()   {}
func (*Ident) Symbol()      {}
func (*Ident) Stmt()        {}
func (*Ident) Expr()        {}
func (*Block) Symbol()      {}
func (*BasicLit) Symbol()   {}
func (*BinaryExpr) Symbol() {}
func (*AssignStmt) Symbol() {}
func (*Operator) Symbol()   {}
func (*CallExpr) Symbol()   {}
func (*CallExpr) Stmt()     {}

type FuncDecl struct {
	Block    Block
	FuncName Ident
}

type Ident struct {
	Type  string
	Value string
	Pos   int
	End   int
}

type CallExpr struct {
	LParen int
	RParen int
	Value  Ident
	Pos    int
	End    int
}

type BasicLit struct {
	Kind  string
	Value string
}

type Stmt interface {
	Stmt()
}

type Block struct {
	LBrace int
	RBrace int
	Stmts  []Stmt
}

type Term struct {
	LParen   int
	RParen   int
	ExprStmt Stmt
}

func (*BasicLit) Expr() {}
func (*Term) Expr()     {}
func (*Term) Stmt()     {}

type ExprStmt interface {
	Expr()
}

type BinaryExpr struct {
	Lhs ExprStmt
	Op  Operator
	Rhs ExprStmt
}

func (b *BinaryExpr) Stmt() {}
func (b *BinaryExpr) Expr() {}

type AssignStmt struct {
	Lhs Ident
	Rhs Stmt
}

func (*AssignStmt) Stmt() {}

func (*TypeCast) Stmt() {}
func (*TypeCast) Expr() {}

type Operator struct {
	Type string
}

type TypeCast struct {
	Kind   string
	Value  Ident
	LParen int
	RParen int
}
