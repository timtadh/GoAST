package walk

// import "os"
import "fmt"
import . "go/ast"
import "go/token"

type DummyNode struct {
    start token.Pos
    end   token.Pos
    Name  string
    List  []Node
}

func NewDummyNode(name string, start, end token.Pos, list []Node) *DummyNode {
    return &DummyNode{start: start, end: end, Name: name, List: list}
}
func (self *DummyNode) Pos() token.Pos { return self.start }
func (self *DummyNode) End() token.Pos { return self.end }

// Helper functions.
func walkIdentList(v Visitor, list []*Ident) {
    if len(list) == 0 {
        return
    }
    GoAST_Walk(v,
        NewDummyNode("Idents",
            list[0].Pos(),
            list[len(list)-1].End(),
            func() []Node {
                nodes := make([]Node, 0, len(list))
                for _, c := range list {
                    nodes = append(nodes, Node(c))
                }
                return nodes
            }()))
}


func walkExprList(v Visitor, list []Expr) {
    if len(list) == 0 {
        return
    }
    GoAST_Walk(v,
        NewDummyNode("Exprs",
            list[0].Pos(),
            list[len(list)-1].End(),
            func() []Node {
                nodes := make([]Node, 0, len(list))
                for _, c := range list {
                    nodes = append(nodes, Node(c))
                }
                return nodes
            }()))
}


func walkStmtList(v Visitor, list []Stmt) {
    if len(list) == 0 {
        return
    }
    GoAST_Walk(v,
        NewDummyNode("Stmts",
            list[0].Pos(),
            list[len(list)-1].End(),
            func() []Node {
                nodes := make([]Node, 0, len(list))
                for _, c := range list {
                    nodes = append(nodes, Node(c))
                }
                return nodes
            }()))
}


func walkDeclList(v Visitor, list []Decl) {
    if len(list) == 0 {
        return
    }
    GoAST_Walk(v,
        NewDummyNode("Decls",
            list[0].Pos(),
            list[len(list)-1].End(),
            func() []Node {
                nodes := make([]Node, 0, len(list))
                for _, c := range list {
                    nodes = append(nodes, Node(c))
                }
                return nodes
            }()))
}


// I am forking the version of walk in the go stdlib. Why this crazyness?
// because I want to insert nodes and this is the cleanest way to do it.

