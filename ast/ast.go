package ast

import (
	"bytes"
	"fmt"
	"strings"

	"tiny-script/token"
)

// Node represents a node in Lox.
// All node types should implement the Node interface.
type Node interface {
	node()
	String() string
}

// Expr represents an expression that can be evaluated to a value.
type Expr interface {
	Node
	expr()
}

// Stmt represent a statement in Lox.
type Stmt interface {
	Node
	stmt()
}

func (*Ident) node() {}

func (*Literal) node() {}

func (*AssignExpr) node()       {}
func (*BinaryExpr) node()       {}
func (*CallExpr) node()         {}
func (*GetExpr) node()          {}
func (*GroupingExpr) node()     {}
func (*LogicalExpr) node()      {}
func (*SetExpr) node()          {}
func (*SuperExpr) node()        {}
func (*ThisExpr) node()         {}
func (*UnaryExpr) node()        {}
func (*VariableExpr) node()     {}
func (*ArrayLiteralExpr) node() {}
func (*IndexExpr) node()        {}

func (*BlockStmt) node()    {}
func (*ClassStmt) node()    {}
func (*ExprStmt) node()     {}
func (*FunctionStmt) node() {}
func (*IfStmt) node()       {}
func (*PrintStmt) node()    {}
func (*ReturnStmt) node()   {}
func (*VarStmt) node()      {}
func (*LetStmt) node()      {}
func (*WhileStmt) node()    {}
func (*ImportStmt) node()   {}

// Ident represents an identifier.
type Ident struct {
	Name string
}

func (ident *Ident) String() string { return ident.Name }

type Literal struct {
	Token token.Token
	Value string
}

func (*Literal) expr() {}

func (lit *Literal) String() string {
	switch lit.Token {
	case token.Nil:
		return "null"
	case token.True:
		return "true"
	case token.False:
		return "false"
	case token.Number, token.String:
		return lit.Value
	default:
		panic("unhandled default case")
	}
}

type (
	// AssignExpr 赋值表达式
	AssignExpr struct {
		Left  *VariableExpr
		Value Expr
	}
	// BinaryExpr 二元运算符表达式
	BinaryExpr struct {
		Left     Expr
		Operator token.Token
		Right    Expr
	}
	// CallExpr 函数调用表达式
	CallExpr struct {
		Callee    Expr
		Arguments []Expr
	}
	// GetExpr 对象的获取表达式
	GetExpr struct {
		Object Expr
		Name   string
	}
	// GroupingExpr 括号表达式
	GroupingExpr struct {
		Expression Expr
	}
	// LogicalExpr
	LogicalExpr struct {
		Left     Expr
		Operator token.Token
		Right    Expr
	}
	// SetExpr 对象的设置字段表达式
	SetExpr struct {
		Object Expr
		Name   string
		Value  Expr
	}
	// SuperExpr 暂未实现
	SuperExpr struct {
		// Method  Ident
		Keyword token.Token
		Method  token.Token
	}
	// ThisExpr this
	ThisExpr  struct{}
	UnaryExpr struct {
		Operator token.Token
		Right    Expr
	}
	// VariableExpr 定义变量表达式
	VariableExpr struct {
		Name     string
		Distance int // -1 represents global variable.
	}
	// ArrayLiteralExpr 数组字面量表达式
	ArrayLiteralExpr struct {
		Elements []Expr
		Distance int // -1 represents global variable.
	}
	// IndexExpr 数组索引表达式
	IndexExpr struct {
		Left  Expr
		Index Expr
	}
)

func (*AssignExpr) expr()       {}
func (*BinaryExpr) expr()       {}
func (*CallExpr) expr()         {}
func (*GetExpr) expr()          {}
func (*GroupingExpr) expr()     {}
func (*LogicalExpr) expr()      {}
func (*SetExpr) expr()          {}
func (*SuperExpr) expr()        {}
func (*ThisExpr) expr()         {}
func (*UnaryExpr) expr()        {}
func (*VariableExpr) expr()     {}
func (*ArrayLiteralExpr) expr() {}
func (*IndexExpr) expr()        {}

func (e *AssignExpr) String() string {
	return fmt.Sprintf("%s = %s", e.Left, e.Value)
}

func (e *BinaryExpr) String() string {
	return fmt.Sprintf("(%s %s %s)", e.Left, e.Operator, e.Right)
}

func (e *CallExpr) String() string {
	args := make([]string, len(e.Arguments))
	for i, arg := range e.Arguments {
		args[i] = arg.String()
	}
	return fmt.Sprintf("%s(%s)", e.Callee, strings.Join(args, ", "))
}

func (e *GetExpr) String() string {
	return e.Object.String() + "." + e.Name
}

func (e *GroupingExpr) String() string {
	return fmt.Sprintf("(%s)", e.Expression)
}

func (e *LogicalExpr) String() string {
	return fmt.Sprintf("%s %s %s", e.Left, e.Operator, e.Right)
}

