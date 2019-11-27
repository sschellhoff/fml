package eval

import (
    "testing"
    "language/scanner"
    "language/parser"
    "language/object"
)

func TestPrograms(t *testing.T) {
    tests := []struct {
        input string
        expected interface{}
    }{
        {`
        const create = fun(name) {
            let this = {};
            this.name = name;
            this.setName = fun(name) {
                this.name = name;
            };
            this.getName = fun() {
                return this.name;
            };

            return this;
        };
        let person = create("Hans Maulwurf");
        println(person.getName());
        person.setName("Mr. Meeseeks");
        person.getName();
        `, "Mr. Meeseeks"},
    }

    for _, tt := range tests {
        evaluated := evaluate(t, tt.input)
        testLiteral(t, evaluated, tt.expected)
    }
}

func TestEvalTryCatch(t *testing.T) {
    tests := []struct {
        input string
        expected interface{}
    }{
        {`
        try {
            i;
        } catch exception {
        }
        1337;
        `, 1337,
        },
        {`
        let a = 0;
        try {
            i;
        } catch exception {
            a = 1337;
        }
        a;
        `, 1337,
        },
        {`
        let a = 0;
        try {
            i;
        } catch exception {
            a = exception;
        }
        a;
        `, "unknown identifier: i",
        },
        {`
        let a = 0;
        try {
            i;
        } catch exception {
            a = unknown;
        }
        a;
        `, &object.Error{Message: "unknown identifier: unknown"},
        },
    }

    for _, tt := range tests {
        evaluated := evaluate(t, tt.input)
        testLiteral(t, evaluated, tt.expected)
    }
}

func TestEvalBreakContinue(t *testing.T) {
    tests := []struct {
        input string
        expected interface{}
    }{
        {`let a = 0;
        loop forever {
            a = a + 1;
            if a > 4 {
                break;
            }
        }
        a;
        `,5},
        {`let a = 0;
        loop i in 0..10 {
            if i > 5 {
                continue;
            }
            a = a + i;
        }
        a;
        `, 15},
    }

    for _, tt := range tests {
        evaluated := evaluate(t, tt.input)
        testLiteral(t, evaluated, tt.expected)
    }
}

func TestRangeOperator(t *testing.T) {
    tests := []struct {
        input string
        expected interface{}
    }{
        {"(0..10)[0];", 0},
        {"(0..10)[9];", 9},
        {"(0..10)[10];", nil},
        {"(10..0)[0];", 10},
        {"(10..0)[9];", 1},
        {"(10..0)[10];", nil},
        {"(0..0)[0];", nil},
    }

    for _, tt := range tests {
        evaluated := evaluate(t, tt.input)
        testLiteral(t, evaluated, tt.expected)
    }
}

func TestEvalReturn(t *testing.T) {
    tests := []struct {
        input string
        expected interface{}
    }{
        {`
        fun() {
            return 10;
        }();
        `, 10},
        {`
        fun(a, b) {
            let c = a;
            loop forever {
                c = c + 1;
                if c > b {
                    return c;
                }
            }
        }(1300, 1336);
        `, 1337},
        {`
        let function = fun(a, b) {
            let function2 = fun(c) {
                return c * 2;
            };
            let result = function2(a);
            return result * b;
        };
        function(2, 3);
        `, 12},
    }

    for _, tt := range tests {
        evaluated := evaluate(t, tt.input)
        testLiteral(t, evaluated, tt.expected)
    }
}

func TestScoping(t *testing.T) {
    tests := []struct {
        input string
        expected interface{}
    }{
        {"let a = 0; if true { let a = 2; } a;", 0},
        {"fun(a, b) { if a > b { let c = 1337; } return c; }(2, 4);", &object.Error{Message: "unknown identifier: c"}},
        {"fun(a, b) { if a > b { let c = 1337; } return c; }(8, 4);", &object.Error{Message: "unknown identifier: c"}},
        {"let a = 42; if a == 42 { let a = 1337; } a;", 42},
        {"let a = 42; if a != 42 { let a = 1337; } a;", 42},
    }

    for _, tt := range tests {
        evaluated := evaluate(t, tt.input)
        testLiteral(t, evaluated, tt.expected)
    }
}

