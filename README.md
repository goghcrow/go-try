## 𝐖𝐡𝐚𝐭 𝐢𝐬 𝕘𝕠-𝕥𝕣𝕪

A source to source translator for error-propagating in golang.

The implementation of [***Proposal: A built-in Go error check function, try***](https://github.com/golang/proposal/blob/master/design/32437-try-builtin.md).


## Quick Start

Create source files ending with `_try.go` / `_try_test.go`.

Build tag `//go:build try` required.

Then `go generate -tags try ./...` (or run by IDE whatever).

And it is a good idea to switch custom build tag to `try` when working in goland or vscode,
so IDE will be happy to index and check your code.

## Example

Create file `copyfile_try.go`.

```golang
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
```

run go:generate and output `copyfile.go`.

```golang
//go:build !try

// Code generated by github.com/goghcrow/go-try DO NOT EDIT.
//go:generate go install github.com/goghcrow/go-try/cmd/trygen@main
//go:generate trygen
package example

import (
	. "github.com/goghcrow/go-try/rt"
	"io"
	"os"
)

//goland:noinspection GoUnhandledErrorResult
func CopyFile(src, dst string) (err error) {
	𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟮 := os.Open(src)
	if 𝗲𝗿𝗿𝟮 != nil {
		return 𝗲𝗿𝗿𝟮
	}
	r := 𝘃𝗮𝗹𝟭
	defer r.Close()
	𝘃𝗮𝗹𝟯, 𝗲𝗿𝗿𝟰 := os.Create(dst)
	if 𝗲𝗿𝗿𝟰 != nil {
		return 𝗲𝗿𝗿𝟰
	}
	w := 𝘃𝗮𝗹𝟯
	defer func() {
		w.Close()
		if err != nil {
			os.Remove(dst)
		}
	}()
	_, 𝗲𝗿𝗿𝟱 := io.Copy(w, r)
	if 𝗲𝗿𝗿𝟱 != nil {
		return 𝗲𝗿𝗿𝟱
	}
	𝗲𝗿𝗿𝟲 := w.Close()
	if 𝗲𝗿𝗿𝟲 != nil {
		return 𝗲𝗿𝗿𝟲
	}
	return nil
}
```

## Translating

```golang
package test

func func1[A, R any](a A) (_ R, _ error) { return }

func return1[A any]() (_ A)                          { return }
func ret0Err() (_ error)                             { return }
func ret1Err[A any]() (_ A, _ error)                 { return }
func ret2Err[A, B any]() (_ A, _ B, _ error)         { return }
func ret3Err[A, B, C any]() (_ A, _ B, _ C, _ error) { return }

func consume1[A any](A)             {}
func consume2[A, B any](A, B)       {}
func consume3[A, B, C any](A, B, C) {}

func id[A any](a A) A { return a }

func handleErrorf(err *error, format string, args ...interface{}) {
	if *err != nil {
		*err = fmt.Errorf(format+": %v", append(args, *err)...)
	}
}
```

<table>
<tr>
<td> 

</td>
<td> 

**Source**

</td> 
<td>

**Target**

</td>
</tr>



<tr>
<td> 

**Error Handling**

</td>
<td>

```golang
func error_wrapping() (a int, err error) {
	defer handleErrorf(&err, "something wrong")
	Try(ret1Err[bool]())
	a = 42
	return
}
```

</td>
<td>

```golang
func error_wrapping() (a int, err error) {
    defer handleErrorf(&err, "something wrong")
    _, 𝗲𝗿𝗿𝟭 := ret1Err[bool]()
    if 𝗲𝗿𝗿𝟭 != nil {
        err = 𝗲𝗿𝗿𝟭
        return
    }
    a = 42
    return
}
```

</td>
</tr>


<tr>
<td> 

**For Stmt**

</td>
<td>

```golang
for i := Try(ret1Err[A]()); 
    Try(func1[int, bool](i)); 
    Try(func1[A, C](i)) 
{
    println(i)
}
```

</td>
<td>

```golang
𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟭 := ret1Err[A]()
if 𝗲𝗿𝗿𝟭 != nil {
	return 𝗲𝗿𝗿𝟭
}
for i := 𝘃𝗮𝗹𝟭; ; {
	𝘃𝗮𝗹𝟮, 𝗲𝗿𝗿𝟮 := func1[int, bool](i)
	if 𝗲𝗿𝗿𝟮 != nil {
		return 𝗲𝗿𝗿𝟮
	}
	if !𝘃𝗮𝗹𝟮 {
		break
	}
	println(i)
	_, 𝗲𝗿𝗿𝟯 := func1[A, C](i)
	if 𝗲𝗿𝗿𝟯 != nil {
		return 𝗲𝗿𝗿𝟯
	}
}
```

</td>
</tr>


<tr>
<td> 

**Select Stmt**

</td>
<td>

```golang
type (
	A = int
	B = int
	C = int
	D = int
	E = int
	F = int
	G = int
	H = int
)
select {
    case <-Try(ret1Err[chan A]()):
    case *Try(ret1Err[*B]()), *Try(ret1Err[*bool]()) = <-Try(ret1Err[chan C]()):
    case Try(ret1Err[chan D]()) <- Try(ret1Err[E]()):
    case Try(ret1Err[[]F]())[Try(ret1Err[G]())] = <-Try(ret1Err[chan H]()):
    default:
}
```

</td>
<td>

```golang
type (
	A = int
	B = int
	C = int
	D = int
	E = int
	F = int
	G = int
	H = int
)
𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟭 := ret1Err[chan A]()
if 𝗲𝗿𝗿𝟭 != nil {
	return 𝗲𝗿𝗿𝟭
}
𝘃𝗮𝗹𝟮, 𝗲𝗿𝗿𝟮 := ret1Err[chan C]()
if 𝗲𝗿𝗿𝟮 != nil {
	return 𝗲𝗿𝗿𝟮
}
𝘃𝗮𝗹𝟲, 𝗲𝗿𝗿𝟱 := ret1Err[chan D]()
if 𝗲𝗿𝗿𝟱 != nil {
	return 𝗲𝗿𝗿𝟱
}
𝘃𝗮𝗹𝟳, 𝗲𝗿𝗿𝟲 := ret1Err[E]()
if 𝗲𝗿𝗿𝟲 != nil {
	return 𝗲𝗿𝗿𝟲
}
𝘃𝗮𝗹𝟴, 𝗲𝗿𝗿𝟳 := ret1Err[chan H]()
if 𝗲𝗿𝗿𝟳 != nil {
	return 𝗲𝗿𝗿𝟳
}
select {
    case <-𝘃𝗮𝗹𝟭:
    case 𝘃𝗮𝗹𝟰, 𝗼𝗸𝟭 := <-𝘃𝗮𝗹𝟮:
        𝘃𝗮𝗹𝟯, 𝗲𝗿𝗿𝟯 := ret1Err[*B]()
        if 𝗲𝗿𝗿𝟯 != nil {
            return 𝗲𝗿𝗿𝟯
        }
        *𝘃𝗮𝗹𝟯 = 𝘃𝗮𝗹𝟰
        𝘃𝗮𝗹𝟱, 𝗲𝗿𝗿𝟰 := ret1Err[*bool]()
        if 𝗲𝗿𝗿𝟰 != nil {
            return 𝗲𝗿𝗿𝟰
        }
        *𝘃𝗮𝗹𝟱 = 𝗼𝗸𝟭
    case 𝘃𝗮𝗹𝟲 <- 𝘃𝗮𝗹𝟳:
    case 𝘃𝗮𝗹𝟭𝟭 := <-𝘃𝗮𝗹𝟴:
        𝘃𝗮𝗹𝟵, 𝗲𝗿𝗿𝟴 := ret1Err[[]F]()
        if 𝗲𝗿𝗿𝟴 != nil {
            return 𝗲𝗿𝗿𝟴
        }
        𝘃𝗮𝗹𝟭𝟬, 𝗲𝗿𝗿𝟵 := ret1Err[G]()
        if 𝗲𝗿𝗿𝟵 != nil {
            return 𝗲𝗿𝗿𝟵
        }
        𝘃𝗮𝗹𝟵[𝘃𝗮𝗹𝟭𝟬] = 𝘃𝗮𝗹𝟭𝟭
    default:
}
```

</td>
</tr>

</table>