package visitor

// import "os"
import "fmt"
import "tree"
import "go/ast"
import "reflect"
import "walk"

type AST_Visitor struct {
    parent *tree.Node
    name string
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
    if n == nil {
        if f, ok := finalizers[self.name]; ok {
            f(self.parent)
        }
        return nil
    }
    w := new(AST_Visitor)
    w.parent = self.addKid(self.getlabel(n))
    w.name = w.parent.Label
    return w
}

func (self *AST_Visitor) getlabel(n ast.Node) *tree.Node {
    node_type := reflect.TypeOf(n).String()
    name := node_type[5:]

    if f, ok := visitors[name]; ok {
        return f(name, n)
    } else if f, ok := visitors[node_type]; ok {
        return f(node_type, n)
    } else {
        return tree.NewNode(name)
    }
    panic("unreachable")
}

var visitors = map[string]func(string, ast.Node) *tree.Node{
    "Ident": func(name string, n ast.Node) *tree.Node {
        m := n.(*ast.Ident)
        return tree.NewNode(name).AddKid(tree.NewNode(m.Name))
    },

    "BasicLit": func(name string, n ast.Node) *tree.Node {
        m := n.(*ast.BasicLit)
        parent := tree.NewNode(name)
        parent.
            AddKid(tree.NewNode(fmt.Sprint(m.Kind))).
            AddKid(tree.NewNode(m.Value))
        return parent
    },

    "BinaryExpr": func(name string, n ast.Node) *tree.Node {
        m := n.(*ast.BinaryExpr)
        return tree.NewNode(name).
            AddKid(tree.NewNode(fmt.Sprint(m.Op)))
    },

    "UnaryExpr": func(name string, n ast.Node) *tree.Node {
        m := n.(*ast.UnaryExpr)
        return tree.NewNode(name).
            AddKid(tree.NewNode(fmt.Sprint(m.Op)))
    },

    "AssignStmt": func(name string, n ast.Node) *tree.Node {
        m := n.(*ast.AssignStmt)
        return tree.NewNode(name).
            AddKid(tree.NewNode(fmt.Sprint(m.Tok)))
    },

    "*walk.DummyNode": func(name string, n ast.Node) *tree.Node {
        m := n.(*walk.DummyNode)
        return tree.NewNode(m.Name)
    },
}

var finalizers = map[string]func(*tree.Node) {
    "Idents": func(root *tree.Node) {
        root.Children = func() []*tree.Node {
            children := make([]*tree.Node, 0, len(root.Children))
            for _, c := range root.Children {
                children = append(children, c.Children[0])
            }
            return children
        }()
    },

    "Type": func(root *tree.Node) {
        root.Children = root.Children[0].Children
    },

    "ElemType" : func(root *tree.Node) {
        root.Children = root.Children[0].Children
    },
}
