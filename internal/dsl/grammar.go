package dsl

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2"
)

type Program struct {
	Header     *Header      `@@`
	Statements []*Statement `{ @@ }`
}

type Header struct {
	Repo   string `"#!" @Ident ":"`
	Branch string `@Ident`
}

type Statement struct {
	Let     *LetStmt     ` @@`
	Command *CommandStmt `| @@`
	For     *ForStmt     `| @@`
	If      *IfStmt      `| @@`
	Expr    *ExprStmt    `| @@`
}

type ExprStmt struct {
	Expr *Expr `@@`
}

type LetStmt struct {
	Let   string `"let"`
	Name  string `@Ident`
	Eq    string `"="`
	Value *Expr  `@@`
}

type CommandStmt struct {
	Name    string           `@Ident`
	Options []*CommandOption `@@*`
}

type CommandOption struct {
	Dot   string `"."`
	Name  string `@Ident`
	Value *Expr  `[ @@ ]`
}

type ForStmt struct {
	For   string       `"for"`
	Var   string       `@Ident`
	In    string       `"in"`
	Range *Expr        `@@`
	Body  []*Statement `"{" { @@ } "}"`
}

type IfStmt struct {
	If   string       `"if"`
	Cond *Expr        `@@`
	Then []*Statement `"{" { @@ } "}"`
	Else *ElseStmt    `[ @@ ]`
}

type ElseStmt struct {
	// Expr *Expr `@@`
	Else string       `"else"`
	Body []*Statement `"{" { @@ } "}"`
}

type Expr struct {
	Ident  *string     ` @Ident`
	String *string     `| @String`
	Number *float64    `| @Int`
	Semver *Semver     `| @@`
	Bool   *BoolLit    `| @@`
	Array  []*Expr     `| "[" [ @@ { "," @@ } ] "]"`
	Call   *CallExpr   `| @@`
	Binary *BinaryExpr `| @@`
}

type BoolLit struct {
	Value string `@("true" | "false")`
}

type Semver struct {
	V     *string `[@"v"]`
	Major int     `@Int`
	Dot1  string  `"."`
	Minor string  `@Int`
	Dot2  string  `"."`
	Patch string  `@Int`
	Pre   *string `["-" @Ident ]`
	Build *string `["+" @Ident ]`
}

type CallExpr struct {
	Name string  `@Ident`
	LPar string  `"("`
	Args []*Expr `[ @@ { "," @@ } ]`
	RPar string  `")"`
}

type BinaryExpr struct {
	Left     *Expr  `@@`
	Operator string `@("==" | "!=" | "<=" | ">=" | "<" | ">" | "+" | "-" | "*" | "/" | "&&" | "||")`
	Right    *Expr
}

func (s *Semver) String() string {
	var b strings.Builder
	if s.V != nil {
		b.WriteString(*s.V)
	}
	b.WriteString(fmt.Sprintf("%d.%d.%d", s.Major, s.Minor, s.Patch))
	if s.Pre != nil {
		b.WriteString("-")
		b.WriteString(*s.Build)
	}
	if s.Build != nil {
		b.WriteString("+")
		b.WriteString(*s.Build)
	}

	return b.String()
}

var parser = participle.MustBuild[Program]()

func ParseScript(input string) (*Program, error) {
	return parser.ParseString("", input)
}
