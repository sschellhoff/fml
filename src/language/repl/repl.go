package repl

import (
    "io"
    "bufio"
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "language/scanner"
    "language/parser"
    "language/ast"
    "language/eval"
    "language/object"
)

func Start(in io.Reader, out io.Writer) {
    // signal interrupt code
    channel := make(chan os.Signal)
    signal.Notify(channel, os.Interrupt, syscall.SIGTERM)
    go func() {
        <- channel
        cleanup()
        os.Exit(1)
    }()

    // main code starts here
    inputScanner := bufio.NewScanner(in)
    environment := object.NewEnvironment()


    for {
        fmt.Printf(PROMPT)

        scanned := inputScanner.Scan()

        if !scanned {
            return
        }

        code := inputScanner.Text()

        program, errors := parse(code)

        if len(errors) > 0 {
            printErrors(errors, out)
            continue
        }

        evaluate(program, environment, out)
    }
}

func parse(code string) (*ast.Program, []error) {
    s := scanner.New(code)
    p := parser.New(s)
    return  p.Parse()
}

func evaluate(program *ast.Program, env *object.Environment, out io.Writer) {
    evaluated := eval.Eval(program, env, make(map[string]*object.Module))

    io.WriteString(out, evaluated.String())

    io.WriteString(out, "\n")
}

func cleanup() {
    fmt.Printf("\nSee you soon!\n")
}

func printErrors(errors []error, out io.Writer) {
    for _, msg := range errors {
        io.WriteString(out, "\t"+msg.Error()+"\n")
    }
}

const PROMPT = "> "
