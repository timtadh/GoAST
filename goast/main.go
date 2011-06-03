package main

import "fmt"
import "os"
import "flag"
import "goast"

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
    pack := flag.String("pack", "", "The package you want the AST of.")
    flag.Parse()
    args := flag.Args()
    if *pack == "" {
        error("A package name is required.")
        usage(exitcodes["usage"])
    }
    if len(args) < 1 {
        error("A path to a go file is required.")
        usage(exitcodes["usage"])
    } else if len(args) > 1 {
        error("To many arguments, only one file is allowed.")
        usage(exitcodes["usage"])
    }

    path := args[0]
    ast, err := goast.ParsePackage(path, *ext, *pack)
    if err != nil {
        error(err.String())
        usage(exitcodes["usage"])
    }
    fmt.Println(ast.Dotty())
    error("GoAST complete")
}
