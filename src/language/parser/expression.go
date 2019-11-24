package parser

import (
    "fmt"
    "strconv"
    "language/ast"
    "language/token"
)

func (p *Parser) expression() ast.Expression {
    return p.expressionWithPrecedence(LOWEST)
}

func (p *Parser) expressionWithPrecedence(prec int) ast.Expression {
    tok := p.peek()
    prefix, ok := p.prefixParseFunctions[tok.Type]
    if !ok {
        p.pushNewError("expected atomic or prefix expression", tok) // could not parse token: tok
        return nil
    }
    lhs := prefix()

    for prec < p.getPrecedence(p.peek()) {

        tok = p.peek()

        infix, ok := p.infixParseFunctions[tok.Type]
        if !ok {
            return lhs
        }
        lhs = infix(lhs)
    }

    return lhs
}

func (p *Parser) unary() ast.Expression {
    op := p.advance()
    if !isUnaryOperator(op) {
        p.pushNewError(fmt.Sprintf("unexpected unary operator: %s", string(op.Type)), op)
        return nil
    }
    rhs := p.expressionWithPrecedence(PREFIX)
    if rhs == nil {
        return nil
    }

    return &ast.UnaryExpression{Op: op, Rhs: rhs}
}

func (p *Parser) infix(lhs ast.Expression) ast.Expression {
    op := p.advance()
    if !isInfixOperator(op) {
        p.pushNewError(fmt.Sprintf("unexpected infix operator: %s", string(op.Type)), op)
        return nil
    }
    prec := p.getPrecedence(op)
    if isRightAssoc(prec) {
        prec--
    }
    rhs := p.expressionWithPrecedence(prec)
    if rhs == nil {
        return nil
    }
    
    return &ast.InfixExpression{Op: op, Lhs: lhs, Rhs: rhs}
}

func (p *Parser) funcLit() ast.Expression {
    p.openFunctionDefinition()
    defer p.closeFunctionDefinition()
    if !p.match(token.FUN) {
        p.pushNewError("Expected function literal", p.peek())
        return nil
    }
    if !p.match(token.LPAREN) {
        p.pushNewError("expected (", p.peek())
        return nil
    }
    params, hadError := p.functionParameters()
    if hadError {
        return nil
    }
    if !p.match(token.RPAREN) {
        p.pushNewError("expected )", p.peek())
        return nil
    }
    body := p.block()
    if body == nil {
        return nil
    }

    return &ast.FunctionLiteralExpression{Parameters: params, Body: body}
}

func (p *Parser) functionParameters() ([]string, bool) {
    params := make([]string, 0)
    if p.peek().Type == token.IDENTIFIER {
        param := p.advance()
        params = append(params, param.Literal)
    }
    for p.peek().Type == token.COMMA {
        p.advance()
        if p.peek().Type != token.IDENTIFIER {
            p.pushNewError("", p.peek())
            return params, true
        }
        param := p.advance()
        params = append(params, param.Literal)
    }

    return params, false
}

func (p *Parser) conditional(cond ast.Expression) ast.Expression {
    if !p.match(token.QUESTION) {
        p.pushNewError("expected ?", p.peek())
        return nil
    }
    thenExpr := p.expression()
    if !p.match(token.COLON) {
        p.pushNewError("expected :", p.peek())
        return nil
    }
    elseExpr := p.expression()

    return &ast.ConditionalExpression{Cond: cond, Then: thenExpr, Else: elseExpr}
}

func (p *Parser) call(function ast.Expression) ast.Expression {
    if !p.match(token.LPAREN) {
        p.pushNewError("Expected (", p.peek())
        return nil
    }

    arguments, hadError := p.functionArguments()
    if hadError {
        return nil
    }

    if !p.match(token.RPAREN) {
        p.pushNewError("Expected )", p.peek())
        return nil
    }

    return &ast.CallExpression{Function: function, Arguments: arguments}
}

func (p *Parser) functionArguments() ([]ast.Expression, bool) {
    args := []ast.Expression{}

    if p.peek().Type != token.RPAREN {
        arg := p.expression()
        if arg == nil {
            return args, true
        }
        args = append(args, arg)
    }
    for p.peek().Type == token.COMMA {
        p.advance()
        
        arg := p.expression()
        if arg == nil {
            return args, true
        }
        args = append(args, arg)
    }

    return args, false
}

func (p *Parser) grouping() ast.Expression {
    if !p.match(token.LPAREN) {
        p.pushNewError("expected '('", p.peek())
        return nil
    }

    expr := p.expression()
     if expr == nil {
         return nil
     }

    if !p.match(token.RPAREN) {
        p.pushNewError("expected ')'", p.peek())
        return nil
    }

    return expr
}

func (p *Parser) parseInt() ast.Expression {
    if p.is(token.INT) {
        tok := p.advance()
        intValue, err := strconv.ParseInt(tok.Literal, 10, 64)
        if err != nil {
            p.pushError(err)
            return nil
        }
        return &ast.IntegerLiteralExpression{Value: intValue}
    }

    p.pushNewError("expected integer-literal", p.peek())
    return nil
}

