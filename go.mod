module github.com/goghcrow/go-try

go 1.19

require (
	github.com/goghcrow/go-ansi v1.0.1
	github.com/goghcrow/go-loader v0.0.4-0.20240221113906-cab11067771f
	github.com/goghcrow/go-matcher v0.0.5-0.20240221112341-6675288f4167
	golang.org/x/tools v0.18.0
)

require golang.org/x/mod v0.15.0 // indirect

//replace github.com/goghcrow/go-matcher => ./../go-matcher
//replace github.com/goghcrow/go-loader => ./../go-loader
//replace github.com/goghcrow/go-ast-matcher => ./../go-ast-matcher
