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

func (self *AST_Visitor) addKid(node *tree.Node) *tree.Node {
    if self.parent != nil {
        self.parent.AddKid(node)
    }
    return node
}

func (self *AST_Visitor) Visit(n ast.Node) ast.Visitor {
    if n == nil { return nil }
    w := new(AST_Visitor)
    w.parent = self.addKid(getlabel(n))
    return w
}

func getlabel(n ast.Node) *tree.Node {
    node_type := reflect.TypeOf(n).String()[5:]
    if f, ok := visitors[node_type]; !ok {
        return tree.NewNode(node_type)
    } else {
        return f(node_type, n)
    }
    panic("unreachable")
}

var visitors = map[string]func(string, ast.Node)*tree.Node {
    "Ident" : func (node_type string, n ast.Node) *tree.Node  {
        m := n.(*ast.Ident)
        parent := tree.NewNode(node_type)
        parent.AddKid(tree.NewNode(m.Name))
        if m.Obj != nil {
            parent.AddKid(
                func(obj *ast.Object) *tree.Node {
                    parent := tree.NewNode("Object").
                        AddKid(tree.NewNode(obj.Kind.String()))
                    if obj.Decl != nil {
                        parent.AddKid(
                            func(d interface{}) *tree.Node {
                                parent := tree.NewNode("Decl")
                                switch decl := d.(type) {
                                case *ast.Field:
                                case *ast.FuncDecl:
                                    parent.
                                        AddKid(getlabel(decl).
                                            AddKid(getlabel(decl.Name)),
                                        )
                                case *ast.LabeledStmt:
                                case *ast.ImportSpec:
                                case *ast.ValueSpec:
                                case *ast.TypeSpec:
                                }
                                return parent
                            }(obj.Decl),
                        )
                    }
                    return parent
                }(m.Obj),
            )
        }
        return parent
    },

    "BasicLit" : func (node_type string, n ast.Node) *tree.Node  {
        m := n.(*ast.BasicLit)
        parent := tree.NewNode(node_type)
        parent.
            AddKid(tree.NewNode(fmt.Sprint(m.Kind))).
            AddKid(tree.NewNode(m.Value))
        return parent
    },
}

