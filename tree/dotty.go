package tree

import "fmt"
import "strings"

func (self *Node) Dotty() string {
    header := "digraph AST {\n"
    node := "  %v[label=\"%v\", shape=\"rect\"];"
    edge := "  %v -> %v;"
    footer := "\n}\n"

    nodes := make([]string, 0, 10)
    edges := make([]string, 0, 10)

    var dotnode func(*Node, string)
    dotnode = func(n *Node, parent string) {
        name := fmt.Sprintf("node%d", len(nodes))
        nodes = append(nodes,
            fmt.Sprintf(node, name, strings.Replace(n.Label, "\"", "\\\"", -1)))
        if parent != "" {
            edges = append(edges, fmt.Sprintf(edge, parent, name))
        }
        for _, child := range n.Children {
            dotnode(child, name)
        }
    }

    dotnode(self, "")
    snodes := strings.Join(nodes, "\n")
    sedges := strings.Join(edges, "\n")
    return strings.Join([]string{header, snodes, sedges, footer}, "\n")
}
