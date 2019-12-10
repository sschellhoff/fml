package parser

import (
    "testing"
    "language/ast"
    "language/scanner"
)

func TestProgram(t *testing.T) {
    tests := []struct {
        input string
        numStatements int
    }{
        {"hello;", 1},
        {"hello; there;", 2},
        {"1; 2; 3;", 3},
    }

    for _, tt := range tests {
        program := parseProgram(t, tt.input)

        handleProgramLength(t, program, tt.numStatements)
    }
}

func TestMissingSemicolon(t *testing.T) {
    input := `
    import "some module path" as my_module
    let a = 1 + 1337 * 42
    const f = fun(a, b) {
        return a + b
    }
    f(13, 37)
    if a == f(1, 2) {
        print(true)
    }
    `

    program := parseProgram(t, input)

    handleProgramLength(t, program, 5)

    _, ok := program.Statements[0].(*ast.ImportStatement)
    if !ok {
        t.Fatalf("Expected import, got %T", program.Statements[0])
    }

    _, ok = program.Statements[1].(*ast.LetStatement)
    if !ok {
        t.Fatalf("Expected let, got %T", program.Statements[1])
    }

    _, ok = program.Statements[2].(*ast.ConstStatement)
    if !ok {
        t.Fatalf("Expected const, got %T", program.Statements[2])
    }
    
    _, ok = program.Statements[3].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("Expected call expression statement, got %T", program.Statements[3])
    }

    _, ok = program.Statements[4].(*ast.IfStatement)
    if !ok {
        t.Fatalf("Expected if-statement, got %T", program.Statements[4])
    }
}

func TestImport(t *testing.T) {
    input := "import \"some module path\" as my_module;"

    program := parseProgram(t, input)

    handleProgramLength(t, program, 1)

    importStmt, ok := program.Statements[0].(*ast.ImportStatement)
    if !ok {
        t.Fatalf("Expected import, got %T", program.Statements[0])
    }

    if importStmt.Name != "my_module" {
        t.Fatalf("Expected another module alias, got \"%s\", expected \"%s\"", importStmt.Name, "my_module")
    }

    if importStmt.Path != "some module path" {
        t.Fatalf("Expected another modulepath, got \"%s\", expected \"%s\"", importStmt.Path, "some module path")
    }
}

func TestTryCatch(t *testing.T) {
    input := "try {} catch exception {}"

    program := parseProgram(t, input)

    handleProgramLength(t, program, 1)

    _, ok := program.Statements[0].(*ast.TryCatchStatement)

    if !ok {
        t.Fatalf("Expected try-catch, got %T", program.Statements[0])
    }
}

func TestBreakAndContinueStatements(t *testing.T) {
    tests := []struct {
        input string
        expectError bool
        expected string
    }{
        {"break;", true, "line: 1, column: 1, Literal: \"\" [BREAK]: Break is only allowed inside a loop"},
        {"continue;", true, "line: 1, column: 1, Literal: \"\" [CONTINUE]: Continue is only allowed inside a loop"},
        {"loop forever { break; }", false, "break;"},
        {"loop forever { continue; }", false, "continue;"},
        {"loop forever { fun(a, b) { break; }; }", true, "line: 1, column: 28, Literal: \"\" [BREAK]: Break is only allowed inside a loop"},
        {"loop forever { fun(a, b) { continue; }; }", true, "line: 1, column: 28, Literal: \"\" [CONTINUE]: Continue is only allowed inside a loop"},
        {"loop forever { if true { break; } }", false, "if true { break; } else { }"},
        {"loop forever { if false { break; } }", false, "if false { break; } else { }"},
    }

    for _, tt := range tests {
        s := scanner.New(tt.input)
        p := New(s)
        program, err := p.Parse()
        hadError := len(err) != 0
        if hadError != tt.expectError {
            if tt.expectError {
                t.Fatalf("expected an error \"%s\" but got none", tt.expected)
            } else {
                handleParserErrors(t, err)
            }
        }
        if hadError {
            if len(err) != 1 {
                t.Fatalf("Expected the parser to result in 1 error but got %d", len(err))
            }
            errorMsg := err[0].Error()
            if errorMsg != tt.expected {
                t.Fatalf("Expected error msg to be \"%s\" but got \"%s\"", tt.expected, errorMsg)
            }
            continue
        }
        handleProgramLength(t, program, 1)
        loopStmt, ok := program.Statements[0].(*ast.WhileStatement)
        if !ok {
            t.Fatalf("expected WhileStatement but got %T", program.Statements[0])
        }
        if len(loopStmt.Body.Statements) != 1 {
            t.Fatalf("expected loop body to have 1 statement, but got %d", len(loopStmt.Body.Statements))
        }
        result := loopStmt.Body.Statements[0].String()
        if result != tt.expected {
            t.Fatalf("expected statement to be %s but was %s", tt.expected, result)
        }
    }
}