func TestEvalLoop(t *testing.T) {
    tests := []struct {
        input string
        expected interface{}
    }{
        {"let a = 0; loop a < 5 { a = a + 1; } a;", 5},
        {"let a = 0; let arr = [1, 2, 3]; loop e in arr { a = a + e; } a;", 6},
        {"let a = 0; let arr = [1, 2, 3]; let arr2 = [4, 5, 6]; loop i, e in arr { a = a + e * arr2[i]; } a;", 32},
        {`
        let hash = {42: 1337, 9: 3};
        let result = 0;
        loop k in hash {
            result = result + k * hash[k];
        }
        result;
        `, 56181},
        {`
        let hash = {1: 1, 4: 2, 9: 3};
        let result = 0;
        loop k, v in hash {
            result = result + k / v;
        }
        result;
        `, 6,
        },
    }

    for _, tt := range tests {
        evaluated := evaluate(t, tt.input)
        testLiteral(t, evaluated, tt.expected)
    }
}

func TestEvalHash(t *testing.T) {
    input := `let two = "two";
    {
        "one": 10 - 9,
        two: 1 + 1,
        "thr" + "ee": 6 / 2,
        4 : 4,
        true: 5,
        false: 6
    };`
    expected := map[object.HashKey]int64{
        (&object.String{Value: "one"}).HashKey(): 1,
        (&object.String{Value: "two"}).HashKey(): 2,
        (&object.String{Value: "three"}).HashKey(): 3,
        (&object.Integer{Value: 4}).HashKey(): 4,
        TRUE.HashKey(): 5,
        FALSE.HashKey(): 6,
    }

    evaluated := evaluate(t, input)
    result, ok := evaluated.(*object.Hash)
    if !ok {
        t.Fatalf("expected Hash, got %T (%+v)", evaluated, evaluated)
    }

    if len(result.Pairs) != len(expected) {
        t.Fatalf("Hash has wrong number of pairs, expected %d, got %d", len(expected), len(result.Pairs))
    }

    for expectedKey, expectedValue := range expected {
        pair, ok := result.Pairs[expectedKey]
        if !ok {
            t.Errorf("No pair for given key")
        }

        testIntegerObject(t, pair.Value, expectedValue)
    }
}

func TestEvalIndex(t *testing.T) {
    tests := []struct {
        input string
        expected interface{}
    }{
        {"[0, 1, 2, 3][0];", 0},
        {"[0, 1, 2, 3][1];", 1},
        {"[0, 1, 2, 3][3];", 3},
        {"[][-1];", nil},
        {"[][1];", nil},
        {"{}[0];", nil},
        {"{0: 0, 1: 1, 2: 2}[0];", 0},
        {"{0: 0, 1: 1, 2: 2}[1];", 1},
        {"{0: 0, 1: 1, 2: 2}[2];", 2},
        {"{0: 0, 1: 1, 2: 2}[3];", nil},
        {"{\"some\": \"thing\"}[\"some\"];", "thing"},
        {"{\"some\": \"thing\"}[\"not\"];", nil},
        {"{true: false}[true];", false},
        {"{false: true}[false];", true},
        {"{true: false}[false];", nil},
        {"{true: 1, 1: false}[true];", 1},
        {"{true: 1, 1: false}[1];", false},
    }

    for _, tt := range tests {
        evaluated := evaluate(t, tt.input)
        testLiteral(t, evaluated, tt.expected)
    }
}

