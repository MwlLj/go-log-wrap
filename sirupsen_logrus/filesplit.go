package logwrap

/*
** 实现 io.Write 接口
*/
type Output struct {
    name string
}

func (self *Output) Write(data []byte) (n int, err error) {
    return 0, nil
}

func (self *Output) Close() {
}

func NewOutput(name string) *Output {
    o := Output{
        name: name,
    }
    return &o
}
