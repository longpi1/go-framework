package gpn

import (
	_ "fmt"
	"log"
	_ "log"
	"net/http"
)

// Engine implement the interface of ServeHTTP
type (
	HandlerFunc func(*Context)

	RouterGroup struct {
		prefix   string
		meddlers []HandlerFunc
		parent   *RouterGroup
		engine   *Engine
	}

	Engine struct {
		*RouterGroup
		router *router
		groups []*RouterGroup // 存储组
	}
)

// New is the constructor of gee.Engine
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}


//创建新的RouterGroup
//所有组共享相同的Engine实例
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.engine.addRoute("GET", pattern, handler)
}


func (group *RouterGroup)  POST(pattern string, handler HandlerFunc) {
	group.engine.addRoute("POST", pattern, handler)
}

func (group *RouterGroup)  Run(addr string) (err error) {
	return http.ListenAndServe(addr, group.engine)
}

func (group *RouterGroup) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	group.engine.router.handle(c)
}