//go:build try

//go:generate go install github.com/goghcrow/go-try/cmd/trygen@main
//go:generate trygen
package example

import (
	"io"
	"os"

	. "github.com/goghcrow/go-try"
)

//goland:noinspection GoUnhandledErrorResult
func CopyFile(src, dst string) (err error) {
	r := Try(os.Open(src))
	defer r.Close()

	w := Try(os.Create(dst))
	defer func() {
		w.Close()
		if err != nil {
			os.Remove(dst) // only if a “try” fails
		}
	}()

	Try(io.Copy(w, r))
	Try0(w.Close())
	return nil
}
