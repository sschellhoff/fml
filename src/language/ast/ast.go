package ast

import (
    "bytes"
)

type Node interface {
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

func (p *Program) String() string {
    var out bytes.Buffer
    for _, s := range p.Statements {
        out.WriteString(s.String())
    }
    return out.String()
}


type Module struct {
    Statements []Statement
}

func (m *Module) String() {
    var out bytes.Buffer
    for _, s := range m.Statements {
        out.WriteString(s.String())
    }
    return out.String()
}
