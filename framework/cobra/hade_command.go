package cobra

import "github.com/wojh217/hade/framework"

func (c *Command) SetContainer(container framework.Container) {
	c.container = container
}

func (c *Command) GetContainer() framework.Container{
	return c.Root().container
}