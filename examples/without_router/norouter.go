package main

import (
	"encoding/json"
	"log"

	middleware "github.com/KiboMibo/FHRMiddleware"

	"github.com/valyala/fasthttp"
)

func main() {
	// Initialize middlewares
	md := middleware.NewChain(SetHeader)

	address := "localhost:8080"
	log.Println("Listen:", address)
	// Apply to handler
	if err := fasthttp.ListenAndServe(address, md.Apply(TestHandle)); err != nil {
		panic(err)
	}
}

func TestHandle(ctx *fasthttp.RequestCtx) {
	r, err := json.Marshal(map[string]interface{}{"data": "Hello World", "X-Test-Header": ctx.Response.Header.Peek("X-Test-Header")})
	if err != nil {
		panic(err)
	}

	// Test panic handler
	// uncomment below
	// foo := make([]int, 1)
	// fmt.Println(foo(5))

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
