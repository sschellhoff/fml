package eval

import (
    "fmt"
    "os"
    "path/filepath"
    "language/ast"
    "language/object"
    "language/token"
    "language/frontend"
)

var (
    TRUE = &object.Boolean{Value: true}
    FALSE = &object.Boolean{Value: false}
    NULL = &object.Null{}
    MODULEPATH = ""
)

func Eval(node ast.Node, env *object.Environment, modules map[string]*object.Module) object.Object {
    switch node := node.(type) {
    case *ast.Program:
        return evalProgram(node, env, modules)

    case *ast.ImportStatement:
        path := node.Path
        if MODULEPATH == "" {
            cwd, err := os.Getwd()
            if err != nil {
                return makeError(err.Error())
            }
            MODULEPATH = cwd
        }
        if !filepath.IsAbs(path) {
            path = filepath.Join(MODULEPATH, path)
        }
        name := node.Name
        // don't load module if already loaded
        if foundModule, ok := modules[path]; ok {
            if !env.AddConst(name, foundModule) {
                return makeError("Cannot define module with this name, it is already taken")
            }
            return NULL
        }
        old_MODULEPATH := MODULEPATH
        MODULEPATH = filepath.Dir(path)
        defer func() {
            MODULEPATH = old_MODULEPATH
        }()
        moduleCode, errs := frontend.Build(path)
        if len(errs) > 0 {
            return makeParserErrors(errs)
        }
        module := &object.Module{Env: object.NewEnvironment(), Path: path}
        if !env.AddConst(name, module) {
            return makeError("Cannot define module with this name, it is already taken")
        }
        moduleEnv := module.Env
        modules[path] = module
        return Eval(moduleCode, moduleEnv, modules)

    case *ast.BlockStatement:
        return evalBlockStatement(node, env, modules)

    case *ast.FunctionLiteralExpression:
        parameters := node.Parameters
        body := node.Body
        return &object.Function{Parameters: parameters, Body: body, Env: env}

    case *ast.CallExpression:
        function := Eval(node.Function, env, modules)
        if isError(function) {
            return function
        }
        args := evalExpressions(node.Arguments, env, modules)
        if len(args) == 1 && isError(args[0]) {
            return args[0]
        }

        return applyFunction(function, args, modules)

    case *ast.TryCatchStatement:
        try := Eval(node.Try, env, modules)
        if isReturn(try) || isBreakOrContinue(try) {
            return try
        }
        if isError(try) {
            catchableError, ok := try.(*object.Error)
            if !ok {
                return try
            } else {
                errorInfo := catchableError.Message
                catchEnv := object.NewEnclosingEnvironment(env)
                catchEnv.Add(node.Info, &object.String{Value: errorInfo})
                catch := Eval(node.Catch, catchEnv, modules)
                if isErrorOrReturn(catch) || isBreakOrContinue(try) {
                    return catch
                }
            }
        }

    case *ast.IfStatement:
        cond := Eval(node.Cond, env, modules)
        if isError(cond) {
            return cond
        }
        if isTruthy(cond){
            return Eval(node.Then, env, modules)
        } else {
            return Eval(node.Else, env, modules)
        }
    case *ast.WhileStatement:
        for {
            head := Eval(node.Head, env, modules)
            if isError(head) {
                return head
            }
            if !isTruthy(head) {
                return NULL
            }
            body := Eval(node.Body, env, modules)
            if isErrorOrReturn(body) {
                return body
            }
            if isBreak(body) {
                return NULL
            }
            if isContinue(body) {
                continue
            }
        }
        return NULL

    case *ast.RangeLoopStatement:
        theRange := Eval(node.RangeExpr, env, modules)
        if isError(theRange) {
            return theRange
        }
        name := node.Name
        loopEnv := object.NewEnclosingEnvironment(env)
        loopEnv.Add(name, NULL)
        switch rangeHolder := theRange.(type) {
        case *object.Array:
            for _, e := range rangeHolder.Elements {
                loopEnv.Set(name, e)
                body := Eval(node.Body, loopEnv, modules)
                if isErrorOrReturn(body) {
                    return body
                }
                if isBreak(body) {
                    return NULL
                }
                if isContinue(body) {
                    continue
                }
            }
            return NULL
        case *object.Hash:
            for _, p := range rangeHolder.Pairs {
                key := p.Key
                loopEnv.Set(name, key)
                body := Eval(node.Body, loopEnv, modules)
                if isErrorOrReturn(body) {
                    return body
                }
                if isBreak(body) {
                    return NULL
                }
                if isContinue(body) {
                    continue
                }
            }
        default:
            return makeError("Can only range over array or hash, got %s", theRange.Type())
        }
        return NULL

    case *ast.KVRangeLoopStatement:
        theRange := Eval(node.RangeExpr, env, modules)
        if isError(theRange) {
            return theRange
        }
        indexName := node.IndexName
        elementName := node.ElementName
        loopEnv := object.NewEnclosingEnvironment(env)
        loopEnv.Add(indexName, NULL)
        loopEnv.Add(elementName, NULL)
        switch rangeHolder := theRange.(type) {
        case *object.Array:
            for i, e := range rangeHolder.Elements {
                loopEnv.Set(indexName, &object.Integer{Value: int64(i)})
                loopEnv.Set(elementName, e)
                body := Eval(node.Body, loopEnv, modules)
                if isErrorOrReturn(body) {
                    return body
                }
                if isBreak(body) {
                    return NULL
                }
                if isContinue(body) {
                    continue
                }
            }
            return NULL
        case *object.Hash:
            for _, p := range rangeHolder.Pairs {
                key := p.Key
                value := p.Value
                loopEnv.Set(indexName, key)
                loopEnv.Set(elementName, value)
                body := Eval(node.Body, loopEnv, modules)
                if isErrorOrReturn(body) {
                    return body
                }
                if isBreak(body) {
                    return NULL
                }
                if isContinue(body) {
                    continue
                }
            }
        default:
            return makeError("Can only range over array or hash, got %s", theRange.Type())
        }
        return NULL

    case *ast.LetStatement:
        value := Eval(node.Initializer, env, modules)
        if isError(value) {
            return value
        }

        if !env.Add(node.Name, value) {
            return makeError("Cannot define variable")
        }

    case *ast.ConstStatement:
        value := Eval(node.Initializer, env, modules)
        if isError(value) {
            return value
        }

        if !env.AddConst(node.Name, value) {
            return makeError("Cannot define constant")
        }

    case *ast.ExpressionStatement:
        return Eval(node.Expr, env, modules)

    case *ast.IntegerLiteralExpression:
        return &object.Integer{Value: node.Value}

    case *ast.FloatLiteralExpression:
        return &object.Float{Value: node.Value}

    case *ast.BoolLiteralExpression:
        return boolToBoolean(node.Value)

    case *ast.StringLiteralExpression:
        return &object.String{Value: node.Value}

    case *ast.IdentifierExpression:
        return evalIdentifier(node.Name, env)

    case *ast.NullLiteralExpression:
        return NULL

    case *ast.ArrayLiteral:
        elements := evalExpressions(node.Elements, env, modules)
        if len(elements) == 1 && isError(elements[0]) {
            return elements[0]
        }
        return &object.Array{Elements: elements}

    case *ast.HashLiteral:
        return evalHashLiteral(node, env, modules)

    case *ast.IndexExpression:
        lhs := Eval(node.Left, env, modules)
        if isError(lhs) {
            return lhs
        }

        idx := Eval(node.Index, env, modules)
        if isError(idx) {
            return idx
        }
        return evalIndex(lhs, idx)

    case *ast.BreakStatement:
        return &object.Break{}

    case *ast.ContinueStatement:
        return &object.Continue{}

    case *ast.ReturnStatement:
        result := Eval(node.Result, env, modules)
        if isError(result) {
            return result
        }
        return &object.Return{Value: result}

    case *ast.UnaryExpression:
        return evalUnary(node, env, modules)

    case *ast.InfixExpression:
        return evalInfix(node, env, modules)

    case *ast.ConditionalExpression:
        return evalConditional(node, env, modules)
    default:
        return makeError("Unknown expression of type: %T", node)
    }
    return NULL
}

