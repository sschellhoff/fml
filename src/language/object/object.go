package object

import (
    "fmt"
    "bytes"
    "strings"
    "hash/fnv"
    "language/ast"
)

const (
    INTEGER_OBJECT = "INTEGER"
    FLOAT_OBJECT = "FLOAT"
    BOOLEAN_OBJECT = "BOOL"
    STRING_OBJECT = "STRING"
    NULL_OBJECT = "NULL"
    RETURN_OBJECT = "RETURN"
    BREAK_OBJECT = "BREAK"
    CONTINUE_OBJECT = "CONTINUE"
    FUNCTION_OBJECT = "FUNCTION"
    BUILTIN_OBJECT = "BUILTIN"
    ARRAY_OBJECT = "ARRAY"
    HASH_OBJECT = "HASH"
    MODULE_OBJECT = "MODULE"
    ERROR_OBJECT = "ERROR"
    PARSER_ERRORS_OBJECT = "PARSERERRORS"
)

type ObjectType string

type Object interface {
    Type() ObjectType
    String() string
}

type Hashable interface {
    HashKey() HashKey
}

type Integer struct {
    Value int64
}

func (i *Integer) Type() ObjectType {
    return INTEGER_OBJECT
}

func (i *Integer) String() string {
    return fmt.Sprintf("%d", i.Value)
}

func (i *Integer) HashKey() HashKey {
    return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}


type Float struct {
    Value float64
}

func (f *Float) Type() ObjectType {
    return FLOAT_OBJECT
}

func (f *Float) String() string {
    return fmt.Sprintf("%f", f.Value)
}


type Boolean struct {
    Value bool
}

func (b *Boolean) Type() ObjectType {
    return BOOLEAN_OBJECT
}

func (b *Boolean) String() string {
    return fmt.Sprintf("%t", b.Value)
}

func (b *Boolean) HashKey() HashKey {
    var value uint64
    if b.Value {
        value = 1
    } else {
        value = 0
    }
    return HashKey{Type: b.Type(), Value: value}
}


type String struct {
    Value string
}

func (s *String) Type() ObjectType {
    return STRING_OBJECT
}

func (s *String) String() string {
    return s.Value
}

func (s *String) HashKey() HashKey {
    h := fnv.New64a()
    h.Write([]byte(s.Value))

    return HashKey{Type: s.Type(), Value: h.Sum64()}
}


type Null struct {
}

func (n *Null) Type() ObjectType {
    return NULL_OBJECT
}

func (n *Null) String() string {
    return "null"
}


type Return struct {
    Value Object
}

func (r *Return) Type() ObjectType {
    return RETURN_OBJECT
}

func (r *Return) String() string {
    return r.Value.String()
}


type Break struct {
}

func (b *Break) Type() ObjectType {
    return BREAK_OBJECT
}

func (b *Break) String() string {
    return "break;"
}


type Continue struct {
}

func (c *Continue) Type() ObjectType {
    return CONTINUE_OBJECT
}

func (c *Continue) String() string {
    return "continue;"
}


type Error struct {
    Message string
}

func (e *Error) Type() ObjectType {
    return ERROR_OBJECT
}

func (e *Error) String() string {
    return "ERROR: " + e.Message
}


type ParserErrors struct {
    Errors []error
}

func (p *ParserErrors) Type() ObjectType {
    return PARSER_ERRORS_OBJECT
}

func (p *ParserErrors) String() string {
    var out bytes.Buffer

    strValues := []string{"Parser errors:"}

    for _, e := range p.Errors {
        strValues = append(strValues, e.Error())
    }

    out.WriteString(strings.Join(strValues, "\n"))

    return out.String()
}


type Function struct {
    Parameters []string
    Body *ast.BlockStatement
    Env *Environment
}

func (f *Function) Type() ObjectType {
    return FUNCTION_OBJECT
}

func (f *Function) String() string {
    var out bytes.Buffer

    out.WriteString("fun(")
    out.WriteString(strings.Join(f.Parameters, ", "))
    out.WriteString(")")
    out.WriteString(f.Body.String())

    return out.String()
}


type BuiltinFunction func(args ...Object) Object

type Builtin struct {
    Function BuiltinFunction
}

func (b *Builtin) Type() ObjectType {
    return BUILTIN_OBJECT
}

func (b *Builtin) String() string {
    return "builtin function"
}


type Array struct {
    Elements []Object
}

func (a *Array) Type() ObjectType {
    return ARRAY_OBJECT
}

func (a *Array) String() string {
    var out bytes.Buffer

    elements := []string{}

    for _, e := range a.Elements {
        elements = append(elements, e.String())
    }

    out.WriteString("[")
    out.WriteString(strings.Join(elements, ", "))
    out.WriteString("]")

    return out.String()
}


type HashKey struct {
    Type ObjectType
    Value uint64
}


type HashPair struct {
    Key Object
    Value Object
}


type Hash struct {
    Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType {
    return HASH_OBJECT
}

func (h *Hash) String() string {
    var out bytes.Buffer

    pairs := []string{}
    for _, pair := range h.Pairs {
        pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.String(), pair.Value.String()))
    }

    out.WriteString("{")
    out.WriteString(strings.Join(pairs, ", "))
    out.WriteString("}")

    return out.String()
}


type Module struct {
    Path string
    Env *Environment
}

func (m *Module) Type() ObjectType {
    return MODULE_OBJECT
}

func (m *Module) String() string {
    return m.Path
}
