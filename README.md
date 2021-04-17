# FHRMiddleware

Simple chain middleware for fasthttp and fasthttp Router

```bigquery
go get github.com/KiboMibo/FHRMiddleware
```

## Fast example
```go

import middleware "github.com/KiboMibo/FHRMiddleware"

func main() {
	//Init middleware
	md := middleware.NewChain()
    
	//Apply to handler
	fasthttp.ListenAndServe(":8080", md.Apply(FooHandle))
}

```

more [examples](https://github.com/KiboMibo/FHRMiddleware/examples)

## Defaults

By default includes logger and recovery middlewares

CORS is enabled by default but with default options.

## Customize

Set logger at your own.

With [ZAP](https://go.uber.org/zap) example
```go

import (
	middleware "github.com/KiboMibo/FHRMiddleware"
	"go.uber.org/zap"
)

func main() {
    md := middleware.NewChain()
    middleware.SetLogger(zap.NewExample().Sugar())
}

```

Customize CORS policy:
```go
import (
    middleware "github.com/KiboMibo/FHRMiddleware"
    "github.com/lab259/cors"
)

func main() {
    md := middleware.NewChain()
    middleware.SetCors(cors.Options{
        AllowedOrigins: "foo.bar",
        AllowedMethods: []string{"GET"},
        AllowedOrigins: []string{"localhost", "127.0.0.1"},
    })
}
```
See more at https://pkg.go.dev/github.com/lab259/cors