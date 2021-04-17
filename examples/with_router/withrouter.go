package main

import (
	"encoding/json"
	"log"

	middleware "github.com/KiboMibo/FHRMiddleware"

	"github.com/fasthttp/router"
	"github.com/lab259/cors"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

func main() {
	r := router.New()

	// Add custom middleware
	// to split middlewares just make another chain
	// default = middleware.NewChain()
	// setHeader = middleware.NewChain(SetHeader)
	// and apply them to routes
	md := middleware.NewChain(SetHeader)

	// Can use custom logger
	// by default used standart log
	middleware.SetLogger(zap.NewExample().Sugar())

	// By default CORS is in allowAll mode
	// https://pkg.go.dev/github.com/lab259/cors documentation for cors configuration
	middleware.SetCors(cors.Options{
		AllowedMethods: []string{"GET"},
		AllowedOrigins: []string{"localhost", "127.0.0.1"},
	})

	r.GET("/test", md.Apply(TestHandle))

	// make MethodNotAllowed and NotFound logging
	r.MethodNotAllowed = md.Apply(NotAllow)
	r.NotFound = md.Apply(NotFound)

	address := "localhost:8080"
	log.Println("Listen:", address)
	if err := fasthttp.ListenAndServe(":8080", r.Handler); err != nil {
		panic(err)
	}
}

func TestHandle(ctx *fasthttp.RequestCtx) {
	r, err := json.Marshal(map[string]interface{}{
		"data":          "Hello World",
		"X-Test-Header": ctx.Response.Header.Peek("X-Test-Header"),
	})
	if err != nil {
		panic(err)
	}

	ctx.SetContentType("application/json")
	ctx.Write(r)
}

// SetHeader is an example of custom middleware
func SetHeader(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("X-Test-Header", "Blah-Blah")
		next(ctx)
	}
}

func NotAllow(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
	ctx.Write(nil)
}

func NotFound(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusNotFound)
	ctx.Write(nil)
}
