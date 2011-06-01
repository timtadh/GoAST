package Chan

func main() {
    var x chan int
    var y chan<-int
    var z <-chan int

    x <- y
    q := <-x
    <-z
}
