package framework

import (
	"net/http"
	"strings"
)

type Core struct {
	router map[string]map[string]ControllerHandler
}

func NewCore() *Core {
	getRouter := map[string]ControllerHandler{}
	postRouter := map[string]ControllerHandler{}
	putRouter := map[string]ControllerHandler{}
	deleteRouter := map[string]ControllerHandler{}

	router := map[string]map[string]ControllerHandler{}
	router["GET"] = getRouter
	router["POST"] = postRouter
	router["PUT"] = putRouter
	router["DELETE"] = deleteRouter
	return &Core{router: router}
}


// 对应 Method = Get
func (c *Core) Get(url string, handler ControllerHandler) {
	upperUrl := strings.ToUpper(url)
	c.router["GET"][upperUrl] = handler
}

// 对应 Method = POST
func (c *Core) Post(url string, handler ControllerHandler) {
	upperUrl := strings.ToUpper(url)
	c.router["POST"][upperUrl] = handler
}

// 对应 Method = PUT
func (c *Core) Put(url string, handler ControllerHandler) {
	upperUrl := strings.ToUpper(url)
	c.router["PUT"][upperUrl] = handler
}

// 对应 Method = DELETE
func (c *Core) Delete(url string, handler ControllerHandler) {
	upperUrl := strings.ToUpper(url)
	c.router["DELETE"][upperUrl] = handler
}


// 匹配路由，如果没有匹配到，返回nil
func (c *Core) FindRouteByRequest(request *http.Request) ControllerHandler {
	// uri 和 method 全部转换为大写，保证大小写不敏感
	uri := request.URL.Path
	method := request.Method
	upperMethod := strings.ToUpper(method)
	upperUri := strings.ToUpper(uri)

	// 查找第一层map
	if methodHandlers, ok := c.router[upperMethod]; ok {
		// 查找第二层map
		if handler, ok := methodHandlers[upperUri]; ok {
			return handler
		}
	}
	return nil
}

func (c *Core) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := NewContext(r, w)
	router := c.FindRouteByRequest(r)
	if router == nil {
		ctx.Json(404, "not found")
		return
	}

	if err := router(ctx); err != nil {
		ctx.Json(500, "inner error")
		return
	}
}
