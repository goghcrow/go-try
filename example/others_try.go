//go:build try

//go:generate go install github.com/goghcrow/go-try/cmd/trygen@main
//go:generate trygen

package example

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	. "github.com/goghcrow/go-try"
)

func printSum1(a, b string) error {
	x := Try(strconv.Atoi(a))
	y := Try(strconv.Atoi(b))
	fmt.Println("result:", x+y)
	return nil
}

func printSum2(a, b string) error {
	fmt.Println(
		"result:",
		Try(strconv.Atoi(a))+Try(strconv.Atoi(b)),
	)
	return nil
}

func parseHexdump(text string) (_ []byte, _ error) { return }

func localMain() error {
	hex := Try(io.ReadAll(os.Stdin))
	data := Try(parseHexdump(string(hex)))
	Try(os.Stdout.Write(data))
	return nil
}

func main() {
	if err := localMain(); err != nil {
		log.Fatal(err)
	}
}
