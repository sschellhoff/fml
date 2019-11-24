package main

import (
    "os"
    "fmt"
    "language/repl"
    "language/run"
)

func main() {
    cmdArgs := os.Args[1:]
    if len(cmdArgs) == 0 {
        repl.Start(os.Stdin, os.Stdout)
    } else if len(cmdArgs) == 1 {
        run.Run(cmdArgs[0])
    } else {
        fmt.Printf("You can only run this command with 0 or 1 arguments.\nIf you run it without arguments, you start the REPL\nIf you run it with one argument, it gets interpreted as a filepath and the file gets evalauted")
    }
}
