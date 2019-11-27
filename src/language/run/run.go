package run

import (
    "fmt"
    "path/filepath"
    "language/eval"
    "language/object"
    "language/ast"
    "language/frontend"
)

func Run(path string) {
    program, errors := frontend.Build(path)

    if len(errors) > 0 {
        printErrors(errors)
        return
    }
    evaluate(program, path)
}

func evaluate(program *ast.Program, path string) {
    absPath, err := filepath.Abs(path)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    path = absPath
    eval.MODULEPATH = filepath.Dir(path)
    mainModuleEnv := object.NewEnvironment()
    mainModule := &object.Module{Path: path, Env: mainModuleEnv}
    modules := make(map[string]*object.Module)
    modules[path] = mainModule
    result := eval.Eval(program, mainModuleEnv, modules)
    if result.Type() == object.ERROR_OBJECT || result.Type() == object.PARSER_ERRORS_OBJECT {
        fmt.Printf("\t%s\n", result.String())
    }
}

func printErrors(errors []error) {
    for _, msg := range errors {
        fmt.Printf("\t%s\n", msg.Error())
    }
}
