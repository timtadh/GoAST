package defer_


func main() {
    defer func() { print("hello") }()
    X
}