func TestLetStatement(t *testing.T) {
    tests := []struct {
        input string
        name string
        init string
    }{
        {"let a = 1337;", "a", "1337"},
        {"let b = 1337 * 3;", "b", "(1337*3)"},
        {"let hello = 1337;", "hello", "1337"},
    }

    for _, tt := range tests {
        program := parseProgram(t, tt.input)

        handleProgramLength(t, program, 1)

        letStmt, ok := program.Statements[0].(*ast.LetStatement)
        if !ok {
            t.Fatalf("expected LetStatement but got %T", program.Statements[0])
        }
        if letStmt.Name != tt.name {
            t.Fatalf("expected name to be \"%s\", but got \"%s\"", tt.name, letStmt.Name)
        }

        initializer := letStmt.Initializer.String()
        if initializer != tt.init {
            t.Fatalf("expected initializer to be %s, but got %s", tt.init, initializer)
        }
    }
}

func TestConstStatement(t *testing.T) {
    tests := []struct {
        input string
        name string
        init string
    }{
        {"const a = 1337;", "a", "1337"},
        {"const b = 1337 * 3;", "b", "(1337*3)"},
        {"const hello = 1337;", "hello", "1337"},
    }

    for _, tt := range tests {
        program := parseProgram(t, tt.input)

        handleProgramLength(t, program, 1)

        constStmt, ok := program.Statements[0].(*ast.ConstStatement)
        if !ok {
            t.Fatalf("expected ConstStatement but got %T", program.Statements[0])
        }
        if constStmt.Name != tt.name {
            t.Fatalf("expected name to be \"%s\", but got \"%s\"", tt.name, constStmt.Name)
        }

        initializer := constStmt.Initializer.String()
        if initializer != tt.init {
            t.Fatalf("expected initializer to be %s, but got %s", tt.init, initializer)
        }
    }
}

func TestLoopStatement(t *testing.T) {
    tests := []struct {
        input string
        head string
        body string
    }{
        {"loop true { doIt(); }", "true", "{ doIt(); }"},
        {"loop forever { doIt(); }", "true", "{ doIt(); }"},
        {"loop element in 1..11 { println(element); }", "element in (1..11)", "{ println(element); }"},
        {"loop i, e in 1..11 { println(i, \": \", e); }", "i, e in (1..11)", "{ println(i, \": \", e); }"},
    }

    for _, tt := range tests {
        program := parseProgram(t, tt.input)

        handleProgramLength(t, program, 1)

        switch loopStmt := program.Statements[0].(type) {
        case *ast.WhileStatement:
            head := loopStmt.Head.String()
            if head != tt.head {
                t.Fatalf("Expected head to be \"%s\", but got \"%s\"", tt.head, head)
            }

            body := loopStmt.Body.String()
            if body != tt.body {
                t.Fatalf("Expected body to be \"%s\", but got \"%s\"", tt.body, body)
            }
        case *ast.RangeLoopStatement:
            name := loopStmt.Name
            rangeExpr := loopStmt.RangeExpr.String()
            head := name + " in " + rangeExpr
            if head != tt.head {
                t.Fatalf("Expected head to be \"%s\", but got \"%s\"", tt.head, head)
            }

            body := loopStmt.Body.String()
            if body != tt.body {
                t.Fatalf("Expected body to be \"%s\", but got \"%s\"", tt.body, body)
            }
        case *ast.KVRangeLoopStatement:
            indexName := loopStmt.IndexName
            elementName := loopStmt.ElementName
            rangeExpr := loopStmt.RangeExpr.String()
            head := indexName + ", " + elementName + " in " + rangeExpr
            if head != tt.head {
                t.Fatalf("Expected head to be \"%s\", but got \"%s\"", tt.head, head)
            }

            body := loopStmt.Body.String()
            if body != tt.body {
                t.Fatalf("Expected body to be \"%s\", but got \"%s\"", tt.body, body)
            }
        default:
            t.Fatalf("Expected Loop but got %T", program.Statements[0])
        }

    }
}