func evalIndex(lhs object.Object, index object.Object) object.Object {
    switch lhs := lhs.(type) {
    case *object.Array:
        return evalArray(lhs, index)
    case *object.Hash:
        return evalHash(lhs, index)
    case *object.Module:
        return evalModule(lhs, index)
    case *object.String:
        return evalStringIndex(lhs, index)
    default:
        return makeError("Cannot index on %s", lhs.Type())
    }
}

func evalHashLiteral(node *ast.HashLiteral, env *object.Environment, modules map[string]*object.Module) object.Object {
    pairs := make(map[object.HashKey]object.HashPair)

    for keyNode, valueNode := range node.Pairs {
        key := Eval(keyNode, env, modules)
        if isError(key) {
            return key
        }

        hashKey, ok := key.(object.Hashable)
        if !ok {
            return makeError("key is not hashable: %s", key.Type())
        }

        value := Eval(valueNode, env, modules)
        if isError(value) {
            return value
        }

        hashed := hashKey.HashKey()
        pairs[hashed] = object.HashPair{Key: key, Value: value}
    }
    return &object.Hash{Pairs: pairs}
}

func evalHash(lhs *object.Hash, index object.Object) object.Object {
    key, ok := index.(object.Hashable)
    if !ok {
        return makeError("unusable as hashkey: %s", index.Type())
    }

    pair, ok := lhs.Pairs[key.HashKey()]
    if !ok {
        return NULL
    }
    return pair.Value
}

