package parser

import (
    "language/ast"
    "language/token"
)

func (p *Parser) parseModuleLevelStmt() ast.Statement {
    nextType := p.peek().Type
    switch nextType {
    case token.IMPORT:
        return p.parseImport()
    default:
        return p.parseStmt()
    }
}

func (p *Parser) parseStmt() ast.Statement {
    nextType := p.peek().Type
    switch nextType {
    case token.LET:
        return p.parseLet()
    case token.CONST:
        return p.parseConst()
    case token.IF:
        return p.parseIf()
    case token.LOOP:
        return p.parseLoop()
    case token.RETURN:
        return p.parseReturn()
    case token.BREAK:
        return p.parseBreak()
    case token.CONTINUE:
        return p.parseContinue()
    case token.TRY:
        return p.parseTryCatch()
    }
    return p.parseExprStmt()
}

func (p *Parser) parseLet() *ast.LetStatement {
    if !p.match(token.LET) {
        p.pushNewError("Expected let statement", p.peek())
    }

    name := p.advance()
    if name.Type != token.IDENTIFIER {
        p.pushNewError("Expected an identifier", name)
        return nil
    }

    if !p.match(token.ASSIGN) {
        p.pushNewError("Expected =", p.peek())
        return nil
    }

    expr := p.expression()
    if expr == nil {
        return nil
    }

    p.match(token.SEMICOLON)

    return &ast.LetStatement{Name: name.Literal, Initializer: expr}
}

func (p *Parser) parseConst() *ast.ConstStatement {
    if !p.match(token.CONST) {
        p.pushNewError("Expected const statement", p.peek())
    }

    name := p.advance()
    if name.Type != token.IDENTIFIER {
        p.pushNewError("Expected an identifier", name)
        return nil
    }

    if !p.match(token.ASSIGN) {
        p.pushNewError("Expected =", p.peek())
        return nil
    }

    expr := p.expression()
    if expr == nil {
        return nil
    }

    p.match(token.SEMICOLON)

    return &ast.ConstStatement{Name: name.Literal, Initializer: expr}
}

func (p *Parser) parseIf() *ast.IfStatement {
    if !p.match(token.IF) {
        p.pushNewError("Expected if", p.peek())
        return nil
    }

    cond := p.expression()
    if cond == nil {
        return nil
    }

    thenBlock := p.block()
    if thenBlock == nil {
        return nil
    }

    if p.match(token.ELSE) {
        if p.peek().Type == token.IF {
            elseIf := p.parseIf()
            if elseIf == nil {
                return nil
            }
            elseStmts := []ast.Statement{elseIf}
            elseBlock := &ast.BlockStatement{Statements: elseStmts}
            return &ast.IfStatement{Cond: cond, Then: thenBlock, Else: elseBlock}
        }
        elseBlock := p.block()
        if elseBlock == nil {
            return nil
        }
        return &ast.IfStatement{Cond: cond, Then: thenBlock, Else: elseBlock}
    }

    return &ast.IfStatement{Cond: cond, Then: thenBlock, Else: &ast.BlockStatement{Statements: []ast.Statement{}}}
}

func (p *Parser) parseLoop() ast.Statement {
    p.enterLoop()
    defer p.exitLoop()
    if !p.match(token.LOOP) {
        p.pushNewError("expected loop", p.peek())
        return nil
    }

    if p.is(token.FOREVER) {
        p.advance()
        head := &ast.BoolLiteralExpression{Value: true}
        block := p.block()

        if block == nil {
            return nil
        }

        return &ast.WhileStatement{Head: head, Body: block}
    } else if p.are(token.IDENTIFIER, token.IN) {
        name := p.advance()
        p.advance()
        rangeExpr := p.expression()
        if rangeExpr == nil {
            return nil
        }
        block := p.block()

        if block == nil {
            return nil
        }

        return &ast.RangeLoopStatement{Name: name.Literal, RangeExpr: rangeExpr, Body: block}
    } else if p.are(token.IDENTIFIER, token.COMMA) {
        idxName := p.advance()
        p.advance()
        if !p.is(token.IDENTIFIER) {
            p.pushNewError("expected identifier", p.peek())
            return nil
        }
        elementName := p.advance()
        if !p.match(token.IN) {
            p.pushNewError("expected in", p.peek())
            return nil
        }
        rangeExpr := p.expression()
        if rangeExpr == nil {
            return nil
        }
        block := p.block()

        if block == nil {
            return nil
        }

        return &ast.KVRangeLoopStatement{IndexName: idxName.Literal, ElementName: elementName.Literal, RangeExpr: rangeExpr, Body: block}
    } else {
        head := p.expression()

        if head == nil {
            return nil
        }
        block := p.block()

        if block == nil {
            return nil
        }

        return &ast.WhileStatement{Head: head, Body: block}
    }
}

