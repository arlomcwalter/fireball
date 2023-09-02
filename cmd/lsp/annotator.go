package lsp

import (
	"fireball/core"
	"fireball/core/ast"
	"fireball/core/scanner"
	"fireball/core/types"
	"fmt"
	"github.com/MineGame159/protocol"
)

type annotator struct {
	functions map[string]*ast.Func
	hints     []protocol.InlayHint
}

func annotate(decls []ast.Decl) []protocol.InlayHint {
	a := annotator{
		functions: make(map[string]*ast.Func),
		hints:     make([]protocol.InlayHint, 0, 16),
	}

	for _, decl := range decls {
		if f, ok := decl.(*ast.Func); ok {
			a.functions[f.Name.Lexeme] = f
		}
	}

	for _, decl := range decls {
		a.AcceptDecl(decl)
	}

	return a.hints
}

// Declarations

func (a *annotator) VisitStruct(decl *ast.Struct) {
	decl.AcceptChildren(a)
}

func (a *annotator) VisitEnum(decl *ast.Enum) {
	if decl.InferType {
		a.addToken(decl.Name, " "+decl.Type.String(), protocol.InlayHintKindType)
	}

	for _, case_ := range decl.Cases {
		if case_.InferValue {
			a.addToken(case_.Name, fmt.Sprintf(" = %d", case_.Value), protocol.InlayHintKindParameter)
		}
	}

	decl.AcceptChildren(a)
}

func (a *annotator) VisitFunc(decl *ast.Func) {
	decl.AcceptChildren(a)
}

// Statements

func (a *annotator) VisitBlock(stmt *ast.Block) {
	stmt.AcceptChildren(a)
}

func (a *annotator) VisitExpression(stmt *ast.Expression) {
	stmt.AcceptChildren(a)
}

func (a *annotator) VisitVariable(stmt *ast.Variable) {
	if stmt.InferType {
		a.addToken(stmt.Name, " "+stmt.Type.String(), protocol.InlayHintKindType)
	}

	stmt.AcceptChildren(a)
}

func (a *annotator) VisitIf(stmt *ast.If) {
	stmt.AcceptChildren(a)
}

func (a *annotator) VisitFor(stmt *ast.For) {
	stmt.AcceptChildren(a)
}

func (a *annotator) VisitReturn(stmt *ast.Return) {
	stmt.AcceptChildren(a)
}

func (a *annotator) VisitBreak(stmt *ast.Break) {
	stmt.AcceptChildren(a)
}

func (a *annotator) VisitContinue(stmt *ast.Continue) {
	stmt.AcceptChildren(a)
}

// Expressions

func (a *annotator) VisitGroup(expr *ast.Group) {
	expr.AcceptChildren(a)
}

func (a *annotator) VisitLiteral(expr *ast.Literal) {
	expr.AcceptChildren(a)
}

func (a *annotator) VisitInitializer(expr *ast.Initializer) {
	expr.AcceptChildren(a)
}

func (a *annotator) VisitUnary(expr *ast.Unary) {
	expr.AcceptChildren(a)
}

func (a *annotator) VisitBinary(expr *ast.Binary) {
	expr.AcceptChildren(a)
}

func (a *annotator) VisitLogical(expr *ast.Logical) {
	expr.AcceptChildren(a)
}

func (a *annotator) VisitIdentifier(expr *ast.Identifier) {
	expr.AcceptChildren(a)
}

func (a *annotator) VisitAssignment(expr *ast.Assignment) {
	expr.AcceptChildren(a)
}

func (a *annotator) VisitCast(expr *ast.Cast) {
	expr.AcceptChildren(a)
}

func (a *annotator) VisitCall(expr *ast.Call) {
	if false {
		if _, ok := expr.Callee.Type().(*types.FunctionType); ok {
			if ident, ok := expr.Callee.(*ast.Identifier); ok {
				if f, ok := a.functions[ident.Identifier.Lexeme]; ok {
					for i, arg := range expr.Args {
						if i >= len(f.Params) {
							break
						}

						param := f.Params[i]
						a.add(arg.Range().Start, param.Name.Lexeme+": ", protocol.InlayHintKindParameter)
					}
				}
			}
		}
	}

	expr.AcceptChildren(a)
}

func (a *annotator) VisitIndex(expr *ast.Index) {
	expr.AcceptChildren(a)
}

func (a *annotator) VisitMember(expr *ast.Member) {
	expr.AcceptChildren(a)
}

// ast.Acceptor

func (a *annotator) AcceptDecl(decl ast.Decl) {
	decl.Accept(a)
}

func (a *annotator) AcceptStmt(stmt ast.Stmt) {
	stmt.Accept(a)
}

func (a *annotator) AcceptExpr(expr ast.Expr) {
	expr.Accept(a)
}

// Utils

func (a *annotator) addToken(token scanner.Token, text string, kind protocol.InlayHintKind) {
	a.add(core.TokenToPos(token, true), text, kind)
}

func (a *annotator) add(pos core.Pos, text string, kind protocol.InlayHintKind) {
	a.hints = append(a.hints, protocol.InlayHint{
		Position: convertPos(pos),
		Label:    text,
		Kind:     kind,
	})
}