func TestIfStatement(t *testing.T) {
    tests := []struct {
        input string
        cond string
        thenBlock string
        elseBlock interface{}
    }{
        {"if a < b { 1; }", "(a<b)", "{ 1; }", "{ }"},
        {"if a < b { 1; } else { 2; }", "(a<b)", "{ 1; }", "{ 2; }"},
        {"if a < b { 1; } else if a > b { 2; } else { 3; }", "(a<b)", "{ 1; }", "{ if (a>b) { 2; } else { 3; } }"},
    }

    for _, tt := range tests {
        program := parseProgram(t, tt.input)

        handleProgramLength(t, program, 1)

        ifStmt, ok := program.Statements[0].(*ast.IfStatement)
        if !ok {
            t.Fatalf("expected IfStatement but got %T", program.Statements[0])
        }

        cond := ifStmt.Cond.String()
        if cond != tt.cond {
            t.Fatalf("expected confition to be %s but got %s", tt.cond, cond)
        }
        if ifStmt.Then == nil {
            t.Fatalf("then block was nil")
        }

        thenBlock := ifStmt.Then.String()
        if thenBlock != tt.thenBlock {
            t.Fatalf("expected thenBlock to be %s but got %s", tt.thenBlock, thenBlock)
        }

        if ifStmt.Else == nil {
            if tt.elseBlock != nil {
                t.Fatalf("expected elseBlock to be %s, but got nil", tt.elseBlock)
            }
        } else {
            expectedElseBlock, ok := tt.elseBlock.(string)
            if !ok {
                t.Fatalf("expected expectedElseBlock to be of type string bot got %T", tt.elseBlock)
            }
            elseBlock := ifStmt.Else.String()
            if elseBlock != expectedElseBlock {
                t.Fatalf("expected elseBlock to be %s but got %s", expectedElseBlock, elseBlock)
            }
        }
    }
}

func TestIntegerLiteralExpression(t *testing.T) {
    input := "1337;"
    var expected int64 = 1337

    program := parseProgram(t, input)

    handleProgramLength(t, program, 1)

    exprStmt := toExprStmt(t, program.Statements[0])

    testLiteral(t, exprStmt.Expr, expected)
}

func TestFloatLiteralExpression(t *testing.T) {
    input := "13.37;"
    var expected float64 = 13.37

    program := parseProgram(t, input)

    handleProgramLength(t, program, 1)

    exprStmt := toExprStmt(t, program.Statements[0])

    testLiteral(t, exprStmt.Expr, expected)
}

func TestStringLiteralExpression(t *testing.T) {
    input := `"Hello, World!";`
    expected := "Hello, World!"

    program := parseProgram(t, input)

    handleProgramLength(t, program, 1)

    exprStmt := toExprStmt(t, program.Statements[0])

    testLiteral(t, exprStmt.Expr, expected)
}

func TestNullLiteralExpression(t *testing.T) {
    input := "null;"

    program := parseProgram(t, input)

    handleProgramLength(t, program, 1)

    exprStmt := toExprStmt(t, program.Statements[0])


    _, ok := exprStmt.Expr.(*ast.NullLiteralExpression)
    if !ok {
        t.Fatalf("expected NullLiteral but got %T", exprStmt.Expr)
    }
}

func TestBoolLiteralExpression(t *testing.T) {
    tests := []struct {
        input string
        expected bool
    }{
        {"true;", true},
        {"false;", false},
    }

    for _, tt := range tests {
        program := parseProgram(t, tt.input)

        handleProgramLength(t, program, 1)

        exprStmt := toExprStmt(t, program.Statements[0])

        testLiteral(t, exprStmt.Expr, tt.expected)
    }
}

