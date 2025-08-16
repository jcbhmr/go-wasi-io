//go:build wasip2

package main

import (
	"fmt"

	"github.com/jcbhmr/go-wasi-io/v0.2.0/_examples/hello-world/internal/jcbhmr/hello-world/greetings"
)

func init() {
	greetings.Exports.SayHello = func(name string) {
		fmt.Printf("Hello %s!\n", name)
	}
	greetings.Exports.SayHi = func(name string) {
		fmt.Printf("Hi %s!\n", name)
	}
	greetings.Exports.SayHiInFrench = func(name string) {
		fmt.Printf("Salut %s!\n", name)
	}
}

func main() {}
