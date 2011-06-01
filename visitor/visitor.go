package visitor

// import "os"
import "fmt"
import "tree"
import "go/ast"
import "reflect"
import "walk"

type AST_Visitor struct {
    parent *tree.Node
    node   ast.Node
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
        if f, ok := finalizers[self.parent.Label]; ok {
            f(self.node, self.parent)
        }
        return nil
    }
    w := new(AST_Visitor)
    w.parent = self.addKid(self.getlabel(n))
    w.node = n
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
        p := tree.NewNode(name).AddKid(tree.NewNode(m.Name))
        if m.IsExported() {
            p.AddKid(tree.NewNode("Exported"))
        }
        return p
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

    "BranchStmt": func(name string, n ast.Node) *tree.Node {
        m := n.(*ast.BranchStmt)
        return tree.NewNode(name).
            AddKid(tree.NewNode(fmt.Sprint(m.Tok)))
    },

    "IncDecStmt": func(name string, n ast.Node) *tree.Node {
        m := n.(*ast.IncDecStmt)
        return tree.NewNode(name).
            AddKid(tree.NewNode(fmt.Sprint(m.Tok)))
    },

    "*walk.DummyNode": func(name string, n ast.Node) *tree.Node {
        m := n.(*walk.DummyNode)
        return tree.NewNode(m.Name)
    },
}

/*
These functions rewrite the tree after construction. I try to keep these to a minimum. But
sometimes they are necessary for a cleaner tree.
*/
var finalizers = map[string]func(ast.Node, *tree.Node){
    "Idents": func(n ast.Node, root *tree.Node) {
        root.Children = func() []*tree.Node {
            children := make([]*tree.Node, 0, len(root.Children))
            for _, c := range root.Children {
                children = append(children, c.Children[0])
            }
            return children
        }()
    },

    "LabeledStmt": func(n ast.Node, root *tree.Node) {
        root.Children[0].Label = "Label"
    },

    "CallExpr": func(n ast.Node, root *tree.Node) {
        m := n.(*ast.CallExpr)
        if m.Ellipsis.IsValid() {
            root.AddKid(tree.NewNode("HasEllipsis"))
        }
    },

    "Unresolved": func(n ast.Node, root *tree.Node) {
        newkids := make([]*tree.Node, 0, len(root.Children))
        set := make(map[string]bool)
        for _, ident := range root.Children {
            name := ident.Children[0].Label
            if _, has := set[name]; !has {
                set[name] = true
                newkids = append(newkids, tree.NewNode(name))
            }
        }
        root.Children = newkids
    },
}
