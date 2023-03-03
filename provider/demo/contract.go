package demo

const Key = "hade:demo"

// Service 一个真正的服务
type Service interface {
	GetFoo() Foo
}

type Foo struct {
	Name string
}


