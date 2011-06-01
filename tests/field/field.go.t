package branch

type X struct {
    q int "my tag"
}

type anon func(int, string) (s int, q string, f float64)
type no_in func() (int)
type no_out func(int)
type none func()

func main(x, y, z int, q int, r int) {
    return
}
