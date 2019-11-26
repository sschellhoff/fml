# Friendly Multi-paradigm Language
It is highly inspired by [Bob Nystrom](https://twitter.com/munificentbob)s [crafting interpreters](https://craftinginterpreters.com/) and [Thorsten Ball](https://twitter.com/thorstenball)s [Writing an interpreter in go](https://interpreterbook.com/)

You can start the program without any arguments to get a really simple repl. Or you can add a filepath as an argument to run a file of code.

## Buildin
Set GOPATH properly to the starting directory. Then run make in the code directory.

On linux, run `export GOPATH=(pwd)`, then go to the code directory `cd src/language` and run the makefile `make`.\
To run the interpreters REPL: `./interpreter`, to run a file, run `./interpreter filepath`. For example: `./interpreter examples/project_euler_001.fml`.

## Examples
[src/language/examples](https://github.com/sschellhoff/fml/tree/master/src/language/examples)


## Comming soon
* module system
* plugins (for own code wrappers and stuff)
* more tests
* lots of refactoring