func evalArray(lhs *object.Array, index object.Object) object.Object {
    idx, ok := index.(*object.Integer)
    if !ok {
        return makeError("Can only use integer as index on array, got %s", index.Type())
    }
    if idx.Value < 0 || idx.Value >= int64(len(lhs.Elements)) {
        return NULL
    }
    return lhs.Elements[idx.Value]
}

func evalStringIndex(lhs *object.String, index object.Object) object.Object {
    idx, ok := index.(*object.Integer)
    if !ok {
        return makeError("Can only use integer as index on string, got %s", index.Type())
    }
    return builtins["substring"].Function(lhs, idx, &object.Integer{Value: idx.Value+1})
}

func evalModule(lhs *object.Module, index object.Object) object.Object {
    name := index.String()
    result, ok := lhs.Env.Get(name)
    if !ok {
        return makeError("Cannot find %s in module", name)
    }
    return result
}

func evalIdentifier(name string, env *object.Environment) object.Object {
    result, ok := env.Get(name)
    if ok {
        return result
    }
    builtin, ok := builtins[name]
    if ok {
        return builtin
    }
    return makeError("unknown identifier: %s", name)
}

func evalProgram(program *ast.Program, env *object.Environment, modules map[string]*object.Module) object.Object {
    var result object.Object = NULL

    for _, stmt := range program.Statements {
        result = Eval(stmt, env, modules)

        if isError(result) {
            return result
        }
    }

    return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment, modules map[string]*object.Module) object.Object {
    var result object.Object = NULL
    blockEnv := object.NewEnclosingEnvironment(env)

    for _, stmt := range block.Statements {
        result = Eval(stmt, blockEnv, modules)
        
        if isErrorOrReturn(result) || isBreakOrContinue(result) {
            return result
        }
    }

    return result
}

func applyFunction(fn object.Object, args []object.Object, modules map[string]*object.Module) object.Object {
    function, ok := fn.(*object.Function)
    if ok {
        extendedEnv := extendFunctionEnv(function, args)
        evaluated := Eval(function.Body, extendedEnv, modules)
        // TODO check if i need to handle error
        return unwrapReturnValue(evaluated)
    }

    builtin, ok := fn.(*object.Builtin)
    if ok {
        return builtin.Function(args...)
    }

    return makeError("cannot call a non function %T", fn)
}

func evalExpressions(exprs []ast.Expression, env *object.Environment, modules map[string]*object.Module) []object.Object {
    var result []object.Object

    for _, e := range exprs {
        evaluated := Eval(e, env, modules)
        if isError(evaluated) {
            return []object.Object{evaluated}
        }
        result = append(result, evaluated)
    }

    return result
}

