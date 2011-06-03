package goast

import "os"
import "strings"
import "go/parser"
import "go/token"
// import "go/ast"
import "goast/visitor"
import "goast/walk"
import "goast/tree"

func ParseFile(file_path string) (*tree.Node, bool) {
    file, err := parser.ParseFile(token.NewFileSet(), file_path, nil, 0)
    if err != nil {
        return nil, false
    }
    visitor := visitor.New()
    walk.GoAST_Walk(
        visitor,
        file,
    )
    return visitor.AST(), true
}

func ParsePackage(dir_path, ext, package_name string) (*tree.Node, bool) {
    pkgs, ok := ParseDirectory(dir_path, ext)
    if !ok {
        return nil, false
    }
    if pkg, has := pkgs[package_name]; has {
        return pkg, true
    }
    return nil, false
}

func ParseDirectory(dir_path, ext string) (map[string]*tree.Node, bool) {
    pkgs, err := parser.ParseDir(
        token.NewFileSet(),
        dir_path,
        func(finfo *os.FileInfo) bool {
            return strings.HasSuffix(finfo.Name, ext)
        },
        0)
    if err != nil {
        return nil, false
    }
    if len(pkgs) == 0 {
        return nil, false
    }
    pkgasts := make(map[string]*tree.Node)
    for name, node := range pkgs {
        visitor := visitor.New()
        walk.GoAST_Walk(
            visitor,
            node,
        )
        pkgasts[name] = visitor.AST()
    }
    return pkgasts, true
}