func TestHashLiteral(t *testing.T) {
    input := "{\"some\": 1337, \"thing\": 42, \"other\": 111};"

    program := parseProgram(t, input)

    handleProgramLength(t, program, 1)

    exprStmt := toExprStmt(t, program.Statements[0])
    hashExpr, ok := exprStmt.Expr.(*ast.HashLiteral)
    if !ok {
        t.Fatalf("expected Hash but got %T", exprStmt.Expr)
    }
    if len(hashExpr.Pairs) != 3 {
        t.Fatalf("expected Hash to have 3 pairs but got %d", len(hashExpr.Pairs))
    }
}

func TestHashWithSpecialKeysHashLiteral(t *testing.T) {
    input := "{1337: 1337, true: 42, nil: 111};"

    program := parseProgram(t, input)

    handleProgramLength(t, program, 1)

    exprStmt := toExprStmt(t, program.Statements[0])
    hashExpr, ok := exprStmt.Expr.(*ast.HashLiteral)
    if !ok {
        t.Fatalf("expected Hash but got %T", exprStmt.Expr)
    }
    if len(hashExpr.Pairs) != 3 {
        t.Fatalf("expected Hash to have 3 pairs but got %d", len(hashExpr.Pairs))
    }
}

func TestEmptyHashLiteral(t *testing.T) {
    input := "{};"

    program := parseProgram(t, input)

    handleProgramLength(t, program, 1)

    exprStmt := toExprStmt(t, program.Statements[0])
    hashExpr, ok := exprStmt.Expr.(*ast.HashLiteral)
    if !ok {
        t.Fatalf("expected Hash but got %T", exprStmt.Expr)
    }
    if len(hashExpr.Pairs) != 0 {
        t.Fatalf("expected Hash to have 0 pairs but got %d", len(hashExpr.Pairs))
    }
}

func TestArrayLiteral(t *testing.T) {
    input := "[1, 2, 3 * 4, 5 + 6 / 7];"
    
    program := parseProgram(t, input)

    handleProgramLength(t, program, 1)

    exprStmt := toExprStmt(t, program.Statements[0])
    arrExpr, ok := exprStmt.Expr.(*ast.ArrayLiteral)
    if !ok {
        t.Fatalf("expected Array but got %T", exprStmt.Expr)
    }
    if len(arrExpr.Elements) != 4 {
        t.Fatalf("expected Array to have 4 elements but got %d", len(arrExpr.Elements))
    }
}

func TestIndexExpression(t *testing.T) {
    tests := []struct {
        input string
    }{
        {"[0, 1, 2, 3][2];"},
    }

    for _, tt := range tests {
        program := parseProgram(t, tt.input)

        handleProgramLength(t, program, 1)

        exprStmt := toExprStmt(t, program.Statements[0])

        _, ok := exprStmt.Expr.(*ast.IndexExpression)
        if !ok {
            t.Fatalf("expected Index expression but got %T", exprStmt.Expr)
        }
    }
}

func TestFunctionLiteralExpression(t *testing.T) {
    tests := []struct {
        input string
        params []string
    }{
        {"fun(a, b, c) { a+b*c; };", []string{"a", "b", "c"}},
        {"fun(a) {};", []string{"a"}},
        {"fun() {};", []string{}},
    }

    for _, tt := range tests {
        program := parseProgram(t, tt.input)

        handleProgramLength(t, program, 1)

        exprStmt := toExprStmt(t, program.Statements[0])

        funcLitStmt, ok := exprStmt.Expr.(*ast.FunctionLiteralExpression)
        if !ok {
            t.Fatalf("Expected FunctionLiteralExpression but got %T", exprStmt.Expr)
        }

        if len(funcLitStmt.Parameters) != len(tt.params) {
            t.Fatalf("Expected %d parameters but got %d", len(tt.params), len(funcLitStmt.Parameters))
        }

        for i, p := range funcLitStmt.Parameters {
            if p != tt.params[i] {
                t.Fatalf("Expected parameter number %d to be %s but got %s", i + 1, tt.params[i], p)
            }
        }

        if funcLitStmt.Body == nil {
            t.Fatalf("FunctionBody should not be nil")
        }
    }
}

