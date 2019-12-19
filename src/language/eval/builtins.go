package eval

import (
    "strings"
    "strconv"
    "fmt"
    "bufio"
    "os"
    "language/object"
    "language/ast"
)

var builtins = map[string]*object.Builtin{
    "len": &object.Builtin{
        Function: func(args ...object.Object) object.Object {
            if len(args) != 1 {
                return makeBuiltinError("wrong number of arguments, want 1, got %d", len(args))
            }

            arg := args[0]
            switch value := arg.(type) {
            case *object.String:
                r := []rune(value.Value)
                return &object.Integer{Value: int64(len(r))}
            case *object.Array:
                return &object.Integer{Value: int64(len(value.Elements))}
            default:
                return makeBuiltinError("cannot call len on %s", value.Type())
            }
        },
    },
    "first": &object.Builtin{
        Function: func(args ...object.Object) object.Object {
            if len(args) != 1 {
                return makeBuiltinError("wrong number of arguments, want, got %d", len(args))
            }

            arg := args[0]
            switch value := arg.(type) {
            case *object.Array:
                if len(value.Elements) == 0 {
                    return makeBuiltinError("cannot get first argument of an empty array")
                }
                return value.Elements[0]
            default:
                return makeBuiltinError("cannot call first on %s", value.Type())
            }
        },
    },
    "last": &object.Builtin{
        Function: func(args ...object.Object) object.Object {
            if len(args) != 1 {
                return makeBuiltinError("wrong number of arguments, want, got %d", len(args))
            }

            arg := args[0]
            switch value := arg.(type) {
            case *object.Array:
                if len(value.Elements) == 0 {
                    return makeBuiltinError("cannot get last argument of an empty array")
                }
                return value.Elements[len(value.Elements)-1]
            default:
                return makeBuiltinError("cannot call last on %s", value.Type())
            }
        },
    },
    "rest": &object.Builtin{
        Function: func(args ...object.Object) object.Object {
            if len(args) != 1 {
                return makeBuiltinError("wrong number of arguments, want, got %d", len(args))
            }

            arg := args[0]
            switch value := arg.(type) {
            case *object.Array:
                length := len(value.Elements)
                if length == 0 {
                    return makeBuiltinError("cannot get rest of an empty array")
                }
                newElements := make([]object.Object, length-1, length-1)
                copy(newElements, value.Elements[1:length])
                return &object.Array{Elements: newElements}
            default:
                return makeBuiltinError("cannot call rest on %s", value.Type())
            }
        },
    },
    "push": &object.Builtin{
        Function: func(args ...object.Object) object.Object {
            if len(args) != 2 {
                return makeBuiltinError("wrong number of arguments, want 2, got %d", len(args))
            }

            arg := args[0]
            element := args[1]
            switch value := arg.(type) {
            case *object.Array:
                length := len(value.Elements)
                newElements := make([]object.Object, length+1, length+1)
                copy(newElements, value.Elements)
                newElements[length] = element
                return &object.Array{Elements: newElements}
            default:
                return makeBuiltinError("cannot call push on %s", value.Type())
            }
        },
    },
    "makeArray": &object.Builtin{
        Function: func(args ...object.Object) object.Object {
            if len(args) != 2 {
                return makeBuiltinError("wrong number of arguments, want 2, got %d", len(args))
            }

            lengthObj, ok := args[0].(*object.Integer)
            if !ok {
                return makeBuiltinError("first argument must be of type integer, got %s", args[0].Type())
            }
            length := lengthObj.Value
            value := args[1]
            elements := make([]object.Object, length)
            for i := range elements {
                elements[i] = value
            }
            return &object.Array{Elements: elements}
        },
    },
    "print": &object.Builtin{
        Function: func(args ...object.Object) object.Object {
            var strs = []string{}
            for _, s := range args {
                strs = append(strs, s.String())
            }
            fmt.Printf("%s", strings.Join(strs, ", "))

            return NULL
        },
    },
    "println": &object.Builtin{
        Function: func(args ...object.Object) object.Object {
            var strs = []string{}
            for _, s := range args {
                strs = append(strs, s.String())
            }
            fmt.Printf("%s\n", strings.Join(strs, ", "))

            return NULL
        },
    },
    "readline": &object.Builtin{
        Function: func(args ...object.Object) object.Object {
            if len(args) != 1 {
                return makeBuiltinError("wrong number of arguments, want 1, got %d", len(args))
            }
            output, ok := args[0].(*object.String)
            if !ok {
                return makeBuiltinError("expected argument to be of type string")
            }
            fmt.Printf("%s", output.Value)
            scanner := bufio.NewScanner(os.Stdin)
            scanner.Scan()
            return &object.String{Value: scanner.Text()}

            return NULL
        },
    },
    "str": &object.Builtin{
        Function: func(args ...object.Object) object.Object {
            if len(args) != 1 {
                return makeBuiltinError("wrong number of arguments, want 1, got %d", len(args))
            }

            arg := args[0]
            return &object.String{Value: arg.String()}
        },
    },
    "int": &object.Builtin{
        Function: func(args ...object.Object) object.Object {
            if len(args) != 1 {
                return makeBuiltinError("wrong number of arguments, want 1, got %d", len(args))
            }

            arg := args[0]
            switch value := arg.(type) {
            case *object.String:
                s := strings.TrimSpace(value.Value)
                i, err := strconv.ParseInt(s, 10, 64)
                if err != nil {
                    return makeBuiltinError("Cannot convert string \"%s\" to integer", value.Value)
                }
                return &object.Integer{Value: i}
            case *object.Float:
                f := int64(value.Value)
                return &object.Integer{Value: f}
            case *object.Integer:
                return value
            default:
                return makeBuiltinError("cannot convert values of type %s to int", arg.Type())
            }
            return &object.String{Value: arg.String()}
        },
    },
    "float": &object.Builtin{
        Function: func(args ...object.Object) object.Object {
            if len(args) != 1 {
                return makeBuiltinError("wrong number of arguments, want 1, got %d", len(args))
            }

            arg := args[0]
            switch value := arg.(type) {
            case *object.String:
                s := strings.TrimSpace(value.Value)
                f, err := strconv.ParseFloat(s, 64)
                if err != nil {
                    return makeBuiltinError("Cannot convert string \"%s\" to float", value.Value)
                }
                return &object.Float{Value: f}
            case *object.Float:
                return value
            case *object.Integer:
                f := float64(value.Value)
                return &object.Float{Value: f}
            default:
                return makeBuiltinError("cannot convert values of type %s to int", arg.Type())
            }
            return &object.String{Value: arg.String()}
        },
    },
    "substring": &object.Builtin{
        Function: func(args ...object.Object) object.Object {
            if len(args) != 3 {
                return makeBuiltinError("wrong number of arguments, want 3, got %d", len(args))
            }

            arg := args[0]
            switch value := arg.(type) {
            case *object.String:
                start := args[1]
                end := args[2]
                startI, ok := start.(*object.Integer)
                if !ok {
                    return makeBuiltinError("Start index need to be an integer")
                }
                if startI.Value < 0 {
                    return makeBuiltinError("Start index must be >= 0")
                }
                endI, ok := end.(*object.Integer)
                if !ok {
                    return makeBuiltinError("End index need to be an integer")
                }
                strValue := value.Value
                r := []rune(strValue)
                if endI.Value > int64(len(r)) {
                    return makeBuiltinError("End index must be < than strings length")
                }
                if endI.Value < startI.Value {
                    return makeBuiltinError("End index must be >= start index")
                }
                return &object.String{Value: string(r[startI.Value:endI.Value])}
            default:
                return makeBuiltinError("cannot call substring on %s", value.Type())
            }
        },
    },
    "copy": &object.Builtin{
        Function: func(args ...object.Object) object.Object {
            if len(args) != 1 {
                return makeBuiltinError("wrong number of arguments, want 1, got %d", len(args))
            }

            arg := args[0]
            return shallowCopy(arg)
        },
    },
    "deepcopy": &object.Builtin{
        Function: func(args ...object.Object) object.Object {
            if len(args) != 1 {
                return makeBuiltinError("wrong number of arguments, want 1, got %d", len(args))
            }

            arg := args[0]
            return deepCopy(arg)
        },
    },
    "isInt": &object.Builtin{
        Function: func(args ...object.Object) object.Object {
            return isOfTypeHelper(object.INTEGER_OBJECT, args...)
        },
    },
    "isFloat": &object.Builtin{
        Function: func(args ...object.Object) object.Object {
            return isOfTypeHelper(object.INTEGER_OBJECT, args...)
        },
    },
    "isBool": &object.Builtin{
        Function: func(args ...object.Object) object.Object {
            return isOfTypeHelper(object.BOOLEAN_OBJECT, args...)
        },
    },
    "isString": &object.Builtin{
        Function: func(args ...object.Object) object.Object {
            return isOfTypeHelper(object.STRING_OBJECT, args...)
        },
    },
    "isArray": &object.Builtin{
        Function: func(args ...object.Object) object.Object {
            return isOfTypeHelper(object.ARRAY_OBJECT, args...)
        },
    },
    "isHash": &object.Builtin{
        Function: func(args ...object.Object) object.Object {
            return isOfTypeHelper(object.HASH_OBJECT, args...)
        },
    },
    "isFunction": &object.Builtin{
        Function: func(args ...object.Object) object.Object {
            return isOfTypeHelper(object.FUNCTION_OBJECT, args...)
        },
    },
    "isCallable": &object.Builtin{
        Function: func(args ...object.Object) object.Object {
            isFunc := isTruthy(isOfTypeHelper(object.FUNCTION_OBJECT, args...))
            isBuiltin := isTruthy(isOfTypeHelper(object.BUILTIN_OBJECT, args...))

            return boolToBoolean(isFunc || isBuiltin)
        },
    },
    "isBuiltin": &object.Builtin{
        Function: func(args ...object.Object) object.Object {
            return isOfTypeHelper(object.BUILTIN_OBJECT, args...)
        },
    },
    "error": &object.Builtin{
        Function: func(args ...object.Object) object.Object {
            argStrs := []string{}
            for _, arg := range args {
                argStrs = append(argStrs, arg.String())
            }
            return makeBuiltinError(strings.Join(argStrs, ", "))
        },
    },
}

