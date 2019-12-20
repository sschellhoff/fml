package scanner

import (
    "fmt"
    "bytes"
    "strconv"
    "unicode"
    "errors"
    "language/token"
)

type Scanner struct {
    sourcecode []rune
    start_idx int
    current_idx int
    start_last_line int
    line_counter int
    filepath string
}

const nullString string = "\x00"

func New(sourcecode string) *Scanner {
    return &Scanner{sourcecode: []rune(sourcecode), start_idx: 0, current_idx: 0, start_last_line: 0, line_counter: 1}
}

func (s *Scanner) NextToken() token.Token {
    s.skipWhitespace()

    s.start_idx = s.current_idx

    if s.isAtEnd() {
        return s.createEOF()
    }

    switch current_character := s.advance(); current_character {
    case "+":
        if s.match("=") {
            return s.createToken(token.ADDASSIGN)
        }
        return s.createToken(token.ADD)
    case "-":
        if s.match("=") {
            return s.createToken(token.SUBASSIGN)
        }
        return s.createToken(token.SUB)
    case "*":
        if s.match("=") {
            return s.createToken(token.MULTASSIGN)
        }
        return s.createToken(token.MULT)
    case "/":
        if s.match("=") {
            return s.createToken(token.DIVASSIGN)
        }
        if s.match("/") {
            if err := s.readLineComment(); err != nil {
                return s.createError(err.Error())
            }
            return s.NextToken()
        } else if s.match("*") {
            if err := s.readNestedMultilineComment(); err != nil {
                return s.createError(err.Error())
            }
            return s.NextToken()
        }
        return s.createToken(token.DIV)
    case "%":
        if s.match("=") {
            return s.createToken(token.MODASSIGN)
        }
        return s.createToken(token.MOD)
    case "=":
        if s.match("=") {
            return s.createToken(token.EQ)
        }
        return s.createToken(token.ASSIGN)
    case "!":
        if s.match("=") {
            return s.createToken(token.NEQ)
        }
        return s.createToken(token.NEG)
    case "<":
        if s.match("=") {
            return s.createToken(token.LE)
        }
        return s.createToken(token.LT)
    case ">":
        if s.match("=") {
            return s.createToken(token.GE)
        }
        return s.createToken(token.GT)
    case "&":
        if !s.match("&") {
            return s.createError("unexpectected &")
        }
        return s.createToken(token.AND)
    case "|":
        if !s.match("|") {
            return s.createError("unexpected |")
        }
        return s.createToken(token.OR)
    case "(":
        return s.createToken(token.LPAREN)
    case ")":
        return s.createToken(token.RPAREN)
    case "[":
        return s.createToken(token.LBRACKET)
    case "]":
        return s.createToken(token.RBRACKET)
    case "{":
        return s.createToken(token.LBRACE)
    case "}":
        return s.createToken(token.RBRACE)
    case ".":
        if s.match(".") {
            return s.createToken(token.RANGE)
        }
        return s.createToken(token.DOT)
    case ",":
        return s.createToken(token.COMMA)
    case ";":
        return s.createToken(token.SEMICOLON)
    case ":":
        return s.createToken(token.COLON)
    case "?":
        if s.match("?") {
            return s.createToken(token.NULLCOAL)
        }
        return s.createToken(token.QUESTION)
    case "_":
        s.readIdentifier()
        return s.createTokenWithLiteral(token.IDENTIFIER)
    case "\"":
        string_literal, err := s.readString()
        if err != nil {
            return s.createError(err.Error())
        }
        string_token := s.createToken(token.STRING)
        string_token.Literal = string_literal
        return string_token
    case nullString:
        return s.createEOF()
    default:
        current_rune := []rune(current_character)[0]
        if unicode.IsDigit(current_rune) {
            s.readNumber()
            if s.readFloat() {
                return s.createTokenWithLiteral(token.FLOAT)
            }
            return s.createTokenWithLiteral(token.INT)
        } else if unicode.IsLetter(current_rune) {
            s.readIdentifier()
            t := token.TypeFromIdent(s.getLexeme())
            if t == token.IDENTIFIER {
                return s.createTokenWithLiteral(token.IDENTIFIER)
            } else {
                return s.createToken(t)
            }
        }
        return s.createUnexpected()
    }
}

func (s *Scanner) advance() string {
    if s.isAtEnd() {
        return nullString // 0-char '\0'
    }
    if s.peek() == "\n" {
        s.start_last_line = s.current_idx + 1
        s.line_counter++
    }
    s.current_idx++
    return string(s.sourcecode[s.current_idx-1])
}

func (s *Scanner) peek() string {
    if s.isAtEnd() {
        return nullString
    }
    return string(s.sourcecode[s.current_idx])
}

func (s *Scanner) peek2() string {
    wantedIdx := s.current_idx + 1
    if wantedIdx >= len(s.sourcecode) {
        return nullString
    }
    return string(s.sourcecode[wantedIdx])
}

func (s *Scanner) match(c string) bool {
    if s.peek() == c {
        s.advance()
        return true
    }
    return false
}