func TestCallExpression(t *testing.T) {
    tests := []struct {
        input string
        function string
        args []string
    }{
        {"a();", "a", []string{}},
        {"a(1, 2, 3);", "a", []string{"1", "2", "3"}},
        {"a(1, 2*3, 3);", "a", []string{"1", "(2*3)", "3"}},
    }

    for _, tt := range tests {
        program := parseProgram(t, tt.input)

        handleProgramLength(t, program, 1)

        exprStmt := toExprStmt(t, program.Statements[0])

        callExpr, ok := exprStmt.Expr.(*ast.CallExpression)
        if !ok {
            t.Fatalf("Expected CallExpression but got %T", exprStmt.Expr)
        }

        function := callExpr.Function.String()
        if function != tt.function {
            t.Fatalf("Expected function to be %s, but got %s", tt.function, function)
        }

        if len(tt.args) != len(callExpr.Arguments) {
            t.Fatalf("Expected to have %d arguments but got %d", len(tt.args), len(callExpr.Arguments))
        }

        for i, a := range callExpr.Arguments {
            as := a.String()
            if as != tt.args[i] {
                t.Fatalf("Expected argument number %d to be %s but got %s", i + 1, tt.args[i], as)
            }
        }
    }
}

func TestIdentifierExpression(t *testing.T) {
    tests := []struct {
        input string
        expected string
    }{
        {"hallo;", "hallo"},
        {"Hallo;", "Hallo"},
        {"hallo_1337_ggez;", "hallo_1337_ggez"},
        {"_test;", "_test"},
        {"_;", "_"},
        {"__;", "__"},
        {"ẞ;", "ẞ"},
        {"éè;", "éè"},
    }

    for _, tt := range tests {
        program := parseProgram(t, tt.input)

        handleProgramLength(t, program, 1)

        exprStmt := toExprStmt(t, program.Statements[0])

        testLiteral(t, exprStmt.Expr, tt.expected)
    }
}

func TestUnaryExpression(t *testing.T) {
    tests := []struct {
        input string
        operator string
        value interface{}
    }{
        {"-1;", "-", 1},
        {"-5;", "-", 5},
        {"!1337;", "!", 1337},
    }

    for _, tt := range tests {
        program := parseProgram(t, tt.input)
        handleProgramLength(t, program, 1)

        exprStmt := toExprStmt(t, program.Statements[0])

        expr, ok := exprStmt.Expr.(*ast.UnaryExpression)
        if !ok {
            t.Fatalf("expected UnaryExpression but got ????")
        }

        testLiteral(t, expr.Rhs, tt.value)
    }
}

func TestExpression(t *testing.T) {
    tests := []struct {
        input string
        expected string
    }{
        {"\"hello\";", "\"hello\""},
        {"1;", "1"},
        {"--1;", "(-(-1))"},
        {"!!!\"hello world!\";", "(!(!(!\"hello world!\")))"},
        {"1+(2+3);", "(1+(2+3))"},
        {"(1+2)+3;", "((1+2)+3)"},
        {"-1+2;", "((-1)+2)"},
        {"1+2+3;", "((1+2)+3)"},
        {"1+2*3;", "(1+(2*3))"},
        {"1*2+3;", "((1*2)+3)"},
        {"1+2*3-4/5%6&&1||2+1*(4+5);", "((((1+(2*3))-((4/5)%6))&&1)||(2+(1*(4+5))))"},
        {"a = 1337;", "(a=1337)"},
        {"a.b??1337;", "((a[\"b\"])??1337)"},
        {"a=b=1337;", "(a=(b=1337))"},
        {"a += b;", "(a+=b)"},
        {"a -= b;", "(a-=b)"},
        {"a /= b;", "(a/=b)"},
        {"a *= b;", "(a*=b)"},
        {"a %= b;", "(a%=b)"},
        {"something = a == b ? 1 : \"hello world\";", "(something=((a==b)?1:\"hello world\"))"},
        {"add = fun(a, b) { return a + b; };", "(add=fun(a, b){ return (a+b); })"},
    }

    for _, tt := range tests {
        program := parseProgram(t, tt.input)
        handleProgramLength(t, program, 1)

        exprStmt := toExprStmt(t, program.Statements[0])
        expr := exprStmt.Expr

        stringValue := expr.String()
        if stringValue != tt.expected {
            t.Fatalf("expected \"%s\" but got \"%s\"", tt.expected, stringValue)
        }
    }
}