func evalConditional(expr *ast.ConditionalExpression, env *object.Environment, modules map[string]*object.Module) object.Object {
    cond := Eval(expr.Cond, env, modules)
    if isError(cond) {
        return cond
    }

    if isTruthy(cond) {
        return Eval(expr.Then, env, modules)
    } else {
        return Eval(expr.Else, env, modules)
    }
}

func evalUnary(expr *ast.UnaryExpression, env *object.Environment, modules map[string]*object.Module) object.Object {
    op := expr.Op
    rhs := expr.Rhs
    value := Eval(rhs, env, modules)

    if isError(value) {
        return value
    }

    if op.Type == token.NEG {
        return computeNegExpr(value)
    }

    switch typedValue := value.(type) {
    case *object.Integer:
        if op.Type == token.ADD {
            return typedValue
        } else if op.Type == token.SUB {
            return &object.Integer{Value: -typedValue.Value}
        } else {
            return makeError("unsupported unary expression")
        }
    case *object.Float:
        if op.Type == token.ADD {
            return typedValue
        } else if op.Type == token.SUB {
            return &object.Float{Value: -typedValue.Value}
        } else {
            return makeError("unsupported unary expression")
        }
    default:
        return makeError("unsupported unary right hand side type")
    }

}

func evalAssign(left ast.Expression, right ast.Expression, env *object.Environment, modules map[string]*object.Module) object.Object{
    switch lhs := left.(type) {
    case *ast.IdentifierExpression:
        name := lhs.Name
        rhs := Eval(right, env, modules)
        if isError(rhs) {
            return rhs
        }
        ok := env.Set(name, rhs)
        if !ok {
            return makeError("cannot assign %s", name)
        }
        return rhs
    case *ast.IndexExpression:
        rhs := Eval(right, env, modules)
        if isError(rhs) {
            return rhs
        }
        return evalIndexSet(lhs, rhs, env, modules)
    default:
        return makeError("can only assign to variables")
    }
}

func evalIndexSet(expr *ast.IndexExpression, value object.Object, env *object.Environment, modules map[string]*object.Module) object.Object {
    lhs := Eval(expr.Left, env, modules)
    switch lhs := lhs.(type) {
    case *object.Array:
        index := Eval(expr.Index, env, modules)
        if isError(index) {
            return index
        }
        return evalArrayIndexSet(lhs, index, value)
    case *object.Hash:
        index := Eval(expr.Index, env, modules)
        if isError(index) {
            return index
        }
        return evalHashIndexSet(lhs, index, value)
    case *object.Module:
        index := Eval(expr.Index, env, modules)
        if isError(index) {
            return index
        }
        return evalModuleIndexSet(lhs, index, value)
    default:
        return makeError("cannot use index expression on %s", lhs.Type())
    }
}

func evalArrayIndexSet(arr *object.Array, index object.Object, value object.Object) object.Object {
    idx, ok := index.(*object.Integer)
    if !ok {
        return makeError("can only use integer as array index but got %s", index.Type())
    }
    if idx.Value < 0 || idx.Value >= int64(len(arr.Elements)) {
        return makeError("index out of bounds: %d", idx.Value)
    }
    arr.Elements[idx.Value] = value
    return value
}

func evalHashIndexSet(hash *object.Hash, index object.Object, value object.Object) object.Object {
    key, ok := index.(object.Hashable)
    if !ok {
        return makeError("cannot use %s as hash key", index.Type())
    }
    hash.Pairs[key.HashKey()] = object.HashPair{Key: index, Value: value}
    return value
}

func evalModuleIndexSet(module *object.Module, index object.Object, value object.Object) object.Object {
    name := index.String()
    ok := module.Env.Set(name, value)
    if !ok {
        return makeError("")
    }
    return value
}