func TestEvalArray(t *testing.T) {
    tests := []struct {
        input string
        expected []interface{}
    }{
        {"[1, 2, 3];", []interface{}{1, 2, 3}},
        {"[1];", []interface{}{1}},
        {"[];", []interface{}{}},
    }

    for _, tt := range tests {
        evaluated := evaluate(t, tt.input)
        if isError(evaluated) {
            t.Fatalf("got an error: %s", evaluated.String())
        }
        arrObj, ok := evaluated.(*object.Array)
        if !ok {
            t.Fatalf("expected Array but got %T", evaluated)
        }
        if len(arrObj.Elements) != len(tt.expected) {
            t.Fatalf("expected Array of size %d but got %d", len(tt.expected), len(arrObj.Elements))
        }

        for i, e := range arrObj.Elements {
            testLiteral(t, e, tt.expected[i])
        }
    }
}

func TestBuiltin(t *testing.T) {
    tests := []struct {
        input string
        expected interface{}
    }{
        {"len(\"hello\");", 5},
        {`let map = fun(arr, f) {
            let iter = fun(arr, acc) {
                if len(arr) == 0 {
                    return acc;
                } else {
                    return iter(rest(arr), push(acc, f(first(arr))));
                }
            };
            return iter(arr, []);
        };
        const double = fun(a) {
            return a * 2;
        };
        str(map([1, 2, 3], double));
        `,"[2, 4, 6]"},
        {`
        const reduce = fun(arr, initial, f) {
            const iter = fun(arr, result) {
                if len(arr) == 0 {
                    return result;
                } else {
                    return iter(rest(arr), f(result, first(arr)));
                }
            };
            return iter(arr, initial);
        };
        const sum = fun(arr) {
            return reduce(arr, 0, fun(initial, el) { return initial + el; });
        };
        sum([1, 2, 3, 4]);
        `, 10},
        {`
        const myFunc = fun(param) {
            if param > 0 {
                error("param greater than 0");
            } else {
                return param;
            }
        };
        let a = myFunc(0);
        a;
        let b = myFunc(1);
        `, makeError("param greater than 0")},
    }

    for _, tt := range tests {
        evaluated := evaluate(t, tt.input)
        testLiteral(t, evaluated, tt.expected)
    }
}

func TestEvalInt(t *testing.T) {
    tests := []struct {
        input string
        expected int64
    }{
        {"1337;", 1337},
        {"1;", 1},
        {"0;", 0},
    }

    for _, tt := range tests {
        evaluated := evaluate(t, tt.input)
        testIntegerObject(t, evaluated, tt.expected)
    }
}

func TestEvalFloat(t *testing.T) {
    tests := []struct {
        input string
        expected float64
    }{
        {"13.37;", 13.37},
        {"1.0;", 1.0},
        {"0.0;", 0.0},
        {"3.14;", 3.14},
        {"0.1;", 0.1},
    }
    
    for _, tt := range tests {
        evaluated := evaluate(t, tt.input)
        testFloatObject(t, evaluated, tt.expected)
    }
}

func TestEvalBool(t *testing.T) {
    tests := []struct {
        input string
        expected bool
    }{
        {"true;", true},
        {"false;", false},
    }

    for _, tt := range tests {
        evaluated := evaluate(t, tt.input)
        testBooleanObject(t, evaluated, tt.expected)
    }
}

func TestEvalString(t *testing.T) {
    tests := []struct {
        input string
        expected string
    }{
        {"\"hello world\";", "hello world"},
    }

    for _, tt := range tests {
        evaluated := evaluate(t, tt.input)
        testStringObject(t, evaluated, tt.expected)
    }
}

