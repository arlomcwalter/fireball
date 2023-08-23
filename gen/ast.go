package main

import (
	"fmt"
	"os"
	"strings"
)

type item struct {
	name   string
	fields []field
	token  string
	ast    bool
}

type field struct {
	name  string
	type_ string
}

var decls = []item{
	{
		name: "Func",
		fields: []field{
			{name: "Name", type_: "Token"},
			{name: "Params", type_: "[]Param"},
			{name: "Returns", type_: "Type"},
			{name: "Body", type_: "[]Stmt"},
		},
		token: "Name",
		ast:   true,
	},
	{
		name: "Param",
		fields: []field{
			{name: "Name", type_: "Token"},
			{name: "Type", type_: "Type"},
		},
		ast: false,
	},
}

var stmts = []item{
	{
		name: "Block",
		fields: []field{
			{name: "Token_", type_: "Token"},
			{name: "Stmts", type_: "[]Stmt"},
		},
		token: "Token_",
		ast:   true,
	},
	{
		name: "Expression",
		fields: []field{
			{name: "Token_", type_: "Token"},
			{name: "Expr", type_: "Expr"},
		},
		token: "Token_",
		ast:   true,
	},
	{
		name: "Variable",
		fields: []field{
			{name: "Type", type_: "Type"},
			{name: "Name", type_: "Token"},
			{name: "Initializer", type_: "Expr"},
		},
		token: "Name",
		ast:   true,
	},
	{
		name: "If",
		fields: []field{
			{name: "Token_", type_: "Token"},
			{name: "Condition", type_: "Expr"},
			{name: "Then", type_: "Stmt"},
			{name: "Else", type_: "Stmt"},
		},
		token: "Token_",
		ast:   true,
	},
	{
		name: "Return",
		fields: []field{
			{name: "Token_", type_: "Token"},
			{name: "Expr", type_: "Expr"},
		},
		token: "Token_",
		ast:   true,
	},
}

var exprs = []item{
	{
		name: "Group",
		fields: []field{
			{name: "Token_", type_: "Token"},
			{name: "Expr", type_: "Expr"},
		},
		token: "Token_",
		ast:   true,
	},
	{
		name: "Literal",
		fields: []field{
			{name: "Value", type_: "Token"},
		},
		token: "Value",
		ast:   true,
	},
	{
		name: "Unary",
		fields: []field{
			{name: "Op", type_: "Token"},
			{name: "Right", type_: "Expr"},
		},
		token: "Op",
		ast:   true,
	},
	{
		name: "Binary",
		fields: []field{
			{name: "Left", type_: "Expr"},
			{name: "Op", type_: "Token"},
			{name: "Right", type_: "Expr"},
		},
		token: "Op",
		ast:   true,
	},
	{
		name: "Identifier",
		fields: []field{
			{name: "Identifier", type_: "Token"},
		},
		token: "Identifier",
		ast:   true,
	},
	{
		name: "Assignment",
		fields: []field{
			{name: "Assignee", type_: "Expr"},
			{name: "Op", type_: "Token"},
			{name: "Value", type_: "Expr"},
		},
		token: "Op",
		ast:   true,
	},
	{
		name: "Call",
		fields: []field{
			{name: "Token_", type_: "Token"},
			{name: "Callee", type_: "Expr"},
			{name: "Args", type_: "[]Expr"},
		},
		token: "Token_",
		ast:   true,
	},
}

func main() {
	file := os.Getenv("GOFILE")
	w := newWriter()

	switch file {
	case "declarations.go":
		generate(w, "Decl", decls)
	case "statements.go":
		generate(w, "Stmt", stmts)
	case "expressions.go":
		generate(w, "Expr", exprs)
	}

	w.flush(file)
}

func generate(w *writer, kind string, items []item) {
	w.write("package ast")
	w.write("")
	w.write("import \"fireball/core/scanner\"")
	w.write("import \"fireball/core/types\"")

	w.write("")
	w.write("//go:generate go run ../../gen/ast.go")
	w.write("")

	// Visitor
	w.write("type %sVisitor interface {", kind)

	for _, item := range items {
		if item.ast {
			w.write("Visit%s(%s *%s)", item.name, strings.ToLower(kind), item.name)
		}
	}

	w.write("}")
	w.write("")

	// Base
	w.write("type %s interface {", kind)

	w.write("Node")
	w.write("")
	w.write("Accept(visitor %sVisitor)", kind)

	if kind == "Expr" {
		w.write("")
		w.write("Type() types.Type")
		w.write("SetType(type_ types.Type)")
	}

	w.write("}")
	w.write("")

	// Items
	for _, item := range items {
		// Struct
		w.write("type %s struct {", item.name)

		if kind == "Expr" {
			w.write("type_ types.Type")
			w.write("")
		}

		for _, field := range item.fields {
			type_ := field.type_

			if type_ == "Token" {
				type_ = "scanner.Token"
			} else if type_ == "Type" {
				type_ = "types.Type"
			}

			w.write("%s %s", field.name, type_)
		}

		w.write("}")
		w.write("")

		// Node
		if item.ast {
			short := strings.ToLower(item.name)[0]
			method := fmt.Sprintf("func (%c *%s)", short, item.name)

			// Token
			w.write("%s Token() scanner.Token {", method)
			w.write("return %c.%s", short, item.token)
			w.write("}")
			w.write("")

			// Accept
			w.write("%s Accept(visitor %sVisitor) {", method, kind)
			w.write("visitor.Visit%s(%c)", item.name, short)
			w.write("}")
			w.write("")

			// Expr
			if kind == "Expr" {
				// Type
				w.write("%s Type() types.Type {", method)
				w.write("return %c.type_", short)
				w.write("}")
				w.write("")

				// SetType
				w.write("%s SetType(type_ types.Type) {", method)
				w.write("%c.type_ = type_", short)
				w.write("}")
				w.write("")
			}
		}
	}
}

type writer struct {
	str   strings.Builder
	depth int
}

func newWriter() *writer {
	return &writer{
		str:   strings.Builder{},
		depth: 0,
	}
}

func (w *writer) flush(file string) {
	_ = os.WriteFile(file, []byte(w.str.String()), 0666)
}

func (w *writer) write(format string, args ...any) {
	str := fmt.Sprintf(format, args...)

	if strings.HasPrefix(str, "}") {
		w.depth--
	}

	for i := 0; i < w.depth; i++ {
		w.str.WriteRune('\t')
	}

	w.str.WriteString(str)
	w.str.WriteRune('\n')

	if strings.HasSuffix(str, "{") {
		w.depth++
	}
}
