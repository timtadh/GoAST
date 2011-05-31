package visitor

import "fmt"
import "tree"
import "go/ast"
import "reflect"

type AST_Visitor struct {
    parent *tree.Node
}

func New() *AST_Visitor {
    v := new(AST_Visitor)
    v.parent = tree.NewNode("dummy")
    return v
}

func (self *AST_Visitor) AST() *tree.Node {
    return self.parent.Children[0]
}

func (self *AST_Visitor) addKid(name string) *tree.Node {
    node := tree.NewNode(name)
    if self.parent != nil {
        self.parent.AddKid(node)
    }
    return node
}

func (self *AST_Visitor) Visit(n ast.Node) ast.Visitor {
    fmt.Println(n)
    if n == nil { return nil }
    w := new(AST_Visitor)
    w.parent = self.addKid(reflect.TypeOf(n).String()[5:])
    return w
}
