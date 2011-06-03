package tree

import "testing"

func TestHello(t *testing.T) {
    t.Log("Hello")
}

func TestPrintTree(t *testing.T) {
    t.Log(NewNode("root").
        AddKid(NewNode("Hiya")).
        AddKid(
        NewNode("yippy").
            AddKid(NewNode("Level3")),
    ),
    )
}
