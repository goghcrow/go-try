//go:build try

//go:generate go install github.com/goghcrow/go-try/cmd/trygen@main
//go:generate trygen

package example

import (
	"fmt"
	"io"
	"os"

	. "github.com/goghcrow/go-try"
)

// CopyFile example
// from https://github.com/golang/proposal/blob/master/design/32437-try-builtin.md#examples
//
//goland:noinspection GoUnhandledErrorResult
func CopyFile(src, dst string) (err error) {
	defer HandleErrorf(&err, "copy %s %s", src, dst)

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

func HandleErrorf(err *error, format string, args ...any) {
	if *err != nil {
		*err = fmt.Errorf(format+": %v", append(args, *err)...)
	}
}
