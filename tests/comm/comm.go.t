package comm

func main() {
    select {
    case q <- x:
        print("q <- x")
    case y := <-r:
        print("y := <-r")
    case q:
    default:
        print("default")
    }
}