func TestEvalAssign(t *testing.T) {
    tests := []struct {
        input string
        expected interface{}
    }{
        {"let a = 1; a = 3; a;", 3},
        {"let a = 1; let b = a = 3; b;", 3},
        {"let a = 1; a = 3; a;", 3},
        {"const a = 1337; a = 42;", &object.Error{Message: "cannot assign a"}},
        {"let a = [1, 2, 3]; a[0] = 1337; a[0];", 1337},
        {"let h = {true: 1}; h[true] = 42; h[true];", 42},
        {"let h = 13; h += 4;", 17},
        {"let h = 13; h -= 4;", 9},
        {"let h = 13; h *= 2;", 26},
        {"let h = 6; h /= 2;", 3},
        {"let h = 27; h %= 4;", 3},
    }

    for _, tt := range tests {
        evaluated := evaluate(t, tt.input)
        testLiteral(t, evaluated, tt.expected)
    }
}

func TestEvalNull(t *testing.T) {
    input := "null;"
    evaluated := evaluate(t, input)
    testNullObject(t, evaluated)
}

func TestClosure(t *testing.T) {
    tests := []struct {
        input string
        expected interface{}
    }{
        {`
        let newAdder = fun(x) {
            return fun(y) { return x + y; };
        };
        let addTwo = newAdder(2);
        addTwo(3);
        `, 5,
        },
    }

    for _, tt := range tests {
        evaluated := evaluate(t, tt.input)
        testLiteral(t, evaluated, tt.expected)
    }
}

func TestEvalFunctionCall(t *testing.T) {
    tests := []struct {
        input string
        expected interface{}
    }{
        {"fun() { return 1337; }();", 1337},
        {"fun(a) { return a; }(1337);", 1337},
        {"fun(a, b) { return a + b; }(1, 2);", 3},
        {"fun() { if 1 < 2 { if 1 < 2 { return 1; } else { return 2; } } else { return 3; } return 4; }();", 1},
        {"let add = fun(a, b) { return a + b; }; let mult = fun(a, b) { return a * b; }; let calc = fun() { return add(1, 3) + mult(4, 6); }; calc();", 28},
    }

    for _, tt := range tests {
        evaluated := evaluate(t, tt.input)
        testLiteral(t, evaluated, tt.expected)
    }
}

func TestEvalFunctionDefinition(t *testing.T) {
    tests := []struct {
        input string
        parameters []string
        body string
    }{
        {"fun(a, b) { return a + b; };", []string{"a", "b"}, "{ return (a+b); }"},
        {"fun(a) { return a + 5; };", []string{"a"}, "{ return (a+5); }"},
        {"fun() { };", []string{}, "{ }"},
    }

    for _, tt := range tests {
        evaluated := evaluate(t, tt.input)
        funcObj, ok := evaluated.(*object.Function)

        if !ok {
            t.Fatalf("Expected Function but got %T", evaluated)
        }

        if len(funcObj.Parameters) != len(tt.parameters) {
            t.Fatalf("Expected %d parameters but got %d", len(tt.parameters), len(funcObj.Parameters))
        }

        for i, p := range funcObj.Parameters {
            if p != tt.parameters[i] {
                t.Fatalf("Expected parameters number %d to be %s, but got %s", i, tt.parameters[i], p)
            }
        }

        if funcObj.Body.String() != tt.body {
            t.Fatalf("Expected Body to be %s, but got %s", tt.body, funcObj.Body.String())
        }
    }
}

func TestEvalIf(t *testing.T) {
    tests := []struct {
        input string
        expected interface{}
    }{
        {"if true { 1; } else { 2; }", 1},
        {"if false { 1; } else { 2; }", 2},
        {"if true { 1; }", 1},
        {"if false { 1; }", nil},
    }

    for _, tt := range tests {
        evaluated := evaluate(t, tt.input)
        testLiteral(t, evaluated, tt.expected)
    }
}

func TestEvalUnary(t *testing.T) {
    tests := []struct {
        input string
        expected interface{}
    }{
        {"!false;", true},
        {"!true;", false},
        {"!!true;", true},
        {"!!!true;", false},
        {"!5;", false},
        {"!0;", false},
        {"!null;", true},
        {"-1;", -1},
        {"+1;", 1},
    }

    for _, tt := range tests {
        evaluated := evaluate(t, tt.input)
        testLiteral(t, evaluated, tt.expected)
    }
}

