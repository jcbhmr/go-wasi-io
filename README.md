# `wasi:io` bindings for Go

📂 Centralized bindings to `wasi:io` interfaces

<table align=center>
<tr>
<th>Before
<th>After
<tr>
<td>

```
.
└── internal/
    ├── octocat/
    │   └── my-app/
    │       └── my-interface/
    │           └── ...
    └── wasi/
        ├── io/
        │   ├── error/
        │   │   ├── empty.s
        │   │   ├── error.wasm.go
        │   │   └── error.wit.go
        │   ├── poll/
        │   │   ├── empty.s
        │   │   ├── poll.wasm.go
        │   │   └── poll.wit.go
        │   └── streams/
        │       ├── empty.s
        │       ├── streams.wasm.go
        │       └── streams.wit.go
        └── ...
```

<td>

```
.
└── internal/
    └── octocat/
        └── my-app/
            └── my-interface/
                └── ...
```

```go
import (
    "github.com/jcbhmr/go-wasi-io/v0.2.0/error"
    "github.com/jcbhmr/go-wasi-io/v0.2.0/poll"
    "github.com/jcbhmr/go-wasi-io/v0.2.0/streams"
)
```

</table>

✂️ Eliminate duplicate code; use a centralized pregenerated bindings package

## Installation

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=Go&logoColor=FFFFFF)

```sh
go get github.com/jcbhmr/go-wasi-io/v0.2.0
```

## Usage

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=Go&logoColor=FFFFFF)
![WebAssembly](https://img.shields.io/badge/WebAssembly-654FF0?style=for-the-badge&logo=WebAssembly&logoColor=FFFFFF)

To get started, use `wit-bindgen-go` like normal. [BytecodeAlliance's component model docs](https://component-model.bytecodealliance.org/) have [a nice guide to get started](https://component-model.bytecodealliance.org/language-support/go.html).

```wit
package octocat:hello;
interface greetings {
    say-hello: func(name: string);
}
world hello {
    include wasi:cli/imports@0.2.0;
}
```

```sh
go tool wit-bindgen-go generate --out internal wit
```

You should now have a bunch of **local** generated bindings. The `.../wasi/io/*` packages are the ones that this repository centralizes here.

```
.
└── internal/
    ├── octocat/
    │   └── my-app/
    │       └── my-interface/
    │           └── ...
    └── wasi/
        ├── io/
        │   ├── error/
        │   │   ├── empty.s
        │   │   ├── error.wasm.go
        │   │   └── error.wit.go
        │   ├── poll/
        │   │   ├── empty.s
        │   │   ├── poll.wasm.go
        │   │   └── poll.wit.go
        │   └── streams/
        │       ├── empty.s
        │       ├── streams.wasm.go
        │       └── streams.wit.go
        └── ...
```

1. Remove the `internal/wasi/io/` folder completely.
2. Find and replace all instances of `<package-root>/internal/wasi/io/` with `github.com/jcbhmr/go-wasi-io/v0.2.0/`

**💡 Pro tip:** Use your IDE's find-and-replace feature to replace everything across many files at once

😵 Confused? See [`_examples/hello-world`](https://github.com/jcbhmr/go-wasi-io/tree/main/_examples/hello-world) for an example project that uses this package. The [`generate.go`](https://github.com/jcbhmr/go-wasi-io/blob/main/_examples/hello-world/generate.go) script does the post-processing of the `wit-bindgen-go`-generated bindings.

You'll likely want to create a script somewhere so that you can `go generate` to regenerate the bindings easily as you edit your `wit/*.wit` files.

## Development

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=Go&logoColor=FFFFFF)
![WebAssembly](https://img.shields.io/badge/WebAssembly-654FF0?style=for-the-badge&logo=WebAssembly&logoColor=FFFFFF)

**Why `/v0.2.0` instead of `/v0` or `/v0.2`?** \
Normally, yes, `/v0.2` would be more "proper" since semver states that v0.2 ➡ v0.3 can include breaking (API, behaviour, etc.) changes. But then there'd be no way to pin a specific version of this package properly. Go doesn't support version pinning in `go.mod` because the entire module system is based on module identifiers being globally unique in a program. This means that if `/v0.2` were used, then two packages could not independently depend on `wasi:io/poll/pollable@0.2.0` and `wasi:io/poll/pollable@0.2.7` in the same program because both would resolve to the single module path `github.com/jcbhmr/go-wasi-io/v0.2`. Using `/vX.Y.Z` means that each version has its own completely independent module path and multiple versions can be present in the same program, just like `/v2` normally.

Go doesn't yet (that I know of) have a WebAssembly runtime that supports running WASM Components and has a WASI preview 2 environment.