func isOfTypeHelper(wantedType object.ObjectType, args ...object.Object) object.Object {
    if len(args) != 1 {
        return makeBuiltinError("wrong number of arguments, want 1, got %d", len(args))
    }

    arg := args[0]
    return boolToBoolean(arg.Type() == wantedType)
}

func shallowCopy(arg object.Object) object.Object {
    switch value := arg.(type) {
    case *object.Array:
        result := &object.Array{Elements: make([]object.Object, len(value.Elements))}
        for i, v := range value.Elements {
            result.Elements[i] = v
        }
        return result
    case *object.Hash:
        result := &object.Hash{Pairs: make(map[object.HashKey]object.HashPair)}
        for k, v := range value.Pairs {
            result.Pairs[k] = v
        }
        return result
    default:
        // no need to copy
        return arg
    }
}

func deepCopy(arg object.Object) object.Object {
    switch value := arg.(type) {
    case *object.Array:
        result := &object.Array{Elements: make([]object.Object, len(value.Elements))}
        for i, v := range value.Elements {
            result.Elements[i] = deepCopy(v)
        }
        return result
    case *object.Hash:
        result := &object.Hash{Pairs: make(map[object.HashKey]object.HashPair)}
        for k, v := range value.Pairs {
            result.Pairs[k] = object.HashPair{Key: v.Key, Value: deepCopy(v.Value)}
        }
        return result
    default:
        // no need to deepcopy
        return arg
    }
}

func makeBuiltinError(format string, a ...interface{}) *object.Error {
    return &object.Error{Message: fmt.Sprintf(format, a...), StackTrace: []ast.PositionalInfo{}}
}
