package scanner

import (
    "testing"
    "language/token"
)

func TestNextToken(t *testing.T) {
    input := `
    1337+hello-héllo[0..1]
    let a = 13 * 37 / 2 % 7;
    loop (a > b) {
        if return break
        continue ,
        // something €
        <===>==!!=<
        /*
        something
        */

        /*
        /*
        */
        *
        */
        /**/
        /***/
        /*
        // something
        */
    }
    : else in true false null something "jaja"
    a.b();
    13.37
    try catch
    += -= *= /= %=
    `

    tests := []struct {
        expectedType token.TokenType
        expectedLiteral string
    }{
        {token.INT, "1337"},
        {token.ADD, ""},
        {token.IDENTIFIER, "hello"},
        {token.SUB, ""},
        {token.IDENTIFIER, "héllo"},
        {token.LBRACKET, ""},
        {token.INT, "0"},
        {token.RANGE, ""},
        {token.INT, "1"},
        {token.RBRACKET, ""},
        {token.LET, ""},
        {token.IDENTIFIER, "a"},
        {token.ASSIGN, ""},
        {token.INT, "13"},
        {token.MULT, ""},
        {token.INT, "37"},
        {token.DIV, ""},
        {token.INT, "2"},
        {token.MOD, ""},
        {token.INT, "7"},
        {token.SEMICOLON, ""},
        {token.LOOP, ""},
        {token.LPAREN, ""},
        {token.IDENTIFIER, "a"},
        {token.GT, ""},
        {token.IDENTIFIER, "b"},
        {token.RPAREN, ""},
        {token.LBRACE, ""},
        {token.IF, ""},
        {token.RETURN, ""},
        {token.BREAK, ""},
        {token.CONTINUE, ""},
        {token.COMMA, ""},
        {token.LE, ""},
        {token.EQ, ""},
        {token.GE, ""},
        {token.ASSIGN, ""},
        {token.NEG, ""},
        {token.NEQ, ""},
        {token.LT, ""},
        {token.RBRACE, ""},
        {token.COLON, ""},
        {token.ELSE, ""},
        {token.IN, ""},
        {token.TRUE, ""},
        {token.FALSE, ""},
        {token.NULL, ""},
        {token.IDENTIFIER, "something"},
        {token.STRING, "jaja"},
        {token.IDENTIFIER, "a"},
        {token.DOT, ""},
        {token.IDENTIFIER, "b"},
        {token.LPAREN, ""},
        {token.RPAREN, ""},
        {token.SEMICOLON, ""},
        {token.FLOAT, "13.37"},
        {token.TRY, ""},
        {token.CATCH, ""},
        {token.ADDASSIGN, ""},
        {token.SUBASSIGN, ""},
        {token.MULTASSIGN, ""},
        {token.DIVASSIGN, ""},
        {token.MODASSIGN, ""},
    }

    scanner := New(input)

    for i, tt := range tests {
        nextToken := scanner.NextToken()

        if nextToken.Type != tt.expectedType {
            t.Fatalf("tests[%d], line:%d, column:%d - Type was=%q, expected=%q", i, nextToken.Line, nextToken.Column, nextToken.Type, tt.expectedType)
        }

        if nextToken.Literal != tt.expectedLiteral {
            t.Fatalf("tests[%d], line:%d, column:%d - Literal was=%q, expected=%q", i, nextToken.Line, nextToken.Column, nextToken.Literal, tt.expectedLiteral)
        }
    }
    eofToken := scanner.NextToken()
    if eofToken.Type != token.EOF {
        t.Fatalf("there are characters left at the end")
    }
}

func TestPosition(t *testing.T) {
    input := `let a = 1;
1 + 2 * 3;

 4 + 5
/* ja */ 7
`
    tests := []struct {
        expectedLine int
        expectedColumn int
    }{
        {1, 1},
        {1, 5},
        {1, 7},
        {1, 9},
        {1, 10},

        {2, 1},
        {2, 3},
        {2, 5},
        {2, 7},
        {2, 9},
        {2, 10},

        {4, 2},
        {4, 4},
        {4, 6},
        
        {5, 10},
    }

    scanner := New(input)
    
    for i, tt := range tests {
        nextToken := scanner.NextToken()

        if nextToken.Line != tt.expectedLine {
            t.Fatalf("tests[%d] - Line was=%d, expected=%d", i, nextToken.Line, tt.expectedLine)
        }
        if nextToken.Column != tt.expectedColumn {
            t.Fatalf("tests[%d] - Column was=%d, expected=%d", i, nextToken.Column, tt.expectedColumn)
        }
    }
    eofToken := scanner.NextToken()
    eofExpectedLine := 6
    eofExpectedColumn := 1
    if eofToken.Type != token.EOF {
        t.Fatalf("Expected EOF, got=%q", eofToken.Type)
    }
    if eofToken.Line != eofExpectedLine {
        t.Fatalf("EOF - Line was=%d, expected=%d", eofToken.Line, eofExpectedLine)
    }
    if eofToken.Column != eofExpectedColumn {
        t.Fatalf("EOF - Column was=%d, expected=%d", eofToken.Column, eofExpectedColumn)
    }
}
