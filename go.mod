module github.com/goghcrow/go-try

go 1.19

require (
	github.com/goghcrow/go-ast-matcher v0.0.13-0.20240130184623-e1946368edb4
	golang.org/x/tools v0.17.0
)

require golang.org/x/mod v0.14.0 // indirect

//replace github.com/goghcrow/go-ast-matcher => ./../go-ast-matcher
