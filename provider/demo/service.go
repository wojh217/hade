package demo

import "github.com/wojh217/hade/framework"

type DemoService struct {
	Service

	container framework.Container
}

func (s *DemoService) GetFoo() Foo {
	return Foo{
		Name: "I am foo",
	}
}

func NewDemoService(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	return &DemoService{container: container}, nil
}