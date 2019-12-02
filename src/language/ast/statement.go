package ast

import (
    "bytes"
)

type ExpressionStatement struct {
    Expr Expression
}

func (e *ExpressionStatement) statementNode() {}

func (e *ExpressionStatement) String() string {
    var out bytes.Buffer
    out.WriteString(e.Expr.String())
    out.WriteString(";")
    return out.String()
}


type IfStatement struct {
    Cond Expression
    Then *BlockStatement
    Else *BlockStatement
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


type TryCatchStatement struct {
    Try *BlockStatement
    Info string
    Catch *BlockStatement
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


type LetStatement struct {
    Name string
    Initializer Expression
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


type ConstStatement struct {
    Name string
    Initializer Expression
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


type BlockStatement struct {
    Statements []Statement
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


type ReturnStatement struct {
    Result Expression
}

func (r *ReturnStatement) statementNode() {}

func (r *ReturnStatement) String() string {
    var out bytes.Buffer

    out.WriteString("return ")
    out.WriteString(r.Result.String())
    out.WriteString(";")

    return out.String()
}


type BreakStatement struct {
}

func (b *BreakStatement) statementNode() {}

func (b *BreakStatement) String() string {
    return "break;"
}


type ContinueStatement struct {
}

func (c *ContinueStatement) statementNode() {}

func (c *ContinueStatement) String() string {
    return "continue;"
}


type WhileStatement struct {
    Head Expression
    Body *BlockStatement
}

func (w *WhileStatement) statementNode() {}

func (w *WhileStatement) String() string {
    var out bytes.Buffer

    out.WriteString("loop ")
    out.WriteString(w.Head.String())
    out.WriteString(w.Body.String())

    return out.String()
}


type RangeLoopStatement struct {
    Name string
    RangeExpr Expression
    Body *BlockStatement
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


type KVRangeLoopStatement struct {
    IndexName string
    ElementName string
    RangeExpr Expression
    Body *BlockStatement
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


type ImportStatement struct {
    Path string
    Name string
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
