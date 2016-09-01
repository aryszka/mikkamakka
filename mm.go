package main

import (
	"os"
	"io"
	"bufio"
	"log"
)

type typeCheck func(either) bool

type either struct {
	value interface{}
	types []typeCheck
}

func fopen(u string) *either {
	return nil
}

func main() {
	var (
		input []byte
		escaped bool
	)

	r := bufio.NewReader(os.Stdin)
	for {
		b, err := r.ReadByte()

		if err == io.EOF {
			return
		}

		if err != nil {
			log.Fatal(err)
		}

		appendEscaped := func(b byte) bool {
			if !escaped {
				return false
			}

			input = append(input, b)
			escaped = false
			return true
		}

		switch b {
		case '\\':
			if !appendEscaped(b) {
				escaped = true
			}
		case '\n':
			if !appendEscaped(b) && len(input) > 0 {
				println(string(input))
				input = nil
			}
		default:
			input = append(input, b)
		}
	}
}