func TestEvalInfix(t *testing.T) {
    tests := []struct {
        input string
        expected interface{}
    }{
        {"1 + 2;", 3},
        {"1 - 2;", -1},
        {"3 * 2;", 6},
        {"6 / 2;", 3},
        {"11 % 2;", 1},
        {"false && false;", false},
        {"false && true;", false},
        {"true && false;", false},
        {"true && true;", true},
        {"false || false;", false},
        {"false || true;", true},
        {"false || true;", true},
        {"true || true;", true},
        {"null || true;", true},
        {"null && true;", false},
        {"false && a;", false}, // tests short circuit because a is undefined
        {"true || a;", true}, // tests short circuit because a is undefined
        {"\"hello \" + \"world!\";", "hello world!"},
        {"null ?? 1337;", 1337},
        {"42 ?? 1337;", 42},
        {"1 < 2;", true},
        {"2 < 1;", false},
        {"1 > 2;", false},
        {"2 > 1;", true},
        {"1 <= 1;", true},
        {"1 <= 2;", true},
        {"2 <= 1;", false},
        {"1 >= 1;", true},
        {"1 >= 2;", false},
        {"2 >= 1;", true},
        {"1 == 1;", true},
        {"1 == 2;", false},
        {"1 != 1;", false},
        {"1 != 2;", true},
        {"true == true;", true},
        {"false == false;", true},
        {"true == false;", false},
        {"false != false;", false},
        {"true != true;", false},
        {"false != true;", true},
        {"\"hello\" == \"hello\";", true},
        {"\"hello\" == \"world\";", false},
        {"\"hello\" != \"hello\";", false},
        {"\"hello\" != \"world\";", true},
    }

    for _, tt := range tests {
        evaluated := evaluate(t, tt.input)
        testLiteral(t, evaluated, tt.expected)
    }
}

func TestConditional(t *testing.T) {
    tests := []struct {
        input string
        expected interface{}
    }{
        {"1 == 1 ? 1337 : 42;", 1337},
        {"1 != 1 ? 1337 : 42;", 42},
    }

    for _, tt := range tests {
        evaluated := evaluate(t, tt.input)
        testLiteral(t, evaluated, tt.expected)
    }
}

func TestErrorHandling(t *testing.T) {
    tests := []struct {
        input string
        expectedMsg string
    }{
        {"1 + \"hello\";", "operands on infix expressions need to be of the same type"},
        {"-\"hello\";", "unsupported unary right hand side type"},
        {"true < true;", "unsupported infix expression"},
        {"\"hello\" - \"hello\";", "unsupported infix operator on strings"},
        {"nothing;", "unknown identifier: nothing"},
        {"{}[fun(){}];", "unusable as hashkey: FUNCTION"},
        {"fun(a, b) { a = 1337; }(1, 2);", "cannot assign a"},
    }

    for _, tt := range tests {
        evaluated := evaluate(t, tt.input)
        errorObj, ok := evaluated.(*object.Error)

        if !ok {
            t.Fatalf("expected error but got %T", evaluated)
        }
        errorMsg := errorObj.Message

        if errorMsg != tt.expectedMsg {
            t.Fatalf("expected error message to be \"%s\" but got \"%s\"", tt.expectedMsg, errorMsg)
        }
    }
}

func TestLetStatement(t *testing.T) {
    tests := []struct {
        input string
        expected interface{}
    }{
        {"let a = 1337; a;", 1337},
        {"let a = 5; let b = a + a; let c = a + b; c;", 15},
        {"let a = 5; let b = a + a; let c = a + b; b;", 10},
    }

    for _, tt := range tests {
        evaluated := evaluate(t, tt.input)
        testLiteral(t, evaluated, tt.expected)
    }
}

