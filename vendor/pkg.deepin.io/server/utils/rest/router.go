package rest

import (
	"github.com/gin-gonic/gin"
)

// Router control
type Router struct {
	engine     *gin.Engine
	processors []Processor
}

// NewRouter create new router with engine
func NewRouter(engine *gin.Engine) *Router {
	return &Router{
		engine: engine,
	}
}

var _staticRouter *Router

// StaticRouter is the default router
func StaticRouter() *Router {
	// gin.SetMode(gin.ReleaseMode)
	if nil == _staticRouter {
		eng := gin.Default()
		_staticRouter = NewRouter(eng)
	}
	return _staticRouter
}

// Plugin insert a Processor to router
func (r *Router) Plugin(processor Processor) error {
	if nil != processor {
		r.processors = append(r.processors, processor)
		return nil
	}
	panic("Processor must not be nil")
}

// BootStrap start up router to work
func (r *Router) BootStrap(bind string) error {
	for _, p := range r.processors {
		err := p.Run(r.engine)
		if nil != err {
			panic(err)
		}
	}
	r.engine.Run(bind)
	return nil
}
