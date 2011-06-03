package tree

import "fmt"
import "strings"

type Node struct {
    Label    string
    Children []*Node
}

func NewNode(label string) *Node {
    return &Node{
        Label:    label,
        Children: make([]*Node, 0, 10),
    }
}

func (self *Node) AddKid(n *Node) *Node {
    self.Children = append(self.Children, n)
    return self
}

func (self *Node) String() string {
    if self == nil {
        return "<Node nil>"
    }
    var walk func(*Node) []string
    walk = func(n *Node) []string {
        l := make([]string, 0, 10)
        s := fmt.Sprintf("%d:%v", len(n.Children), n.Label)
        l = append(l, s)
        for _, c := range n.Children {
            l = append(l, walk(c)...)
        }
        return l
    }
    l := walk(self)
    return strings.Join(l, "\n")
}
