package frontend

import (
    "io/ioutil"
    "path/filepath"
    "language/scanner"
    "language/parser"
    "language/ast"
)

func Build(path string) (*ast.Program, []error) {
    path, err := filepath.Abs(path)
    if err != nil {
        return nil, []error{err}
    }
    code, err := readFile(path)
    if err != nil {
        return nil, []error{err}
    }
    return parse(code, path)
}

func readFile(path string) (string, error) {
    content, err := ioutil.ReadFile(path)
    if err != nil {
        return "", err
    }
    return string(content), nil
}

func parse(code string, path string) (*ast.Program, []error) {
    s := scanner.New(code)
    p := parser.New(s, path)
    return  p.Parse()
}

