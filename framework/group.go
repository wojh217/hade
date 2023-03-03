package framework

type IGroup interface {
	Get(string, ControllerHandler)
	Post(string, ControllerHandler)
	Put(string, ControllerHandler)
	Delete(string, ControllerHandler)

	Group(string) IGroup
}

type Group struct {
	core *Core
	parent *Group
	prefix string
}

func NewGroup(c *Core, prefix string) IGroup {
	return &Group{
		core: c,
		prefix: prefix,
	}
}

func (g *Group) GetAbsUri() string {
	if g.parent != nil {
		return g.parent.GetAbsUri() + g.prefix
	}
	return g.prefix
}
func (g *Group) Get(uri string, handler ControllerHandler) {
	absUri := g.GetAbsUri() + uri
	g.core.Get(absUri, handler)
}

func (g *Group) Post(uri string, handler ControllerHandler) {
	absUri := g.GetAbsUri() + uri
	g.core.Post(absUri, handler)
}

func (g *Group) Put(uri string, handler ControllerHandler) {
	absUri := g.GetAbsUri() + uri
	g.core.Put(absUri, handler)
}

func (g *Group) Delete(uri string, handler ControllerHandler) {
	absUri := g.GetAbsUri() + uri
	g.core.Delete(absUri, handler)
}

func (g *Group) Group(uri string) IGroup {
	return &Group{
		core: g.core,
		parent: g,
		prefix: uri,
	}
}
