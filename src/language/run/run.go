package run

import (
    "fmt"
    "io/ioutil"
    "language/scanner"
    "language/parser"
    "language/eval"
    "language/object"
    "language/ast"
)

func Run(filepath string) {
    code, err := readFile(filepath)
    if err != nil {
        fmt.Printf("There was an error while reading the input: %s", err.Error())
        return
    }
    environment := object.NewEnvironment()
    program, errors := parse(code)

    if len(errors) > 0 {
        printErrors(errors)
        return
    }
    evaluate(program, environment)
}

func parse(code string) (*ast.Program, []error) {
    s := scanner.New(code)
    p := parser.New(s)
    return  p.Parse()
}

func evaluate(program *ast.Program, env *object.Environment) {
    result := eval.Eval(program, env)
    if result.Type() == object.ERROR_OBJECT {
        fmt.Printf("\t%s\n", result.String())
    }
}

func printErrors(errors []error) {
    for _, msg := range errors {
        fmt.Printf("\t%s\n", msg.Error())
    }
}

func readFile(filepath string) (string, error) {
    content, err := ioutil.ReadFile(filepath)
    if err != nil {
        return "", err
    }
    return string(content), nil
}
