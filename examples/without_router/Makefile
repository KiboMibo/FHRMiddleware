withoutrouter: clean
	go get -u github.com/valyala/fasthttp
	go get -u github.com/fasthttp/router
	go get -u github.com/lab259/cors
	go get -u go.uber.org/zap
	go build -o withoutrouter norouter.go

clean:
	rm -f withoutrouter