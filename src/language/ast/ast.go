package ast

import (
    "bytes"
    "fmt"
)

type PositionalInfo struct {
    Line int
    Column int
    Path string
}

func (p PositionalInfo) String() string {
    return fmt.Sprintf("%s: [line: %d, column: %d]", p.Path, p.Line, p.Column)
}

type Node interface {
    String() string
    Position() PositionalInfo
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
    Path string
    PosInfo PositionalInfo
}

func (p *Program) String() string {
    var out bytes.Buffer
    for _, s := range p.Statements {
        out.WriteString(s.String())
    }
    return out.String()
}

func (p *Program) Position() PositionalInfo {
    return p.PosInfo
}