// Walk traverses an AST in depth-first order: It starts by calling
// v.Visit(node); node must not be nil. If the visitor w returned by
// v.Visit(node) is not nil, Walk is invoked recursively with visitor
// w for each of the non-nil children of node, followed by a call of
// w.Visit(nil).
//
func GoAST_Walk(v Visitor, node Node) {
    if v = v.Visit(node); v == nil {
        return
    }

    // walk children
    // (the order of the cases matches the order
    // of the corresponding node types in ast.go)
    switch n := node.(type) {
    // Comments and fields
    case *Comment:
        // do nothing
    case *DummyNode:
        for _, c := range n.List {
            GoAST_Walk(v, c)
        }

    case *CommentGroup:
        for _, c := range n.List {
            GoAST_Walk(v, c)
        }

    case *Field:
        if n.Doc != nil {
            GoAST_Walk(v, n.Doc)
        }
        walkIdentList(v, n.Names)
        GoAST_Walk(v, NewDummyNode(
            "Type",
            n.Type.Pos(),
            n.Type.End(),
            []Node{n.Type},
        ),
        )
        if n.Tag != nil {
            GoAST_Walk(v, NewDummyNode(
                "Tag",
                n.Tag.Pos(),
                n.Tag.End(),
                []Node{n.Tag},
            ),
            )
        }
        if n.Comment != nil {
            GoAST_Walk(v, n.Comment)
        }

    case *FieldList:
        for _, f := range n.List {
            GoAST_Walk(v, f)
        }

    // Expressions
    case *BadExpr, *Ident, *BasicLit:
        // nothing to do

    case *Ellipsis:
        if n.Elt != nil {
            GoAST_Walk(v, n.Elt)
        }

    case *FuncLit:
        GoAST_Walk(v, n.Type)
        GoAST_Walk(v, n.Body)

    case *CompositeLit:
        if n.Type != nil {
            GoAST_Walk(v, n.Type)
        }
        walkExprList(v, n.Elts)

    case *ParenExpr:
        GoAST_Walk(v, n.X)

    case *SelectorExpr:
        GoAST_Walk(v, n.X)
        GoAST_Walk(v, n.Sel)

    case *IndexExpr:
        GoAST_Walk(v, n.X)
        GoAST_Walk(v, n.Index)

    case *SliceExpr:
        GoAST_Walk(v, n.X)
        if n.Low != nil {
            GoAST_Walk(v,
                NewDummyNode(
                    "Low",
                    n.Low.Pos(),
                    n.Low.End(),
                    []Node{n.Low},
                ),
            )
        }
        if n.High != nil {
            GoAST_Walk(v,
                NewDummyNode(
                    "High",
                    n.High.Pos(),
                    n.High.End(),
                    []Node{n.High},
                ),
            )
        }

    case *TypeAssertExpr:
        GoAST_Walk(v, n.X)
        if n.Type != nil {
            GoAST_Walk(v, n.Type)
        }

    case *CallExpr:
        GoAST_Walk(v, n.Fun)
        walkExprList(v, n.Args)

    case *StarExpr:
        GoAST_Walk(v, n.X)

    case *UnaryExpr:
        GoAST_Walk(v, n.X)

    case *BinaryExpr:
        GoAST_Walk(v, n.X)
        GoAST_Walk(v, n.Y)

    case *KeyValueExpr:
        GoAST_Walk(v, n.Key)
        GoAST_Walk(v, n.Value)

    // Types
    case *ArrayType:
        if n.Len != nil {
            GoAST_Walk(v, n.Len)
        }
        GoAST_Walk(v, NewDummyNode(
            "ElemType",
            n.Elt.Pos(),
            n.Elt.End(),
            []Node{n.Elt},
        ),
        )

    case *StructType:
        GoAST_Walk(v, n.Fields)

    case *FuncType:
        if n.Params.NumFields() > 0 {
            list := n.Params.List
            GoAST_Walk(v,
                NewDummyNode("Params",
                    list[0].Pos(),
                    list[len(list)-1].End(),
                    func() []Node {
                        nodes := make([]Node, 0, len(list))
                        for _, c := range list {
                            nodes = append(nodes, Node(c))
                        }
                        return nodes
                    }()))
        }
        if n.Results != nil && n.Results.NumFields() > 0 {
            list := n.Results.List
            GoAST_Walk(v,
                NewDummyNode("Results",
                    list[0].Pos(),
                    list[len(list)-1].End(),
                    func() []Node {
                        nodes := make([]Node, 0, len(list))
                        for _, c := range list {
                            nodes = append(nodes, Node(c))
                        }
                        return nodes
                    }()))
        }

    case *InterfaceType:
        GoAST_Walk(v, n.Methods)

    case *MapType:
        GoAST_Walk(v, n.Key)
        GoAST_Walk(v, n.Value)

    case *ChanType:
        GoAST_Walk(v, NewDummyNode(
            "Direction",
            n.Pos(),
            n.End(),
            []Node{NewDummyNode(fmt.Sprint(n.Dir), n.Pos(), n.End(), nil)},
        ),
        )
        GoAST_Walk(v, NewDummyNode(
            "ElemType",
            n.Value.Pos(),
            n.Value.End(),
            []Node{n.Value},
        ),
        )

    // Statements
    case *BadStmt:
        // nothing to do

    case *DeclStmt:
        GoAST_Walk(v, n.Decl)

    case *EmptyStmt:
        // nothing to do

    case *LabeledStmt:
        GoAST_Walk(v, n.Label)
        GoAST_Walk(v, n.Stmt)

    case *ExprStmt:
        GoAST_Walk(v, n.X)

    case *SendStmt:
        GoAST_Walk(v, n.Chan)
        GoAST_Walk(v, n.Value)

    case *IncDecStmt:
        GoAST_Walk(v, n.X)

    case *AssignStmt:
        walkExprList(v, n.Lhs)
        walkExprList(v, n.Rhs)

    case *GoStmt:
        GoAST_Walk(v, n.Call)

    case *DeferStmt:
        GoAST_Walk(v, n.Call)

    case *ReturnStmt:
        walkExprList(v, n.Results)

    case *BranchStmt:
        if n.Label != nil {
            GoAST_Walk(v, n.Label)
        }

    case *BlockStmt:
        for _, c := range n.List {
            GoAST_Walk(v, c)
        }

    case *IfStmt:
        if n.Init != nil {
            GoAST_Walk(v,
                NewDummyNode(
                    "Init",
                    n.Init.Pos(),
                    n.Init.End(),
                    []Node{n.Init},
                ),
            )
        }
        GoAST_Walk(v,
            NewDummyNode(
                "Cond",
                n.Cond.Pos(),
                n.Cond.End(),
                []Node{n.Cond},
            ),
        )
        GoAST_Walk(v,
            NewDummyNode(
                "Body",
                n.Body.Pos(),
                n.Body.End(),
                []Node{n.Body},
            ),
        )
        if n.Else != nil {
            GoAST_Walk(v,
                NewDummyNode(
                    "Else",
                    n.Else.Pos(),
                    n.Else.End(),
                    []Node{n.Else},
                ),
            )
        }

    case *CaseClause:
        if n.List == nil {
            GoAST_Walk(v,
                NewDummyNode(
                    "default",
                    n.Case,
                    n.Colon,
                    nil,
                ),
            )
        }
        walkExprList(v, n.List)
        walkStmtList(v, n.Body)

    case *SwitchStmt:
        if n.Init != nil {
            GoAST_Walk(v,
                NewDummyNode(
                    "Init",
                    n.Init.Pos(),
                    n.Init.End(),
                    []Node{n.Init},
                ),
            )
        }
        if n.Tag != nil {
            GoAST_Walk(v,
                NewDummyNode(
                    "Tag",
                    n.Tag.Pos(),
                    n.Tag.End(),
                    []Node{n.Tag},
                ),
            )
        }
        GoAST_Walk(v, n.Body)

    case *TypeSwitchStmt:
        if n.Init != nil {
            GoAST_Walk(v,
                NewDummyNode(
                    "Init",
                    n.Init.Pos(),
                    n.Init.End(),
                    []Node{n.Init},
                ),
            )
            GoAST_Walk(v, n.Init)
        }
        GoAST_Walk(v, n.Assign)
        GoAST_Walk(v, n.Body)

    case *CommClause:
        if n.Comm != nil {
            GoAST_Walk(v, n.Comm)
        } else {
            GoAST_Walk(v, NewDummyNode(
                "default",
                n.Case,
                n.Colon,
                nil,
            ),
            )
        }
        walkStmtList(v, n.Body)

    case *SelectStmt:
        GoAST_Walk(v, n.Body)

    case *ForStmt:
        if n.Init != nil {
            GoAST_Walk(v, NewDummyNode(
                "Init",
                n.Init.Pos(),
                n.Init.End(),
                []Node{n.Init},
            ),
            )
        }
        if n.Cond != nil {
            GoAST_Walk(v, NewDummyNode(
                "Cond",
                n.Cond.Pos(),
                n.Cond.End(),
                []Node{n.Cond},
            ),
            )
        }
        if n.Post != nil {
            GoAST_Walk(v, NewDummyNode(
                "Post",
                n.Post.Pos(),
                n.Post.End(),
                []Node{n.Post},
            ),
            )
        }
        GoAST_Walk(v, NewDummyNode(
            "Body",
            n.Body.Pos(),
            n.Body.End(),
            []Node{n.Body},
        ),
        )

    case *RangeStmt:
        GoAST_Walk(v, NewDummyNode(
            "Key",
            n.Key.Pos(),
            n.Key.End(),
            []Node{n.Key},
        ),
        )
        if n.Value != nil {
            GoAST_Walk(v, NewDummyNode(
                "Value",
                n.Value.Pos(),
                n.Value.End(),
                []Node{n.Value},
            ),
            )
        }
        GoAST_Walk(v, NewDummyNode(
            fmt.Sprint(n.Tok),
            n.TokPos,
            n.TokPos,
            nil,
        ),
        )
        GoAST_Walk(v, NewDummyNode(
            "Rangeable",
            n.X.Pos(),
            n.X.End(),
            []Node{n.X},
        ),
        )
        GoAST_Walk(v, NewDummyNode(
            "Body",
            n.Body.Pos(),
            n.Body.End(),
            []Node{n.Body},
        ),
        )

    // Declarations
    case *ImportSpec:
        if n.Doc != nil {
            GoAST_Walk(v, n.Doc)
        }
        if n.Name != nil {
            GoAST_Walk(v, n.Name)
        }
        GoAST_Walk(v, n.Path)
        if n.Comment != nil {
            GoAST_Walk(v, n.Comment)
        }

    case *ValueSpec:
        walkIdentList(v, n.Names)
        walkExprList(v, n.Values)
        if n.Doc != nil {
            GoAST_Walk(v, n.Doc)
        }
        if n.Type != nil {
            GoAST_Walk(v, n.Type)
        }
        if n.Comment != nil {
            GoAST_Walk(v, n.Comment)
        }

    case *TypeSpec:
        if n.Doc != nil {
            GoAST_Walk(v, n.Doc)
        }
        GoAST_Walk(v, n.Name)
        GoAST_Walk(v, n.Type)
        if n.Comment != nil {
            GoAST_Walk(v, n.Comment)
        }

    case *BadDecl:
        // nothing to do

    case *GenDecl:
        if n.Doc != nil {
            GoAST_Walk(v, n.Doc)
        }
        for _, s := range n.Specs {
            GoAST_Walk(v, s)
        }

    case *FuncDecl:
        if n.Doc != nil {
            GoAST_Walk(v, n.Doc)
        }
        if n.Recv != nil {
            GoAST_Walk(v, NewDummyNode(
                "Recv",
                n.Recv.Pos(),
                n.Recv.End(),
                []Node{n.Recv},
            ),
            )
        }
        GoAST_Walk(v, NewDummyNode(
            "Name",
            n.Name.Pos(),
            n.Name.End(),
            []Node{n.Name},
        ),
        )
        GoAST_Walk(v, NewDummyNode(
            "Type",
            n.Type.Pos(),
            n.Type.End(),
            []Node{n.Type},
        ),
        )
        if n.Body != nil {
            GoAST_Walk(v, NewDummyNode(
                "Body",
                n.Body.Pos(),
                n.Body.End(),
                []Node{n.Body},
            ),
            )
        }

    // Files and packages
    case *File:
        if n.Doc != nil {
            GoAST_Walk(v, n.Doc)
        }
        GoAST_Walk(v, n.Name)
        walkDeclList(v, n.Decls)
        if len(n.Imports) > 0 {
            list := n.Imports
            GoAST_Walk(v,
                NewDummyNode("Imports",
                    list[0].Pos(),
                    list[len(list)-1].End(),
                    func() []Node {
                        nodes := make([]Node, 0, len(list))
                        for _, c := range list {
                            nodes = append(nodes, Node(c))
                        }
                        return nodes
                    }()))
        }
        if len(n.Unresolved) > 0 {
            list := n.Unresolved
            GoAST_Walk(v,
                NewDummyNode("Unresolved",
                    list[0].Pos(),
                    list[len(list)-1].End(),
                    func() []Node {
                        nodes := make([]Node, 0, len(list))
                        for _, c := range list {
                            nodes = append(nodes, Node(c))
                        }
                        return nodes
                    }()))
        }
        for _, g := range n.Comments {
            GoAST_Walk(v, g)
        }
        // don't walk n.Comments - they have been
        // visited already through the individual
        // nodes

    case *Package:
        for _, f := range n.Files {
            GoAST_Walk(v, f)
        }

    default:
        fmt.Printf("ast.Walk: unexpected node type %T", n)
        panic("ast.Walk")
    }

    v.Visit(nil)
}
