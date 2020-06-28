package ast

import (
	"strings"

	"github.com/pandeykartikey/goto/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}

	return ""
}

func (p *Program) String() string {
	var out strings.Builder

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}

func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}

func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode() {}

func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}

func (b *Boolean) String() string {
	return b.Token.Literal
}

type String struct {
	Token token.Token
	Value string
}

func (s *String) expressionNode() {}

func (s *String) TokenLiteral() string {
	return s.Token.Literal
}

func (s *String) String() string {
	return s.Token.Literal
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode() {}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) String() string {
	return i.Value
}

type Assignment struct {
	Token        token.Token
	NameList     *IdentifierList
	ValueList    *ExpressionList
	IsExpression bool // to check whether it acting as expression or statement
}

func (as *Assignment) statementNode() {}

func (as *Assignment) expressionNode() {}

func (as *Assignment) TokenLiteral() string {
	return as.Token.Literal
}

func (as *Assignment) String() string {
	var out strings.Builder
	if as.Token.Literal == "var" {
		out.WriteString("var ")
	}

	out.WriteString(as.NameList.String())

	if as.ValueList != nil {
		out.WriteString(" = ")
		out.WriteString(as.ValueList.String())
	}
	if !as.IsExpression {
		out.WriteString(";")
	}

	return out.String()
}

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}

func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

func (rs *ReturnStatement) String() string {
	var out strings.Builder

	out.WriteString(rs.TokenLiteral())

	if rs.ReturnValue != nil {
		out.WriteString(" ")
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}

func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}

func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}

func (pe *PrefixExpression) String() string {
	var out strings.Builder

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}

func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *InfixExpression) String() string {
	var out strings.Builder

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" ")
	out.WriteString(ie.Operator)
	out.WriteString(" ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}

func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}

func (bs *BlockStatement) String() string {
	var out strings.Builder
	out.WriteString("{ ")

	for _, stmt := range bs.Statements {
		out.WriteString(stmt.String())
	}

	out.WriteString(" }")

	return out.String()
}

type IfStatement struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
	FollowIf    *IfStatement
}

func (is *IfStatement) statementNode() {}

func (is *IfStatement) TokenLiteral() string {
	return is.Token.Literal
}

func (is *IfStatement) String() string {
	var out strings.Builder
	out.WriteString(is.TokenLiteral())
	out.WriteString(" ")
	out.WriteString(is.Condition.String())
	out.WriteString(" ")
	out.WriteString(is.Consequence.String())

	if is.FollowIf != nil {
		out.WriteString(" else ")
		out.WriteString(is.FollowIf.String())
	}

	if is.Alternative != nil {
		out.WriteString(" else ")
		out.WriteString(is.Alternative.String())
	}

	return out.String()
}

type IdentifierList struct {
	Token       token.Token
	Identifiers []*Identifier
}

func (il *IdentifierList) expressionNode() {}

func (il *IdentifierList) TokenLiteral() string {
	return il.Token.Literal
}

func (il *IdentifierList) String() string {
	var out strings.Builder

	if il != nil {

		for idx, param := range il.Identifiers {
			if idx > 0 {
				out.WriteString(",")
			}
			out.WriteString(param.String())
		}
	}

	return out.String()
}

type FuncStatement struct { // TODO: Add return type
	Token         token.Token
	Name          *Identifier
	ParameterList *IdentifierList
	FuncBody      *BlockStatement
}

func (fs *FuncStatement) statementNode() {}

func (fs *FuncStatement) TokenLiteral() string {
	return fs.Token.Literal
}

func (fs *FuncStatement) String() string {
	var out strings.Builder

	out.WriteString(fs.TokenLiteral())
	out.WriteString(" ")
	out.WriteString(fs.Name.String())
	out.WriteString(" (")
	out.WriteString(fs.ParameterList.String())
	out.WriteString(") ")
	out.WriteString(fs.FuncBody.String())

	return out.String()
}

type ExpressionList struct {
	Token       token.Token
	Expressions []*Expression
}

func (el *ExpressionList) expressionNode() {}

func (el *ExpressionList) TokenLiteral() string {
	return el.Token.Literal
}

func (el *ExpressionList) String() string {
	var out strings.Builder

	if el != nil {
		for idx, param := range el.Expressions {
			if idx > 0 {
				out.WriteString(", ")
			}
			out.WriteString((*param).String())
		}
	}

	return out.String()
}

type CallExpression struct {
	Token        token.Token
	FunctionName *Identifier
	ArgumentList *ExpressionList
}

func (ce *CallExpression) expressionNode() {}

func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.Literal
}

func (ce *CallExpression) String() string {
	var out strings.Builder

	out.WriteString(ce.FunctionName.String())
	out.WriteString("(")
	out.WriteString(ce.ArgumentList.String())
	out.WriteString(")")
	return out.String()
}

type ForStatement struct {
	Token     token.Token
	Init      *Assignment
	Condition Expression
	Update    *Assignment
	ForBody   *BlockStatement
}

func (fs *ForStatement) statementNode() {}

func (fs *ForStatement) TokenLiteral() string {
	return fs.Token.Literal
}

func (fs *ForStatement) String() string {
	var out strings.Builder

	out.WriteString(fs.TokenLiteral())
	out.WriteString(" ")
	out.WriteString(fs.Init.String())
	out.WriteString(";")
	out.WriteString(fs.Condition.String())
	out.WriteString(";")
	out.WriteString(fs.Update.String())
	out.WriteString(" ")
	out.WriteString(fs.ForBody.String())

	return out.String()
}

type LoopControlStatement struct {
	Token token.Token
	Value string
}

func (lc *LoopControlStatement) statementNode() {}

func (lc *LoopControlStatement) TokenLiteral() string {
	return lc.Token.Literal
}

func (lc *LoopControlStatement) String() string {
	return lc.Token.Literal + ";"
}

type List struct {
	Token    token.Token // the '['
	Elements *ExpressionList
}

func (l *List) expressionNode() {}

func (l *List) TokenLiteral() string {
	return l.Token.Literal
}

func (l *List) String() string {
	var out strings.Builder

	out.WriteString("[")
	out.WriteString(l.Elements.String())
	out.WriteString("]")
	return out.String()
}

type IndexExpression struct {
	Token token.Token // The [ token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode() {}
func (ie *IndexExpression) TokenLiteral() string {
	return ie.Token.Literal
}
func (ie *IndexExpression) String() string {
	var out strings.Builder
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")
	return out.String()
}