func (e *SetExpr) String() string {
	return fmt.Sprintf("%s.%s = %s", e.Object, e.Name, e.Value)
}

func (e *SuperExpr) String() string {
	return "TODO"
}

func (e *ThisExpr) String() string {
	return "this"
}

func (e *UnaryExpr) String() string {
	return fmt.Sprintf("(%s%s)", e.Operator, e.Right)
}

func (e *VariableExpr) String() string {
	return e.Name
}

func (e *ArrayLiteralExpr) String() string {
	var arr = make([]string, 0, len(e.Elements))
	for _, v := range e.Elements {
		arr = append(arr, v.String())
	}
	buff := bytes.Buffer{}
	buff.WriteString("[")
	buff.WriteString(strings.Join(arr, ","))
	buff.WriteString("]")
	return buff.String()
}

func (e *IndexExpr) String() string {
	return fmt.Sprintf("%s[%s]", e.Left.String(), e.Index.String())
}

type (
	BlockStmt struct {
		Statements []Stmt
	}
	ClassStmt struct {
		Name       string
		SuperClass VariableExpr
		Methods    []*FunctionStmt
	}
	ImportStmt struct {
		Name        string
		Initializer *AssignExpr
	}
	ExprStmt struct {
		Expression Expr
	}
	FunctionStmt struct {
		Name          string
		Params        []*Ident
		Body          []Stmt
		IsInitializer bool
	}
	IfStmt struct {
		Condition  Expr
		ThenBranch Stmt
		ElseBranch Stmt
	}
	PrintStmt struct {
		Expression Expr
	}
	ReturnStmt struct {
		Keyword token.Token
		Value   Expr
	}
	VarStmt struct {
		Name        *Ident
		Initializer Expr
	}
	LetStmt struct {
		Name        *Ident
		Initializer Expr
	}
	WhileStmt struct {
		Condition Expr
		Body      Stmt
	}
)

func (*BlockStmt) stmt()    {}
func (*ClassStmt) stmt()    {}
func (*ExprStmt) stmt()     {}
func (*FunctionStmt) stmt() {}
func (*IfStmt) stmt()       {}
func (*PrintStmt) stmt()    {}
func (*ReturnStmt) stmt()   {}
func (*VarStmt) stmt()      {}
func (*LetStmt) stmt()      {}
func (*WhileStmt) stmt()    {}
func (*ImportStmt) stmt()   {}

func (i *ImportStmt) String() string {
	return "import" + i.Name
}

func (s *BlockStmt) String() string {
	var sb strings.Builder
	sb.WriteString("{ ")
	for _, stmt := range s.Statements {
		sb.WriteString(stmt.String())
	}
	sb.WriteString(" }")
	return sb.String()
}

func (s *ClassStmt) String() string {
	return "class " + s.Name
}

func (s *ExprStmt) String() string {
	return s.Expression.String() + ";"
}

func (s *FunctionStmt) String() string {
	var sb strings.Builder
	sb.WriteString("fun ")
	sb.WriteString(s.Name)
	sb.WriteString("(")
	params := make([]string, len(s.Params))
	for i, p := range s.Params {
		params[i] = p.Name
	}
	sb.WriteString(strings.Join(params, ", "))
	sb.WriteString(") { ")
	for _, stmt := range s.Body {
		sb.WriteString(stmt.String())
	}
	sb.WriteString(" }")
	return sb.String()
}

func (s *IfStmt) String() string {
	var sb strings.Builder
	sb.WriteString("if (")
	sb.WriteString(s.Condition.String())
	sb.WriteString(") ")
	sb.WriteString(s.ThenBranch.String())
	if s.ElseBranch != nil {
		sb.WriteString(" else ")
		sb.WriteString(s.ElseBranch.String())
	}
	return sb.String()
}

func (s *PrintStmt) String() string {
	var sb strings.Builder
	sb.WriteString("print ")
	sb.WriteString(s.Expression.String())
	sb.WriteRune(';')
	return sb.String()
}

func (s *ReturnStmt) String() string {
	str := "return"
	if s.Value != nil {
		str += " " + s.Value.String()
	}
	return str + ";"
}

func (s *VarStmt) String() string {
	var sb strings.Builder
	sb.WriteString("var ")
	sb.WriteString(s.Name.String())
	sb.WriteString(" = ")
	sb.WriteString(s.Initializer.String())
	sb.WriteRune(';')
	return sb.String()
}

func (s *LetStmt) String() string {
	var sb strings.Builder
	sb.WriteString("let ")
	sb.WriteString(s.Name.String())
	sb.WriteString(" = ")
	sb.WriteString(s.Initializer.String())
	sb.WriteRune(';')
	return sb.String()
}

func (s *WhileStmt) String() string {
	var sb strings.Builder
	sb.WriteString("while (")
	sb.WriteString(s.Condition.String())
	sb.WriteString(") ")
	sb.WriteString(s.Body.String())
	return sb.String()
}
