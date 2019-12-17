package ast

import (
    "bytes"
    "fmt"
    "strconv"
    "strings"
    "language/token"
)

type IntegerLiteralExpression struct {
    Value int64
    PosInfo PositionalInfo
}

func (i *IntegerLiteralExpression) expressionNode() {}

func (i *IntegerLiteralExpression) String() string {
    return strconv.FormatInt(i.Value, 10)
}

func (i *IntegerLiteralExpression) Position() PositionalInfo {
    return i.PosInfo
}


type FloatLiteralExpression struct {
    Value float64
    PosInfo PositionalInfo
}

func (f *FloatLiteralExpression) expressionNode() {}

func (f *FloatLiteralExpression) String() string {
    return fmt.Sprintf("%f", f.Value)
}

func (f *FloatLiteralExpression) Position() PositionalInfo {
    return f.PosInfo
}


type StringLiteralExpression struct {
    Value string
    PosInfo PositionalInfo
}

func (s *StringLiteralExpression) expressionNode() {}

func (s *StringLiteralExpression) String() string {
    var out bytes.Buffer

    out.WriteString("\"")
    out.WriteString(s.Value)
    out.WriteString("\"")
    return out.String()
}

func (s *StringLiteralExpression) Position() PositionalInfo {
    return s.PosInfo
}


type BoolLiteralExpression struct {
    Value bool
    PosInfo PositionalInfo
}

func (b *BoolLiteralExpression) expressionNode() {}

func (b *BoolLiteralExpression) String() string {
    if b.Value {
        return "true"
    }
    return "false"
}

func (b *BoolLiteralExpression) Position() PositionalInfo {
    return b.PosInfo
}


type IdentifierExpression struct {
    Name string
    PosInfo PositionalInfo
}

func (i *IdentifierExpression) expressionNode() {}

func (i *IdentifierExpression) String() string {
    return i.Name
}

func (i *IdentifierExpression) Position() PositionalInfo {
    return i.PosInfo
}


type NullLiteralExpression struct {
    PosInfo PositionalInfo
}

func (n *NullLiteralExpression) expressionNode() {}

func (n *NullLiteralExpression) String() string {
    return "null"
}

func (n *NullLiteralExpression) Position() PositionalInfo {
    return n.PosInfo
}


type UnaryExpression struct {
    Op token.Token
    Rhs Expression
    PosInfo PositionalInfo
}

func (u *UnaryExpression) expressionNode() {}

func (u *UnaryExpression) String() string {
    var out bytes.Buffer

    out.WriteString("(")
    out.WriteString(string(u.Op.Type))
    out.WriteString(u.Rhs.String())
    out.WriteString(")")

    return out.String()
}

func (u *UnaryExpression) Position() PositionalInfo {
    return u.PosInfo
}


type InfixExpression struct {
    Op token.Token
    Lhs Expression
    Rhs Expression
    PosInfo PositionalInfo
}

func (i *InfixExpression) expressionNode() {}

func (i *InfixExpression) String() string {
    var out bytes.Buffer

    out.WriteString("(")
    out.WriteString(i.Lhs.String())
    out.WriteString(string(i.Op.Type))
    out.WriteString(i.Rhs.String())
    out.WriteString(")")

    return out.String()
}

func (i *InfixExpression) Position() PositionalInfo {
    return i.PosInfo
}


type ConditionalExpression struct {
    Cond Expression
    Then Expression
    Else Expression
    PosInfo PositionalInfo
}

func (c *ConditionalExpression) expressionNode() {}

func (c *ConditionalExpression) String() string {
    var out bytes.Buffer

    out.WriteString("(")
    out.WriteString(c.Cond.String())
    out.WriteString("?")
    out.WriteString(c.Then.String())
    out.WriteString(":")
    out.WriteString(c.Else.String())
    out.WriteString(")")

    return out.String()
}

func (c *ConditionalExpression) Position() PositionalInfo {
    return c.PosInfo
}


type FunctionLiteralExpression struct {
    Parameters []string
    Body *BlockStatement
    PosInfo PositionalInfo
}

func (f *FunctionLiteralExpression) expressionNode() {}

func (f *FunctionLiteralExpression) String() string {
    var out bytes.Buffer

    out.WriteString("fun(")
    out.WriteString(strings.Join(f.Parameters, ", "))
    out.WriteString(")")
    out.WriteString(f.Body.String())

    return out.String()
}

func (f *FunctionLiteralExpression) Position() PositionalInfo {
    return f.PosInfo
}


type CallExpression struct  {
    Function Expression
    Arguments []Expression
    PosInfo PositionalInfo
}

func (c *CallExpression) expressionNode() {}

func (c *CallExpression) String() string {
    var out bytes.Buffer

    out.WriteString(c.Function.String())
    out.WriteString("(")

    args := make([]string, 0)
    for _, a := range c.Arguments {
        args = append(args, a.String())
    }
    out.WriteString(strings.Join(args, ", "))

    out.WriteString(")")

    return out.String()
}

func (c *CallExpression) Position() PositionalInfo {
    return c.PosInfo
}


type ArrayLiteral struct {
    Elements []Expression
    PosInfo PositionalInfo
}

func (a *ArrayLiteral) expressionNode() {}

func (a *ArrayLiteral) String() string {
    var out bytes.Buffer

    elements := []string{}

    for _, e := range a.Elements {
        elements = append(elements, e.String())
    }
    out.WriteString("[")
    out.WriteString(strings.Join(elements, ", "))
    out.WriteString("]")

    return out.String()
}

func (a *ArrayLiteral) Position() PositionalInfo {
    return a.PosInfo
}


type AssignExpression struct {
    Left Expression
    Op token.Token
    Value Expression
    PosInfo PositionalInfo
}

func (a *AssignExpression) expressionNode() {}

func (a *AssignExpression) String() string {
    var out bytes.Buffer

    out.WriteString("(")
    out.WriteString(a.Left.String())
    out.WriteString("=")
    out.WriteString(a.Value.String())
    out.WriteString(")")

    return out.String()
}

func (a *AssignExpression) Position() PositionalInfo {
    return a.PosInfo
}


type IndexExpression struct {
    Left Expression
    Index Expression
    PosInfo PositionalInfo
}

func (i *IndexExpression) expressionNode() {}

func (i *IndexExpression) String() string {
    var out bytes.Buffer

    out.WriteString("(")
    out.WriteString(i.Left.String())
    out.WriteString("[")
    out.WriteString(i.Index.String())
    out.WriteString("])")

    return out.String()
}

func (i *IndexExpression) Position() PositionalInfo {
    return i.PosInfo
}


type HashLiteral struct {
    Pairs map[Expression]Expression
    PosInfo PositionalInfo
}

func (h *HashLiteral) expressionNode() {}

func (h *HashLiteral) String() string {
    var out bytes.Buffer

    pairs := []string{}

    for k, v := range h.Pairs {
        pairs = append(pairs, k.String() + ": " + v.String())
    }

    out.WriteString("{")
    out.WriteString(strings.Join(pairs, ", "))
    out.WriteString("}")

    return out.String()
}

func (h *HashLiteral) Position() PositionalInfo {
    return h.PosInfo
}

