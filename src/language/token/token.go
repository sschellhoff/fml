package token

type TokenType string

type Token struct {
    Type TokenType
    Literal string
    Line int
    Column int
}

func New(the_type TokenType, literal string, line, column int) Token {
    return Token{Type: the_type, Literal: literal, Line: line, Column: column}
}

func FromType(the_type TokenType, line, column int) Token {
    return New(the_type, "", line, column)
}

func Compare(t1, t2 Token) bool {
    return t1.Type == t2.Type && t1.Literal == t2.Literal
}

const (
    EOF = "EOF"
    ERROR = "ERROR"

    IDENTIFIER = "IDENTIFIER"
    INT = "INT"
    FLOAT = "FLOAT"
    STRING = "STRING"
    TRUE = "TRUE"
    FALSE = "FALSE"
    NULL = "NULL"

    LET = "LET"
    CONST = "CONST"
    IF = "IF"
    ELSE = "ELSE"
    LOOP = "LOOP"
    IN = "IN"
    FUN = "FUN"
    RETURN = "RETURN"
    BREAK = "BREAK"
    CONTINUE = "CONTINUE"
    FOREVER = "FOREVER"
    TRY = "TRY"
    CATCH = "CATCH"
    IMPORT = "IMPORT"
    AS = "AS"

    ADD = "+"
    SUB = "-"
    MULT = "*"
    DIV = "/"
    MOD = "%"
    AND = "&&"
    OR = "||"
    EQ = "=="
    NEQ = "!="
    LT = "<"
    GT = ">"
    LE = "<="
    GE = ">="
    ASSIGN = "="
    NEG = "!"
    COMMA = ","
    SEMICOLON = ";"
    COLON = ":"
    QUESTION = "?"
    NULLCOAL = "??"
    DOT = "."
    RANGE = ".."
    ADDASSIGN = "+="
    SUBASSIGN = "-="
    MULTASSIGN = "*="
    DIVASSIGN = "/="
    MODASSIGN = "%="

    LPAREN = "("
    RPAREN = ")"
    LBRACE = "{"
    RBRACE = "}"
    LBRACKET = "["
    RBRACKET = "]"
)

var reservedWords = map[string]TokenType{
    "let": LET,
    "const": CONST,
    "if": IF,
    "else": ELSE,
    "loop": LOOP,
    "in": IN,
    "fun": FUN,
    "return": RETURN,
    "break": BREAK,
    "continue": CONTINUE,
    "true": TRUE,
    "false": FALSE,
    "null": NULL,
    "and": AND,
    "or": OR,
    "forever": FOREVER,
    "try": TRY,
    "catch": CATCH,
    "import": IMPORT,
    "as": AS,
}

func TypeFromIdent(value string) TokenType {
    if t, ok := reservedWords[value]; ok {
        return t
    }
    return IDENTIFIER
}
