package FHRMiddleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"runtime"
	"runtime/debug"
	"text/template"

	"github.com/valyala/fasthttp"
)

var tpl *template.Template

// Took this from Negroni recovery (https://github.com/urfave/negroni/blob/master/recovery.go)

// PanicInformation contains all
// elements for printing stack informations.
type PanicInformation struct {
	RecoveredPanic interface{}
	Stack          []byte
	Request        *fasthttp.RequestCtx
}

// StackAsString returns a printable version of the stack
func (p *PanicInformation) StackAsString() string {
	return string(p.Stack)
}

// RequestDescription returns a printable description of the url
func (p *PanicInformation) RequestDescription() string {

	if p.Request == nil {
		return "Request is nil"
	}

	var queryOutput string

	if p.Request.QueryArgs().Len() > 0 {
		queryOutput = "?" + p.Request.QueryArgs().String()
	}
	return fmt.Sprintf("%s %s%s", p.Request.Method(), p.Request.Path(), queryOutput)
}

// RecoveryHandle handles panics in chain
func RecoveryHandle(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		defer func() {
			if err := recover(); err != nil {
				logger.Info(err)
				ctx.SetStatusCode(fasthttp.StatusInternalServerError)

				stack := make([]byte, 1024*8)
				stack = stack[:runtime.Stack(stack, false)]
				infos := &PanicInformation{RecoveredPanic: err, Request: ctx}
				infos.Stack = stack
				logger.Infof("provided PanicHandlerFunc panic'd: %s, trace:\n%s\n", err, debug.Stack())

				body, err := PrepareTemplate(ctx, infos)
				if err != nil {
					logger.Info(err)
				}

				ctx.Write(body)
			}
		}()
		next(ctx)
	}
}

// PrepareTemplate render panic template
func PrepareTemplate(ctx *fasthttp.RequestCtx, infos *PanicInformation) ([]byte, error) {
	if tpl == nil {
		tpl = template.Must(template.New("PanicPage").Parse(panicHTML))
	}

	buf := &bytes.Buffer{}
	if err := tpl.Execute(buf, infos); err != nil {
		j, err := json.Marshal(infos)
		if err != nil {
			return nil, err
		}
		ctx.SetContentType("application/json")
		return j, nil
	}
	ctx.SetContentType("text/html")
	return buf.Bytes(), nil
}

var panicHTML = `
<html>
<head>
    <title>DON'T PANIC</title></head>
<style type="text/css">
    html, body {
        font-family: Helvetica, Arial, Sans;
        color: #333333;
        background-color: #ffffff;
        margin: 0px;
    }
    h1 {
        color: #ffffff;
        background-color: #f14c4c;
        padding: 20px;
        border-bottom: 1px solid #2b3848;
    }
    .block {
        margin: 2em;
    }

    .panic-stack-raw pre {
        padding: 1em;
        background: #f6f8fa;
        border: dashed 1px;
    }
    .panic-interface-title {
        font-weight: bold;
    }
</style>
<body>
<h1>DON'T PANIC</h1>

<div class="panic-interface block">
    <h3>{{ .RequestDescription }}</h3>
</div>

{{ if .Stack }}
<div class="panic-stack-raw block">
    <h3>Runtime Stack</h3>
    <pre>{{ .StackAsString }}</pre>
</div>
{{ end }}
</body>
</html>
`
