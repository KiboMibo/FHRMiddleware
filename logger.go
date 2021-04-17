package FHRMiddleware

import (
	"log"

	"github.com/valyala/fasthttp"
)

// Interface for compability with other loggers
type Logger interface {
	Info(args ...interface{})
	Infof(template string, args ...interface{})
}

// Package for standart logger
type DefaultLogger struct{}

func (d DefaultLogger) Info(args ...interface{}) {
	log.Println(args...)
}

func (d DefaultLogger) Infof(template string, args ...interface{}) {
	log.Printf(template, args...)
}

// RequestLogger writes logs for requests
func RequestLogger(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		next(ctx)
		logger.Infof("%s %d %s %q from: %s", ctx.Method(), ctx.Response.StatusCode(), fasthttp.StatusMessage(ctx.Response.StatusCode()), ctx.RequestURI(), ctx.RemoteAddr())
	}
}