func evalInfix(expr *ast.InfixExpression, env *object.Environment, modules map[string]*object.Module) object.Object {
    if expr.Op.Type == token.ASSIGN {
        return evalAssign(expr.Lhs, expr.Rhs, env, modules)
    }

    if expr.Op.Type == token.ADDASSIGN {
        newRhs := &ast.InfixExpression{Lhs: expr.Lhs, Rhs: expr.Rhs, Op: token.FromType(token.ADD, expr.Op.Line, expr.Op.Column)}
        return evalAssign(expr.Lhs, newRhs, env, modules)
    }
    if expr.Op.Type == token.SUBASSIGN {
        newRhs := &ast.InfixExpression{Lhs: expr.Lhs, Rhs: expr.Rhs, Op: token.FromType(token.SUB, expr.Op.Line, expr.Op.Column)}
        return evalAssign(expr.Lhs, newRhs, env, modules)
    }
    if expr.Op.Type == token.MULTASSIGN {
        newRhs := &ast.InfixExpression{Lhs: expr.Lhs, Rhs: expr.Rhs, Op: token.FromType(token.MULT, expr.Op.Line, expr.Op.Column)}
        return evalAssign(expr.Lhs, newRhs, env, modules)
    }
    if expr.Op.Type == token.DIVASSIGN {
        newRhs := &ast.InfixExpression{Lhs: expr.Lhs, Rhs: expr.Rhs, Op: token.FromType(token.DIV, expr.Op.Line, expr.Op.Column)}
        return evalAssign(expr.Lhs, newRhs, env, modules)
    }
    if expr.Op.Type == token.MODASSIGN {
        newRhs := &ast.InfixExpression{Lhs: expr.Lhs, Rhs: expr.Rhs, Op: token.FromType(token.MOD, expr.Op.Line, expr.Op.Column)}
        return evalAssign(expr.Lhs, newRhs, env, modules)
    }

    lhs := Eval(expr.Lhs, env, modules)

    if isError(lhs) {
        return lhs
    }

    switch expr.Op.Type {
    case token.AND:
        if !isTruthy(lhs) {
            return FALSE
        }
        return boolToBoolean(isTruthy(Eval(expr.Rhs, env, modules)))
    case token.OR:
        if isTruthy(lhs) {
            return TRUE
        }
        return boolToBoolean(isTruthy(Eval(expr.Rhs, env, modules)))
    case token.NULLCOAL:
        if lhs == NULL {
            return Eval(expr.Rhs, env, modules)
        } else {
            return lhs
        }
    }

    rhs := Eval(expr.Rhs, env, modules)
    if isError(rhs) {
        return rhs
    }

    if lhs.Type() == rhs.Type() {
        if lhs.Type() == object.INTEGER_OBJECT {
            lhsIo, _ := lhs.(*object.Integer)
            rhsIo, _ := rhs.(*object.Integer)
            return evalIntegerInfix(expr.Op, lhsIo, rhsIo)
        }
        if lhs.Type() == object.FLOAT_OBJECT {
            lhsFo := lhs.(*object.Float)
            rhsFo := rhs.(*object.Float)
            return evalFloatInfix(expr.Op, lhsFo, rhsFo)
        }
        if lhs.Type() == object.STRING_OBJECT {
            lhsSo, _ := lhs.(*object.String)
            rhsSo, _ := rhs.(*object.String)
            return evalStringInfix(expr.Op, lhsSo, rhsSo)
        }
        switch expr.Op.Type {
        case token.EQ:
            return boolToBoolean(lhs == rhs)
        case token.NEQ:
            return boolToBoolean(lhs != rhs)
        }
        return makeError("unsupported infix expression")
    }
    switch expr.Op.Type {
    case token.EQ:
        return boolToBoolean(lhs == rhs)
    case token.NEQ:
        return boolToBoolean(lhs != rhs)
    }

    return makeError("operands on infix expressions need to be of the same type")
}

func evalStringInfix(op token.Token, lhs *object.String, rhs *object.String) object.Object {
    switch op.Type {
    case token.ADD:
        return &object.String{Value: lhs.Value + rhs.Value}
    case token.EQ:
        return boolToBoolean(lhs.Value == rhs.Value)
    case token.NEQ:
        return boolToBoolean(lhs.Value != rhs.Value)
    }
    return makeError("unsupported infix operator on strings")
}