func (s *Scanner) match2(c0 string, c1 string) bool {
    if s.peek() == c0 && s.peek2() == c1 {
        s.advance()
        s.advance()
        return true
    }
    return false
}

func (s *Scanner) getLexeme() string {
    return string(s.sourcecode[s.start_idx:s.current_idx])
}

func (s *Scanner) isAtEnd() bool {
    return s.current_idx == len(s.sourcecode)
}

func (s *Scanner) isNum() bool {
    return unicode.IsDigit([]rune(s.peek())[0])
}

func (s *Scanner) isHex() bool {
    r := []rune(s.peek())[0]
    return s.isNum() || (r >= 'A' && r <= 'F') || (r >= 'a' && r <= 'f')
}

func (s *Scanner) isNum2() bool {
    return unicode.IsDigit([]rune(s.peek2())[0])
}

func (s *Scanner) isAlpha() bool {
    return unicode.IsLetter([]rune(s.peek())[0])
}

func (s *Scanner) isAlphaNum() bool {
    return s.isAlpha() || s.isNum()
}

func (s *Scanner) isAlphaNumUnderscore() bool {
    return s.isAlphaNum() || s.isUnderscore()
}

func (s *Scanner) isUnderscore() bool {
    return s.peek() == "_"
}

func (s *Scanner) createToken(tokenType token.TokenType) token.Token {
    return token.FromType(tokenType, s.line_counter, s.currentColumn())
}

func (s *Scanner) createTokenWithLiteral(tokenType token.TokenType) token.Token {
    return token.New(tokenType, string(s.sourcecode[s.start_idx:s.current_idx]), s.line_counter, s.currentColumn())
}

func (s *Scanner) createError(msg string) token.Token {
    return token.New(token.ERROR, msg, s.line_counter, s.currentColumn())
}

func (s *Scanner) createUnexpected() token.Token {
    msg := fmt.Sprintf("unexpected lexeme '%s'", string(s.sourcecode[s.start_idx:s.current_idx]))
    return s.createError(msg)
}

func (s *Scanner) createEOF() token.Token {
    return token.FromType(token.EOF, s.line_counter, s.currentColumn())
}

func (s *Scanner) currentColumn() int {
    return s.start_idx - s.start_last_line + 1
}

func (s *Scanner) readNumber() {
    for s.isNum() {
        s.advance()
    }
}

func (s *Scanner) readFloat() bool {
    if s.peek() == "." && s.isNum2() {
        s.advance()
        s.advance()
        s.readNumber()
        return true
    }
    return false
}

func (s *Scanner) readIdentifier() {
    for s.isAlphaNumUnderscore() {
        s.advance()
    }
}

func (s *Scanner) readString() (string, error) {
    var out bytes.Buffer
    c := s.advance()
    for c != "\"" {
        if c == "\\" {
            e := s.advance()
            switch e {
            case "\"":
                out.WriteString("\"")
            case "\\":
                out.WriteString("\\")
            case "n":
                out.WriteString("\n")
            case "t":
                out.WriteString("\t")
            case "u":
                var hex bytes.Buffer
                hex.WriteString("'\\u")
                for i := 0; i < 4; i++ {
                    if s.isHex() {
                        hex.WriteString(s.advance())
                    } else {
                        return "", errors.New("Expected unicode sequence of length 4")
                    }
                }
                hex.WriteString("'")
                unquot, err := strconv.Unquote(hex.String())
                if err != nil {
                    return "", errors.New("cannot convert unicode sequence")
                }
                out.WriteString(unquot)
            case "U":
                var hex bytes.Buffer
                hex.WriteString("'\\U")
                for i := 0; i < 8; i++ {
                    if s.isHex() {
                        hex.WriteString(s.advance())
                    } else {
                        return "", errors.New("Expected unicode sequence of length 8")
                    }
                }
                hex.WriteString("'")
                unquot, err := strconv.Unquote(hex.String())
                if err != nil {
                    return "", errors.New("cannot convert unicode sequence")
                }
                out.WriteString(unquot)
            default:
                return "", errors.New("unexpected escape sequence")
            }
        } else {
            out.WriteString(c)
        }
        if s.isAtEnd() {
            return "", errors.New("unexpected end of file in string")
        }
        c = s.advance()
    }

    return string(out.String()), nil
}

func (s *Scanner) readLineComment() error {
    for !s.match("\n") {
        if s.isAtEnd() {
            return errors.New("unexpected end of file in comment")
        }
        s.advance()
    }
    return nil
}

func (s *Scanner) readNestedMultilineComment() error {
    for true {
        if s.isAtEnd() {
            return errors.New("unexpected end of file in multiline comment")
        }
        if s.match2("*", "/") {
            return nil
        } else if s.match2("/", "*") {
            if err := s.readNestedMultilineComment(); err != nil {
                return err
            }
        }
        s.advance()
    }
    return nil // never called
}

func (s *Scanner) skipWhitespace() {
    for unicode.IsSpace([]rune(s.peek())[0]) {
        s.advance()
    }
}
