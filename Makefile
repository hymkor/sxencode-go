all:
	go fmt ./...
	go build

test-gmnlisp:
	go run example.go | gmnlisp test.lsp

test-oki:
	go run example.go > sample.log
	echo "(load \"test-oki.lsp\")" | islisp

test-sbcl:
	go run example.go | sbcl --load "test-sbcl.lsp"

test-marshal:
	go run example-marshal.go | gmnlisp test.lsp

example:
	go run example.go
