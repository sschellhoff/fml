package ast

import (
    "bytes"
)

type ExpressionStatement struct {
    Expr Expression
    PosInfo PositionalInfo
}

func (e *ExpressionStatement) statementNode() {}

func (e *ExpressionStatement) String() string {
    var out bytes.Buffer
    out.WriteString(e.Expr.String())
    out.WriteString(";")
    return out.String()
}

func (e *ExpressionStatement) Position() PositionalInfo {
    return e.PosInfo
}


type IfStatement struct {
    Cond Expression
    Then *BlockStatement
    Else *BlockStatement
    PosInfo PositionalInfo
}

func (i *IfStatement) statementNode() {}

func (i *IfStatement) String() string {
    var out bytes.Buffer

    out.WriteString("if ")
    out.WriteString(i.Cond.String())
    out.WriteString(" ")
    out.WriteString(i.Then.String())
    if i.Else != nil {
        out.WriteString(" else ")
        out.WriteString(i.Else.String())
    }

    return out.String()
}

func (i *IfStatement) Position() PositionalInfo {
    return i.PosInfo
}


type TryCatchStatement struct {
    Try *BlockStatement
    Info string
    Catch *BlockStatement
    PosInfo PositionalInfo
}

func (t *TryCatchStatement) statementNode() {}

func (t *TryCatchStatement) String() string {
    var out bytes.Buffer

    out.WriteString("try ")
    out.WriteString(t.Try.String())
    out.WriteString(" catch ")
    out.WriteString(t.Info)
    out.WriteString(" ")
    out.WriteString(t.Catch.String())

    return out.String()
}

func (t *TryCatchStatement) Position() PositionalInfo {
    return t.PosInfo
}


type LetStatement struct {
    Name string
    Initializer Expression
    PosInfo PositionalInfo
}

func (l *LetStatement) statementNode() {}

func (l *LetStatement) String() string {
    var out bytes.Buffer

    out.WriteString("let ")
    out.WriteString(l.Name)
    out.WriteString(" = ")
    out.WriteString(l.Initializer.String())
    out.WriteString(";")

    return out.String()
}

func (l *LetStatement) Position() PositionalInfo {
    return l.PosInfo
}


type ConstStatement struct {
    Name string
    Initializer Expression
    PosInfo PositionalInfo
}

func (c *ConstStatement) statementNode() {}

func (c *ConstStatement) String() string {
    var out bytes.Buffer

    out.WriteString("const ")
    out.WriteString(c.Name)
    out.WriteString(" = ")
    out.WriteString(c.Initializer.String())
    out.WriteString(";")

    return out.String()
}

func (c *ConstStatement) Position() PositionalInfo {
    return c.PosInfo
}


type BlockStatement struct {
    Statements []Statement
    PosInfo PositionalInfo
}

func (b *BlockStatement) statementNode() {}

func (b *BlockStatement) String() string {
    var out bytes.Buffer

    out.WriteString("{ ")
    for _, s := range b.Statements {
        out.WriteString(s.String())
        out.WriteString(" ")
    }
    out.WriteString("}")
    
    return out.String()
}

func (b *BlockStatement) Position() PositionalInfo {
    return b.PosInfo
}


type ReturnStatement struct {
    Result Expression
    PosInfo PositionalInfo
}

func (r *ReturnStatement) statementNode() {}

func (r *ReturnStatement) String() string {
    var out bytes.Buffer

    out.WriteString("return ")
    out.WriteString(r.Result.String())
    out.WriteString(";")

    return out.String()
}

func (r *ReturnStatement) Position() PositionalInfo {
    return r.PosInfo
}


type BreakStatement struct {
    PosInfo PositionalInfo
}

func (b *BreakStatement) statementNode() {}

func (b *BreakStatement) String() string {
    return "break;"
}

func (b *BreakStatement) Position() PositionalInfo {
    return b.PosInfo
}


type ContinueStatement struct {
    PosInfo PositionalInfo
}

func (c *ContinueStatement) statementNode() {}

func (c *ContinueStatement) String() string {
    return "continue;"
}

func (c *ContinueStatement) Position() PositionalInfo {
    return c.PosInfo
}


type WhileStatement struct {
    Head Expression
    Body *BlockStatement
    PosInfo PositionalInfo
}

func (w *WhileStatement) statementNode() {}

func (w *WhileStatement) String() string {
    var out bytes.Buffer

    out.WriteString("loop ")
    out.WriteString(w.Head.String())
    out.WriteString(w.Body.String())

    return out.String()
}

func (w *WhileStatement) Position() PositionalInfo {
    return w.PosInfo
}


type RangeLoopStatement struct {
    Name string
    RangeExpr Expression
    Body *BlockStatement
    PosInfo PositionalInfo
}

func (r *RangeLoopStatement) statementNode() {}

func (r *RangeLoopStatement) String() string {
    var out bytes.Buffer

    out.WriteString("loop ")
    out.WriteString(r.Name)
    out.WriteString(" in ")
    out.WriteString(r.RangeExpr.String())
    out.WriteString(r.Body.String())

    return out.String()
}

func (r *RangeLoopStatement) Position() PositionalInfo {
    return r.PosInfo
}


type KVRangeLoopStatement struct {
    IndexName string
    ElementName string
    RangeExpr Expression
    Body *BlockStatement
    PosInfo PositionalInfo
}

func (k *KVRangeLoopStatement) statementNode() {}

func (k *KVRangeLoopStatement) String() string {
    var out bytes.Buffer

    out.WriteString("loop ")
    out.WriteString(k.IndexName)
    out.WriteString(", ")
    out.WriteString(k.ElementName)
    out.WriteString(" in ")
    out.WriteString(k.RangeExpr.String())
    out.WriteString(k.Body.String())

    return out.String()
}

func (k *KVRangeLoopStatement) Position() PositionalInfo {
    return k.PosInfo
}


type ImportStatement struct {
    Path string
    Name string
    PosInfo PositionalInfo
}

func (i *ImportStatement) statementNode() {}

func(i *ImportStatement) String() string {
    var out bytes.Buffer

    out.WriteString("import ")
    out.WriteString(i.Path)
    out.WriteString(" as ")
    out.WriteString(i.Name)

    return out.String()
}

func (i *ImportStatement) Position() PositionalInfo {
    return i.PosInfo
}

