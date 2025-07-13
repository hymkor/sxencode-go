all:
	go fmt ./...
	go build

test-gmnlisp:
	go run example.go | gmnlisp test.lsp

test-oki:
	go run example.go > sample.log
	echo "(load \"test-oki.lsp\")" | islisp

example:
	go run example.go
