# 𝐖𝐡𝐚𝐭 𝐢𝐬 𝕘𝕠-𝕥𝕣𝕪

A src2src translator for propagating the error in golang.

# [WIP]Quick Start

Create source files ending with `_try.go` / `_try_test.go`.

Build tag `//go:build try` required.

Then `go generate -tags try ./...` (or run by IDE whatever).

And it is a good idea to switch custom build tag to `try` when working in goland or vscode,
so IDE will be happy to index and check your code.

```golang
//go:build try

//go:generate go install github.com/goghcrow/go-try/cmd/trygen
//go:generate trygen

package main

import (
	. "github.com/goghcrow/go-try"
)


```
