package call

func main() {
    q()
    r(q, func(){ return 5 })
    s(r, v...)
    return q, r, s()
}
