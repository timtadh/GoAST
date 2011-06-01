package funcdecl

func forward_func()
func (self *obj) forward_method()

func simplefunc() { body }
func (self *obj) simplemeth() { body }


func func_args(x,y,z int) (int, string, int, int) { body }
func (self *obj) meth_args(x,y,z int) (int, string, int, int) { body }

func littest() {
    func() {}
}
