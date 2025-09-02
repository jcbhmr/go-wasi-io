//go:generate rm -rf ./.out/bindings/ ./error/ ./poll/ ./streams/
//go:generate go tool wit-bindgen-go generate --out ./.out/bindings/ ./wit/
//go:generate mv ./.out/bindings/wasi/io/error/ ./error/
//go:generate mv ./.out/bindings/wasi/io/poll/ ./poll/
//go:generate mv ./.out/bindings/wasi/io/streams/ ./streams/
//go:generate rm -rf ./error/empty.s ./poll/empty.s ./streams/empty.s
//go:generate go tool jet -g "*.go" -e "DO NOT EDIT.\n" "DO NOT EDIT.\n\n//go:build wasip2\n" -e "\\.out/bindings/wasi/io/" "" ./error/ ./poll/ ./streams/

package io_test
