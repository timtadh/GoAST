package main

import "fmt"
import "os"
import "strings"
import "go/parser"
import "go/token"
import "go/ast"
import "visitor"
// import "tree"
// import "reflect"

const usage_msg = `goast gofile.go

GoAST produces an AST for the given go file.
`

var exitcodes = map[string]int{
    "ok":      0,
    "usage":   1,
    "badpath": 2,
    "barf":    3,
}

func usage(status int) {
    fmt.Fprint(os.Stderr, usage_msg)
    os.Exit(status)
}

func error(s string) {
    fmt.Fprintln(os.Stderr, s)
}

func main() {
    if len(os.Args[1:]) < 1 {
        error("A path to a go file is required.")
        usage(exitcodes["usage"])
    } else if len(os.Args[1:]) > 1 {
        error("To many arguments, only one file is allowed.")
        usage(exitcodes["usage"])
    }

    path := os.Args[1]
    fmt.Println("path = ", path)
    if pkg, err := parser.ParseDir(token.NewFileSet(), path,
        func(finfo *os.FileInfo) bool {
            fmt.Println(finfo.Name, strings.HasSuffix(finfo.Name, ".go"))
            return strings.HasSuffix(finfo.Name, ".go")
        },
        0); err != nil {
        error("could not read path")
        error(err.String())
        usage(exitcodes["badpath"])
    } else {
        if len(pkg) != 1 {
            error("the directory contained more than one package, I will barf.")
            os.Exit(exitcodes["barf"])
        }
        for name, node := range pkg {
            fmt.Println(name)
            visitor := visitor.New()
            ast.Walk(
                visitor,
                node,
            )
            AST := visitor.AST()
            fmt.Println(AST)
        }
    }
}