func testLiteral(t *testing.T, expr ast.Expression, expectedValue interface{}) {
    t.Helper()

    switch value := expectedValue.(type) {
    case int:
        intExpr, ok := expr.(*ast.IntegerLiteralExpression)
        if !ok {
            t.Fatalf("Expected IntegerLiteralExpression but got %T", expr)
        }
        testInteger(t, intExpr, int64(value))
    case int64:
        intExpr, ok := expr.(*ast.IntegerLiteralExpression)
        if !ok {
            t.Fatalf("Expected IntegerLiteralExpression but got %T", expr)
        }
        testInteger(t, intExpr, value)
    case float64:
        floatExpr, ok := expr.(*ast.FloatLiteralExpression)
        if !ok {
            t.Fatalf("Expected FloatLiteralExpression but got %T", expr)
        }
        testFloat(t, floatExpr, value)
    case bool:
        boolExpr, ok := expr.(*ast.BoolLiteralExpression)
        if !ok {
            t.Fatalf("Expected BoolLiteralExpression but got %T", expr)
        }
        testBool(t, boolExpr, value)
    case string:
        strExpr, ok := expr.(*ast.StringLiteralExpression)
        if ok {
            testString(t, strExpr, value)
        } else {
            idfExpr, ok := expr.(*ast.IdentifierExpression)
            if !ok {
                t.Fatalf("Expected IdentifierExpression or StringExpression but got %T", expr)
            }
            testIdentifier(t, idfExpr, value)
        }
    default:
        t.Fatalf("no literal expression: %T", expr)
    }
}

func testInteger(t *testing.T, expr *ast.IntegerLiteralExpression, expectedValue int64) {
    t.Helper()
    
    if expr.Value != expectedValue {
        t.Fatalf("Expected integer to be=%d, bot got=%d", expectedValue, expr.Value)
    }
}

func testFloat(t *testing.T, expr *ast.FloatLiteralExpression, expectedValue float64) {
    t.Helper()

    if expr.Value != expectedValue {
        t.Fatalf("Expected float to be=%f, but got=%f", expectedValue, expr.Value)
    }
}

func testString(t *testing.T, expr *ast.StringLiteralExpression, expectedValue string) {
    t.Helper()
    
    if expr.Value != expectedValue {
        t.Fatalf("Expected string to be=\"%s\", bot got\"%s\"", expectedValue, expr.Value)
    }
}

func testIdentifier(t *testing.T, expr *ast.IdentifierExpression, expectedValue string) {
    t.Helper()
    
    if expr.Name != expectedValue {
        t.Fatalf("Expected identifier to be=\"%s\", bot got\"%s\"", expectedValue, expr.Name)
    }
}

func testBool(t *testing.T, expr *ast.BoolLiteralExpression, expectedValue bool) {
    t.Helper()

    if expr.Value != expectedValue {
        t.Fatalf("Expected bool to be=%T, bzt got=%T", expectedValue, expr.Value)
    }
}

func handleParserErrors(t *testing.T, errors []error) {
    t.Helper()


    if errors != nil && len(errors) > 0 {
        for _, e := range errors {
            t.Errorf(e.Error())
        }
        t.Fatalf("There were %d parser errors", len(errors))
    }
}

func handleProgramLength(t *testing.T, program *ast.Program, expectedLength int) {
    t.Helper()

    if len(program.Statements) != expectedLength {
        t.Fatalf("Expected a program of length 1, but got length=%d", len(program.Statements))
    }
}

func parseProgram(t *testing.T, input string) *ast.Program {
    t.Helper()

    s := scanner.New(input)
    p := New(s)
    program, err := p.Parse()
    
    handleParserErrors(t, err)

    return program
}

func toExprStmt(t *testing.T, stmt ast.Statement) *ast.ExpressionStatement {
    exprStmt, ok := stmt.(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("expected ExpressionStatement but got %T", stmt)
    }
    return exprStmt
}
