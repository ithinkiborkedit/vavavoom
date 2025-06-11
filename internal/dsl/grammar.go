package dsl

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type Program struct {
	Header     *Header      `@@`
	Statements []*Statement `{ @@ }`
}

type Header struct {
	Bang   string `@HeaderBang`
	Repo   string `@RepoURL`
	Colon  string `@Colon`
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
	Number *float64    `| @Number`
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
	Major int     `@Number`
	Dot1  string  `"."`
	Minor string  `@Number`
	Dot2  string  `"."`
	Patch string  `@Number`
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

var DslLexer = lexer.MustSimple([]lexer.SimpleRule{
	{"HeaderBang", `#!`},
	{"RepoURL", `[^:\s]+`}, // repo: any non-colon, non-whitespace string (includes URLs)
	{"Colon", `:`},
	{"Let", `let\b`},
	{"For", `for\b`},
	{"In", `in\b`},
	{"If", `if\b`},
	{"Else", `else\b`},
	{"True", `true\b`},
	{"False", `false\b`},
	{"Number", `[-+]?\d*\.?\d+([eE][-+]?\d+)?`},
	{"String", `"(?:\\.|[^"])*"`},
	{"Ident", `[a-zA-Z_][a-zA-Z0-9_]*`},
	{"LBrace", `\{`},
	{"RBrace", `\}`},
	{"LBracket", `\[`},
	{"RBracket", `\]`},
	{"LParen", `\(`},
	{"RParen", `\)`},
	{"Comma", `,`},
	{"Assign", `=`},
	{"Semicolon", `;`},
	// Operators (add more as needed)
	{"Op", `==|!=|<=|>=|&&|\|\||[+\-*/<>]`},
	{"Dot", `\.`},
	{"Whitespace", `[ \t\n\r]+`},
	{"Comment", `//[^\n]*`},
})

var parser = participle.MustBuild[Program](
	participle.Lexer(DslLexer),
)

func ParseScript(input string) (*Program, error) {
	return parser.ParseString("", input)
}
