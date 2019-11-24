package parser

import (
    "fmt"
    "errors"
    "language/scanner"
    "language/ast"
    "language/token"
)

type Parser struct {
    scanner *scanner.Scanner
    tokenBuffer []token.Token
    prefixParseFunctions map[token.TokenType]prefixParseFunction
    infixParseFunctions map[token.TokenType]infixParseFunction
    errors []error
    numberOfEnclosingFunctions int
    isInLoopStack []bool
}

func New(scanner *scanner.Scanner) *Parser {
    bufferSize := 2
    p := &Parser{scanner: scanner, tokenBuffer: make([]token.Token, bufferSize), errors: make([]error, 0), numberOfEnclosingFunctions: 0, isInLoopStack: []bool{false}}

    p.prefixParseFunctions = make(map[token.TokenType]prefixParseFunction)
    p.registerPrefixFunctions()
    p.infixParseFunctions = make(map[token.TokenType]infixParseFunction)
    p.registerInfixFunctions()

    for i := 0; i < bufferSize; i++ {
        p.advance()
    }

    return p
}

func (p *Parser) Parse() (*ast.Program, []error) {
    result := ast.Program{Statements: make([]ast.Statement, 0)}

    for !p.HadErrors(){
        stmt := p.parseStmt()
        if stmt != nil {
            result.Statements = append(result.Statements, stmt)
        } else {
            break
        }
        if p.isAtEnd() {
            break
        }
    }

    if !p.isAtEnd() && len(p.errors) == 0 {
        p.pushNewError("There are unparsed tokens left", p.peek())
    }
    return &result, p.errors
}

func (p *Parser) pushError(err error) {
    p.errors = append(p.errors, err)
}

func (p *Parser) pushNewError(msg string, at token.Token) {
    positionalMsg := fmt.Sprintf("line: %d, column: %d, Literal: \"%s\" [%s]: %s", at.Line, at.Column, at.Literal, string(at.Type), msg)
    p.pushError(errors.New(positionalMsg))
}

func (p *Parser) advance() token.Token {
    idxOfLastBufferElement := len(p.tokenBuffer) - 1
    result := p.tokenBuffer[0]
    if result.Type == token.ERROR {
        p.pushNewError(result.Literal, result)
    }

    for i := 0; i < idxOfLastBufferElement; i++ {
        p.tokenBuffer[i] = p.tokenBuffer[i+1]
    }
    p.tokenBuffer[idxOfLastBufferElement] = p.scanner.NextToken()

    return result
}

func (p *Parser) peek() token.Token {
    return p.tokenBuffer[0]
}

func (p *Parser) peek2() token.Token {
    return p.tokenBuffer[1]
}

func (p *Parser) match(ttype token.TokenType) bool {
    if p.is(ttype) {
        p.advance()
        return true
    }
    return false
}

func (p *Parser) is(ttype token.TokenType) bool {
    return p.peek().Type == ttype
}

func (p *Parser) are(fstType token.TokenType, sndType token.TokenType) bool {
    return p.is(fstType) && p.peek2().Type == sndType
}

func (p *Parser) isAtEnd() bool {
    return p.peek().Type == token.EOF
}

func (p *Parser) HadErrors() bool {
    return len(p.errors) > 0
}

func (p *Parser) openFunctionDefinition() {
    p.numberOfEnclosingFunctions++
    p.pushLoopStack(false)
}

func (p *Parser) closeFunctionDefinition() {
    p.numberOfEnclosingFunctions--
    p.popLoopStack()
}

func (p *Parser) isInFunctionDefinition() bool {
    return p.numberOfEnclosingFunctions > 0
}

func (p *Parser) enterLoop() {
    p.pushLoopStack(true)
}

func (p *Parser) exitLoop() {
    p.popLoopStack()
}

func (p *Parser) isInLoop() bool {
    n := len(p.isInLoopStack) - 1
    return p.isInLoopStack[n]
}

func (p *Parser) pushLoopStack(enterLoop bool) {
    p.isInLoopStack = append(p.isInLoopStack, enterLoop)
}

func (p *Parser) popLoopStack() {
    n := len(p.isInLoopStack) - 1
    p.isInLoopStack = p.isInLoopStack[:n]
}
