package main

import "fmt"
import "os"
import "flag"
import "strings"
import "go/parser"
import "go/token"
// import "go/ast"
import "visitor"
import "walk"
// import "tree"
// import "reflect"

const usage_msg = `goast [-ext=.go] path

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
    ext := flag.String("ext", ".go", "The extension of the files to parse.")
    flag.Parse()
    args := flag.Args()
    if len(args) < 1 {
        error("A path to a go file is required.")
        usage(exitcodes["usage"])
    } else if len(args) > 1 {
        error("To many arguments, only one file is allowed.")
        usage(exitcodes["usage"])
    }

    path := args[0]
    if pkgs, err := parser.ParseDir(token.NewFileSet(), path,
        func(finfo *os.FileInfo) bool {
            if strings.HasSuffix(finfo.Name, *ext) {
                error("found " + finfo.Name)
            }
            return strings.HasSuffix(finfo.Name, *ext)
        },
        0); err != nil {
        error("could not read path")
        error(err.String())
        usage(exitcodes["badpath"])
    } else {
        if len(pkgs) == 0 {
            error("no files found.")
            usage(exitcodes["badpath"])
        }
        //         if len(pkg) != 1 {
        //             error("the directory contained more than one package, I will barf.")
        //             os.Exit(exitcodes["barf"])
        //         }
        for _, node := range pkgs {
            visitor := visitor.New()
            walk.GoAST_Walk(
                visitor,
                node,
            )
            AST := visitor.AST()
            fmt.Println(AST.Dotty())
            break
        }
    }
}
