package framework

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Core struct {
	router map[string]*Tree

	middlewares []ControllerHandler  // 用于一组使用
}

func NewCore() *Core {
	router := map[string]*Tree{}

	router["GET"] = NewTree()
	router["POST"] = NewTree()
	router["PUT"] = NewTree()
	router["DELETE"] = NewTree()
	return &Core{router: router}
}
func disPlayNode(n *node) {
	fmt.Printf("node: %p, segment: %s, handlers: %p\n", n, n.segment, n.handlers)
	for _, cnode := range n.childs {
		disPlayNode(cnode)
	}
}

func (c *Core) DisplayTree() {
	root := c.router["GET"].root
	disPlayNode(root)

}

// 对应 Method = Get
func (c *Core) Get(url string, handlers ...ControllerHandler) {
	allHandlers := make([]ControllerHandler, len(c.middlewares) + len(handlers))
	copy(allHandlers, append(c.middlewares, handlers...))


	if err := c.router["GET"].AddRouter(url, allHandlers); err != nil {
		log.Fatal("add router error: ", err)
	}
}

// 对应 Method = POST
func (c *Core) Post(url string, handlers ...ControllerHandler) {
	allHandlers := append(c.middlewares, handlers...)
	if err := c.router["POST"].AddRouter(url, allHandlers); err != nil {
		log.Fatal("add router error: ", err)
	}
}

// 对应 Method = PUT
func (c *Core) Put(url string, handlers ...ControllerHandler) {
	allHandlers := append(c.middlewares, handlers...)
	if err := c.router["PUT"].AddRouter(url, allHandlers); err != nil {
		log.Fatal("add router error: ", err)
	}
}

// 对应 Method = DELETE
func (c *Core) Delete(url string, handlers ...ControllerHandler) {
	allHandlers := append(c.middlewares, handlers...)
	if err := c.router["DELETE"].AddRouter(url, allHandlers); err != nil {
		log.Fatal("add router error: ", err)
	}
}

func (c *Core) Group(prefix string) IGroup {
	return NewGroup(c, prefix)
}

// 匹配路由，如果没有匹配到，返回nil
func (c *Core) FindRouteByRequest(request *http.Request) []ControllerHandler {
	// uri 和 method 全部转换为大写，保证大小写不敏感
	uri := request.URL.Path
	method := request.Method
	upperMethod := strings.ToUpper(method)

	// 查找第一层map
	if methodHandlers, ok := c.router[upperMethod]; ok {
		return methodHandlers.FindHandlers(uri)
	}
	return nil
}

// ServeHTTP 会多线程并发访问, 要考虑并发处理是否会产生竞争
func (c *Core) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := NewContext(r, w)
	handlers := c.FindRouteByRequest(r)   //
	if handlers == nil || len(handlers) == 0 {
		ctx.Json(404, "not found")
		return
	}

	// 之前handler只写到core中，现在因为中间件只接收ctx参数，因此也要写入ctx中
	//all_handlers := make([]ControllerHandler, len(c.middlewares) + len(handlers))
	//copy(all_handlers, c.middlewares)
	//copy(all_handlers[len(c.middlewares):], handlers)
	//ctx.SetHandlers(all_handlers)

	fmt.Printf("handlers: %v\n", handlers)
	ctx.SetHandlers(handlers)

	// 这里使用ctx.Next()函数，内部按顺序逐个调用handler链表
	if err := ctx.Next(); err != nil {
		ctx.Json(500, "inner error")
		return
	}
}


func (c *Core) Use(middlewares ...ControllerHandler) {
	c.middlewares = append(c.middlewares, middlewares...)
}

func (c *Core) GetMiddleWares() []ControllerHandler {
	return c.middlewares
}