func (p *Parser) parseFloat() ast.Expression {
    if p.is(token.FLOAT) {
        tok := p.advance()
        floatValue, err := strconv.ParseFloat(tok.Literal, 64)
        if err != nil {
            p.pushError(err)
            return nil
        }
        return &ast.FloatLiteralExpression{Value: floatValue}
    }

    p.pushNewError("expected float-literal", p.peek())
    return nil
}

func (p *Parser) parseString() ast.Expression {
    if p.is(token.STRING) {
        tok := p.advance()
        stringValue := tok.Literal
        return &ast.StringLiteralExpression{Value: stringValue}
    }

    p.pushNewError("expected string-literal", p.peek())
    return nil
}

func (p *Parser) parseIdentifier() ast.Expression {
    if p.is(token.IDENTIFIER) {
        tok := p.advance()
        name := tok.Literal
        return &ast.IdentifierExpression{Name: name}
    }

    p.pushNewError("expected identifier", p.peek())
    return nil
}

func (p *Parser) parseBool() ast.Expression {
    if p.match(token.TRUE) {
        return &ast.BoolLiteralExpression{Value: true}
    }
    if p.match(token.FALSE) {
        return &ast.BoolLiteralExpression{Value: false}
    }
    p.pushNewError("expected boolean", p.peek())
    return nil
}

func (p *Parser) parseNull() ast.Expression {
    if p.match(token.NULL) {
        return &ast.NullLiteralExpression{}
    }
    p.pushNewError("expected null", p.peek())
    return nil
}

func (p *Parser) parseHash() ast.Expression {
    if !p.match(token.LBRACE) {
        p.pushNewError("expected hash", p.peek())
        return nil
    }
    pairs := make(map[ast.Expression]ast.Expression)

    for p.peek().Type != token.RBRACE {
        key := p.expression()

        if !p.match(token.COLON) {
            p.pushNewError("expected :", p.peek())
            return nil
        }

        value := p.expression()

        pairs[key] = value

        if p.peek().Type != token.RBRACE && !p.match(token.COMMA) {
            p.pushNewError("expected , or }", p.peek())
            return nil
        }
    }

    if !p.match(token.RBRACE) {
        p.pushNewError("expected }", p.peek())
        return nil
    }
    return &ast.HashLiteral{Pairs: pairs}
}

func (p *Parser) parseArray() ast.Expression {
    if !p.match(token.LBRACKET) {
        p.pushNewError("expected array", p.peek())
        return nil
    }

    elements := []ast.Expression{}

    if p.peek().Type != token.RBRACKET {
        elem := p.expression()
        if elem == nil {
            return nil
        }
        elements = append(elements, elem)
        for p.peek().Type == token.COMMA {
            p.advance()
            elem = p.expression()
            if elem == nil {
                return nil
            }
            elements = append(elements, elem)
        }
    }

    if !p.match(token.RBRACKET) {
        p.pushNewError("missing ]", p.peek())
        return nil
    }
    return &ast.ArrayLiteral{Elements: elements}
}

func (p *Parser) property(lhs ast.Expression) ast.Expression {
    if !p.match(token.DOT) {
        p.pushNewError("expected property expression", p.peek())
        return nil
    }

    if p.peek().Type != token.IDENTIFIER {
        p.pushNewError("expected identifier", p.peek())
        return nil
    }
    name := p.advance()
    index := &ast.StringLiteralExpression{Value: name.Literal}

    return &ast.IndexExpression{Left: lhs, Index: index}
}

func (p *Parser) index(lhs ast.Expression) ast.Expression {
    if !p.match(token.LBRACKET) {
        p.pushNewError("expected index expression", p.peek())
        return nil
    }

    index := p.expression()
    if index == nil {
        return nil
    }

    if !p.match(token.RBRACKET) {
        p.pushNewError("expected ]", p.peek())
        return nil
    }
    return &ast.IndexExpression{Left: lhs, Index: index}
}

func isUnaryOperator(tok token.Token) bool {
    switch tok.Type {
    case token.ADD:
        return true
    case token.SUB:
        return true
    case token.NEG:
        return true
    }
    return false
}

func isInfixOperator(tok token.Token) bool {
    switch tok.Type {
    case token.ADD:
        return true
    case token.SUB:
        return true
    case token.MULT:
        return true
    case token.DIV:
        return true
    case token.MOD:
        return true
    case token.AND:
        return true
    case token.OR:
        return true
    case token.EQ:
        return true
    case token.NEQ:
        return true
    case token.LT:
        return true
    case token.GT:
        return true
    case token.LE:
        return true
    case token.GE:
        return true
    case token.ASSIGN:
        return true
    case token.ADDASSIGN:
        return true
    case token.SUBASSIGN:
        return true
    case token.MULTASSIGN:
        return true
    case token.DIVASSIGN:
        return true
    case token.MODASSIGN:
        return true
    case token.NULLCOAL:
        return true
    case token.DOT:
        return true
    case token.RANGE:
        return true
    case token.QUESTION:
        return true
    }
    return false
}

func isRightAssoc(prec int) bool {
    return prec == TERNARY || prec == ASSIGN || prec == NULLCOALESCING || prec == PREFIX // prefix should never be returned since it is unary, not infix
}
