package goast

import "os"
import "strings"
import "go/parser"
import "go/token"
// import "go/ast"
import "goast/visitor"
import "goast/walk"
import "goast/tree"

func ParseFile(file_path string) (*tree.Node, os.Error) {
    file, err := parser.ParseFile(token.NewFileSet(), file_path, nil, 0)
    if err != nil {
        return nil, err
    }
    visitor := visitor.New()
    walk.GoAST_Walk(
        visitor,
        file,
    )
    return visitor.AST(), nil
}

func ParsePackage(dir_path, ext, package_name string) (*tree.Node, os.Error) {
    pkgs, err := ParseDirectory(dir_path, ext)
    if err != nil {
        return nil, err
    }
    if pkg, has := pkgs[package_name]; has {
        return pkg, nil
    }
    return nil, os.NewError("The supplied package name was not found.")
}

func ParseDirectory(dir_path, ext string) (map[string]*tree.Node, os.Error) {
    pkgs, err := parser.ParseDir(
        token.NewFileSet(),
        dir_path,
        func(finfo *os.FileInfo) bool {
            return strings.HasSuffix(finfo.Name, ext)
        },
        0)
    if err != nil {
        return nil, err
    }
    if len(pkgs) == 0 {
        return nil, os.NewError("No packages found.")
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
    return pkgasts, nil
}
