package framework

type IGroup interface {
	Get(string, ...ControllerHandler)
	Post(string, ...ControllerHandler)
	Put(string, ...ControllerHandler)
	Delete(string, ...ControllerHandler)

	Use(...ControllerHandler)

	Group(string) IGroup
}

type Group struct {
	core *Core
	parent *Group
	prefix string
	middlewares []ControllerHandler  // 用于一组使用
}

func NewGroup(c *Core, prefix string) IGroup {
	return &Group{
		core: c,
		prefix: prefix,
	}
}
func (g *Group) Use(middlewares ...ControllerHandler) {
	g.middlewares = append(g.middlewares, middlewares...)
}

func (g *Group) GetAbsUri() string {
	if g.parent != nil {
		return g.parent.GetAbsUri() + g.prefix
	}
	return g.prefix
}

func (g *Group) GetParentMiddlewares() []ControllerHandler {
	if g.parent != nil {
		return append(g.parent.GetParentMiddlewares(), g.middlewares...)
	}
	return g.middlewares
}
func (g *Group) Get(uri string, handlers ...ControllerHandler) {
	absUri := g.GetAbsUri() + uri
	allHandlers := append(g.GetParentMiddlewares(), handlers...)
	g.core.Get(absUri, allHandlers...)
}

func (g *Group) Post(uri string, handlers ...ControllerHandler) {
	absUri := g.GetAbsUri() + uri
	allHandlers := append(g.GetParentMiddlewares(), handlers...)
	g.core.Post(absUri, allHandlers...)
}

func (g *Group) Put(uri string, handlers ...ControllerHandler) {
	absUri := g.GetAbsUri() + uri
	allHandlers := append(g.GetParentMiddlewares(), handlers...)
	g.core.Put(absUri, allHandlers...)
}

func (g *Group) Delete(uri string, handlers ...ControllerHandler) {
	absUri := g.GetAbsUri() + uri
	allHandlers := append(g.GetParentMiddlewares(), handlers...)
	g.core.Delete(absUri, allHandlers...)
}

func (g *Group) Group(uri string) IGroup {
	return &Group{
		core: g.core,
		parent: g,
		prefix: uri,
	}
}
