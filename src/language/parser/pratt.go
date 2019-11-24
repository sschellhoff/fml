package parser

import (
    "language/ast"
    "language/token"
)

func (p *Parser) registerPrefixFunctions() {
    p.prefixParseFunctions[token.INT] = p.parseInt
    p.prefixParseFunctions[token.FLOAT] = p.parseFloat
    p.prefixParseFunctions[token.STRING] = p.parseString
    p.prefixParseFunctions[token.NULL] = p.parseNull
    p.prefixParseFunctions[token.TRUE] = p.parseBool
    p.prefixParseFunctions[token.FALSE] = p.parseBool
    p.prefixParseFunctions[token.IDENTIFIER] = p.parseIdentifier
    p.prefixParseFunctions[token.LBRACKET] = p.parseArray
    p.prefixParseFunctions[token.LBRACE] = p.parseHash
    p.prefixParseFunctions[token.SUB] = p.unary
    p.prefixParseFunctions[token.ADD] = p.unary
    p.prefixParseFunctions[token.NEG] = p.unary
    p.prefixParseFunctions[token.LPAREN] = p.grouping
    p.prefixParseFunctions[token.FUN] = p.funcLit
}

func (p *Parser) registerInfixFunctions() {
    p.infixParseFunctions[token.ADD] = p.infix
    p.infixParseFunctions[token.SUB] = p.infix
    p.infixParseFunctions[token.MULT] = p.infix
    p.infixParseFunctions[token.DIV] = p.infix
    p.infixParseFunctions[token.MOD] = p.infix
    p.infixParseFunctions[token.AND] = p.infix
    p.infixParseFunctions[token.OR] = p.infix
    p.infixParseFunctions[token.EQ] = p.infix
    p.infixParseFunctions[token.NEQ] = p.infix
    p.infixParseFunctions[token.LT] = p.infix
    p.infixParseFunctions[token.GT] = p.infix
    p.infixParseFunctions[token.LE] = p.infix
    p.infixParseFunctions[token.GE] = p.infix
    p.infixParseFunctions[token.ASSIGN] = p.infix
    p.infixParseFunctions[token.ADDASSIGN] = p.infix
    p.infixParseFunctions[token.SUBASSIGN] = p.infix
    p.infixParseFunctions[token.MULTASSIGN] = p.infix
    p.infixParseFunctions[token.DIVASSIGN] = p.infix
    p.infixParseFunctions[token.MODASSIGN] = p.infix
    p.infixParseFunctions[token.DOT] = p.property
    p.infixParseFunctions[token.LPAREN] = p.call
    p.infixParseFunctions[token.LBRACKET] = p.index
    p.infixParseFunctions[token.RANGE] = p.infix
    p.infixParseFunctions[token.NULLCOAL] = p.infix
    p.infixParseFunctions[token.QUESTION] = p.conditional
}

func (p *Parser) getPrecedence(op token.Token) int {
    prec, ok := precedences[op.Type]
    if !ok {
        return LOWEST
    }
    return prec
}

type (
    prefixParseFunction func() ast.Expression
    infixParseFunction func(ast.Expression) ast.Expression
)

const (
    _ int = iota
    LOWEST
    ASSIGN
    TERNARY
    NULLCOALESCING
    DISJUNCTION
    CONJUNCTION
    EQUALS
    COMPARE
    SUM
    PRODUCT
    RANGE
    PREFIX
    POSTFIX
)

var precedences = map[token.TokenType]int {
    token.RANGE: RANGE,
    token.NULLCOAL: NULLCOALESCING,
    token.ASSIGN: ASSIGN,
    token.ADDASSIGN: ASSIGN,
    token.SUBASSIGN: ASSIGN,
    token.MULTASSIGN: ASSIGN,
    token.DIVASSIGN: ASSIGN,
    token.MODASSIGN: ASSIGN,
    token.QUESTION: TERNARY,
    token.AND: CONJUNCTION,
    token.OR: DISJUNCTION,
    token.EQ: EQUALS,
    token.NEQ: EQUALS,
    token.LT: COMPARE,
    token.GT: COMPARE,
    token.LE: COMPARE,
    token.GE: COMPARE,
    token.ADD: SUM,
    token.SUB: SUM,
    token.MULT: PRODUCT,
    token.DIV: PRODUCT,
    token.MOD: PRODUCT,
    token.DOT: POSTFIX,
    token.LPAREN: POSTFIX,
    token.LBRACKET: POSTFIX,
}
