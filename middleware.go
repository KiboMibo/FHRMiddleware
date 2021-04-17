package FHRMiddleware

import (
	"github.com/lab259/cors"
	"github.com/valyala/fasthttp"
)

// Like Alice middleware https://github.com/justinas/alice

type Middleware func(next fasthttp.RequestHandler) fasthttp.RequestHandler

type Chain struct {
	handlers []Middleware
}

var (
	logger Logger
	crs    *cors.Cors
)

func SetLogger(l Logger) {
	logger = l
}

func SetCors(opts cors.Options) {
	crs = cors.New(opts)
}

func NewChain(handlers ...Middleware) *Chain {
	if logger == nil {
		logger = DefaultLogger{}
	}

	if crs == nil {
		crs = cors.Default()
	}

	return &Chain{
		handlers: append([]Middleware{RecoveryHandle, RequestLogger}, handlers...),
	}
}

func (c *Chain) Add(m ...Middleware) *Chain {
	c.handlers = append(c.handlers, m...)
	return c
}

func (c *Chain) Apply(h fasthttp.RequestHandler) (r fasthttp.RequestHandler) {
	for i := range c.handlers {
		h = c.handlers[len(c.handlers)-1-i](crs.Handler(h))
	}

	return h
}