func evalIntegerInfix(op token.Token, lhs *object.Integer, rhs *object.Integer) object.Object {
    switch op.Type {
    case token.ADD:
        return &object.Integer{Value: lhs.Value + rhs.Value}
    case token.SUB:
        return &object.Integer{Value: lhs.Value - rhs.Value}
    case token.MULT:
        return &object.Integer{Value: lhs.Value * rhs.Value}
    case token.DIV:
        return &object.Integer{Value: lhs.Value / rhs.Value}
    case token.MOD:
        return &object.Integer{Value: lhs.Value % rhs.Value}
    case token.LT:
        return boolToBoolean(lhs.Value < rhs.Value)
    case token.GT:
        return boolToBoolean(lhs.Value > rhs.Value)
    case token.LE:
        return boolToBoolean(lhs.Value <= rhs.Value)
    case token.GE:
        return boolToBoolean(lhs.Value >= rhs.Value)
    case token.EQ:
        return boolToBoolean(lhs.Value == rhs.Value)
    case token.NEQ:
        return boolToBoolean(lhs.Value != rhs.Value)
    case token.RANGE:
        elements := []object.Object{}
        if lhs.Value < rhs.Value {
            for i := lhs.Value; i < rhs.Value; i++ {
                elements = append(elements, &object.Integer{Value: i})
            }
        } else if lhs.Value > rhs.Value {
            for i := lhs.Value; i > rhs.Value; i-- {
                elements = append(elements, &object.Integer{Value: i})
            }
        }
        return &object.Array{Elements: elements}
    default:
        return makeError("unsupported infix operator on integers")
    }
}

func evalFloatInfix(op token.Token, lhs *object.Float, rhs *object.Float) object.Object {
    switch op.Type {
    case token.ADD:
        return &object.Float{Value: lhs.Value + rhs.Value}
    case token.SUB:
        return &object.Float{Value: lhs.Value - rhs.Value}
    case token.MULT:
        return &object.Float{Value: lhs.Value * rhs.Value}
    case token.DIV:
        return &object.Float{Value: lhs.Value / rhs.Value}
    case token.LT:
        return boolToBoolean(lhs.Value < rhs.Value)
    case token.GT:
        return boolToBoolean(lhs.Value > rhs.Value)
    case token.LE:
        return boolToBoolean(lhs.Value <= rhs.Value)
    case token.GE:
        return boolToBoolean(lhs.Value >= rhs.Value)
    case token.EQ:
        return boolToBoolean(lhs.Value == rhs.Value)
    case token.NEQ:
        return boolToBoolean(lhs.Value != rhs.Value)
    default:
        return makeError("unsupported infix operator on floats")
    }
}

func isTruthy(obj object.Object) bool {
    switch obj {
    case FALSE:
        return false
    case NULL:
        return false
    default:
        return true
    }
}

func computeNegExpr(obj object.Object) object.Object {
    if isTruthy(obj) {
        return FALSE
    } else {
        return TRUE
    }
}

func isError(obj object.Object) bool {
    return obj.Type() == object.ERROR_OBJECT || obj.Type() == object.PARSER_ERRORS_OBJECT
}

func isReturn(obj object.Object) bool {
    return obj.Type() == object.RETURN_OBJECT
}

func isBreak(obj object.Object) bool {
    return obj.Type() == object.BREAK_OBJECT
}

func isContinue(obj object.Object) bool {
    return obj.Type() == object.CONTINUE_OBJECT
}

func isBreakOrContinue(obj object.Object) bool {
    return isBreak(obj) || isContinue(obj)
}

func isErrorOrReturn(obj object.Object) bool {
    return isError(obj) || isReturn(obj)
}

func boolToBoolean(value bool) *object.Boolean {
    if value {
        return TRUE
    } else {
        return FALSE
    }
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
    env := object.NewEnclosingEnvironment(fn.Env)
    
    for i, p := range fn.Parameters {
        env.AddConst(p, args[i])
    }

    return env
}

func unwrapReturnValue(obj object.Object) object.Object {
    if retVal, ok := obj.(*object.Return); ok {
        return retVal.Value
    }
    return obj
}

func makeError(format string, a ...interface{}) *object.Error {
    return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func makeParserErrors(errs []error) *object.ParserErrors {
    return &object.ParserErrors{Errors: errs}
}