func (p *Parser) parseReturn() *ast.ReturnStatement {
    if !p.match(token.RETURN) {
        p.pushNewError("expected return statement", p.peek())
        return nil
    }

    if !p.isInFunctionDefinition() {
        p.pushNewError("return is only allowed in function definitions", p.peek())
        return nil
    }

    if p.match(token.SEMICOLON) {
        return &ast.ReturnStatement{Result: &ast.NullLiteralExpression{}}
    }
    result := p.expression()

    p.match(token.SEMICOLON)

    return &ast.ReturnStatement{Result: result}
}

func (p *Parser) parseBreak() *ast.BreakStatement {
    if !p.isInLoop() {
        p.pushNewError("Break is only allowed inside a loop", p.peek())
        return nil
    }

    if !p.match(token.BREAK) {
        p.pushNewError("Expected break statement", p.peek())
        return nil
    }

    p.match(token.SEMICOLON)

    return &ast.BreakStatement{}
}

func (p *Parser) parseContinue() *ast.ContinueStatement {
    if !p.isInLoop() {
        p.pushNewError("Continue is only allowed inside a loop", p.peek())
        return nil
    }

    if !p.match(token.CONTINUE) {
        p.pushNewError("Expected continue statement", p.peek())
        return nil
    }

    p.match(token.SEMICOLON)

    return &ast.ContinueStatement{}
}

func (p *Parser) parseTryCatch() *ast.TryCatchStatement {
    if !p.match(token.TRY) {
        p.pushNewError("Expected try", p.peek())
        return nil
    }
    tryBlock := p.block()
    if tryBlock == nil {
        return nil
    }
    if !p.match(token.CATCH) {
        p.pushNewError("Expected catch", p.peek())
        return nil
    }
    identifier := p.advance()
    if identifier.Type != token.IDENTIFIER {
        p.pushNewError("Expected identifier", identifier)
    }
    catchBlock := p.block()
    if catchBlock == nil {
        return nil
    }
    return &ast.TryCatchStatement{Try: tryBlock, Catch: catchBlock, Info: identifier.Literal}
}

func (p *Parser) block() *ast.BlockStatement {
    if !p.match(token.LBRACE) {
        p.pushNewError("Expected BlockStatement", p.peek())
        return nil
    }

    stmts := make([]ast.Statement, 0)

    for p.peek().Type != token.RBRACE {
        if p.peek().Type == token.EOF {
            p.pushNewError("Unexpected end of file", p.peek())
            return nil
        }

        stmt := p.parseStmt()
        if stmt == nil {
            return nil
        }
        if p.HadErrors() {
            return nil
        }
        stmts = append(stmts, stmt)
    }

    if !p.match(token.RBRACE) {
        p.pushNewError("Expected }", p.peek())
        return nil
    }
    return &ast.BlockStatement{Statements: stmts}
}

func (p *Parser) parseExprStmt() *ast.ExpressionStatement {
    expr := p.expression()

    if expr == nil {
        return nil
    }
    p.match(token.SEMICOLON)
    return &ast.ExpressionStatement{Expr: expr}
}

func (p *Parser) parseImport() *ast.ImportStatement {
    if !p.match(token.IMPORT) {
        p.pushNewError("expected import", p.peek())
        return nil
    }

    path := p.advance()
    if path.Type != token.STRING {
        p.pushNewError("expected string", path)
        return nil
    }

    if !p.match(token.AS) {
        p.pushNewError("expected as", p.peek())
        return nil
    }

    name := p.advance()
    if name.Type != token.IDENTIFIER {
        p.pushNewError("expected identifier", name)
        return nil
    }

    p.match(token.SEMICOLON)

    return &ast.ImportStatement{Path: path.Literal, Name: name.Literal}
}