func TestConstStatement(t *testing.T) {
    tests := []struct {
        input string
        expected interface{}
    }{
        {"const a = 1337; a;", 1337},
        {"const a = 5; const b = a + a; const c = a + b; c;", 15},
        {"const a = 5; const b = a + a; let c = a + b; b;", 10},
        {"const a = 5; let b = 3; let c = a + b; const d = c + 3; d;", 11},
    }

    for _, tt := range tests {
        evaluated := evaluate(t, tt.input)
        testLiteral(t, evaluated, tt.expected)
    }
}

func testLiteral(t *testing.T, obj object.Object, expected interface{}) {
    t.Helper()

    expectedError, ok := expected.(*object.Error)
    if ok {
        testErrorObject(t, obj, expectedError)
        return
    }

    errorObj, ok := obj.(*object.Error)
    if ok {
        t.Fatalf("got error: %s", errorObj.Message)
    }

    switch expected := expected.(type) {
    case int:
        testIntegerObject(t, obj, int64(expected))
    case int64:
        testIntegerObject(t, obj, expected)
    case float64:
        testFloatObject(t, obj, expected)
    case bool:
        testBooleanObject(t, obj, expected)
    case string:
        testStringObject(t, obj, expected)
    default:
        if expected == nil {
            testNullObject(t, obj)
            return
        }
        t.Fatalf("unknown expected literal of type %T", expected)
    }
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) {
    t.Helper()

    intObj, ok := obj.(*object.Integer)
    if !ok {
        t.Fatalf("expected Integer but got %T", obj)
    }

    if intObj.Value != expected {
        t.Fatalf("expected Integer to be %d, but got %d", expected, intObj.Value)
    }
}

func testFloatObject(t *testing.T, obj object.Object, expected float64) {
    t.Helper()

    floatObj, ok := obj.(*object.Float)
    if !ok {
        t.Fatalf("expected Float but got %T", obj)
    }

    if floatObj.Value != expected {
        t.Fatalf("expected Float to be %f, but got %f", expected, floatObj.Value)
    }
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) {
    t.Helper()

    boolObj, ok := obj.(*object.Boolean)
    if !ok {
        t.Fatalf("expected Bool but got %T", obj)
    }

    if boolObj.Value != expected {
        t.Fatalf("expected Bool to be %t, but got %t", expected, boolObj.Value)
    }
}

func testStringObject(t *testing.T, obj object.Object, expected string) {
    t.Helper()

    strObj, ok := obj.(*object.String)
    if !ok {
        t.Fatalf("expected String but got %T", obj)
    }

    if strObj.Value != expected {
        t.Fatalf("expected String to be %s, but got %s", expected, strObj.Value)
    }
}

func testNullObject(t *testing.T, obj object.Object) {
    t.Helper()

    _, ok := obj.(*object.Null)
    if !ok {
        t.Fatalf("expected Null bot got %T", obj)
    }
}

func testErrorObject(t *testing.T, obj object.Object, expected *object.Error) {
    t.Helper()

    errObj, ok := obj.(*object.Error)
    if !ok {
        t.Fatalf("expected Error but got %T", obj)
    }

    if errObj.Message != expected.Message {
        t.Fatalf("expected Error to be \"%s\", but got \"%s\"", expected.Message, errObj.Message)
    }
}

func evaluate(t *testing.T, input string) object.Object {
    t.Helper()

    s := scanner.New(input)
    p := parser.New(s)
    program, errors := p.Parse()
    handleParserErrors(t, errors)

    env := object.NewEnvironment()
    return Eval(program, env, make(map[string]*object.Module))
}

func handleParserErrors(t *testing.T, errors []error) {
    t.Helper()

    if len(errors) > 0 {
        for _, err := range errors {
            t.Errorf("%s", err.Error())
        }
        t.Fatalf("There were %d parsing errors", len(errors))
    }
}
