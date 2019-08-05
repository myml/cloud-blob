package rest

import "github.com/gin-gonic/gin"

// Processor handle api end point of router
type Processor interface {
	Run(*gin.Engine) error
}

// ProcessorHandler is handler of Processor
type ProcessorHandler func(*gin.Engine) error

// DefaultProcessor is simple processor implementation
type DefaultProcessor struct {
	Handler ProcessorHandler
}

// Run process engine with handler
func (p *DefaultProcessor) Run(engine *gin.Engine) error {
	return p.Handler(engine)
}
