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

To get started, use `wit-bindgen-go` like normal. [The Bytecode Alliance's component model docs](https://component-model.bytecodealliance.org/) have [a nice guide to get started](https://component-model.bytecodealliance.org/language-support/go.html).

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

😵 Confused? See [`_examples/hello-world`](https://github.com/jcbhmr/go-wasi-io/tree/main/_examples/hello-world) for an example project that uses this package.

You'll likely want to create a script somewhere so that you can `go generate` to regenerate the bindings easily as you edit your `wit/*.wit` files.

<details><summary>Example <code>generate.go</code> file</summary>

```go
//go:build generate

//go:generate go run $GOFILE

package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	log.Printf("Removing %q", "internal")
	err := os.RemoveAll("internal")
	if err != nil {
		log.Fatalf("failed to remove directory %q: %v", "internal", err)
	}

	cmd := exec.Command("go", "tool", "wit-bindgen-go", "generate", "--out", "internal", "wit")
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	log.Printf("Running %q", cmd)
	err = cmd.Run()
	if err != nil {
		log.Fatalf("failed to run command %q: %v", cmd, err)
	}

	log.Printf("Removing %q", "internal/wasi/io")
	err = os.RemoveAll("internal/wasi/io")
	if err != nil {
		log.Fatalf("failed to remove directory %q: %v", "internal/wasi/io", err)
	}

	err = filepath.WalkDir("internal", func(path2 string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path2, ".go") {
			log.Printf("Skipping non-Go file %q", path2)
			return nil
		}
		code, err := os.ReadFile(path2)
		if err != nil {
			return fmt.Errorf("failed to read file %q: %w", path2, err)
		}
		code = bytes.ReplaceAll(code, []byte("<package-root>/internal/wasi/io/"), []byte("github.com/jcbhmr/go-wasi-io/v0.2.0/"))
		//                                    👆 Replace this with your Go module path
		log.Printf("Writing %q", path2)
		err = os.WriteFile(path2, code, 0644)
		if err != nil {
			return fmt.Errorf("failed to write file %q: %w", path2, err)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
```

**💡 VS Code tip:** Use `.vscode/settings.json` to mark all `//go:build generate` files as standalone Go scripts.

<div><code>.vscode/settings.json</code></div>

```json
{
    "gopls": {
        "build.standaloneTags": [
            "ignore",
            "generate"
        ]
    }
}
```

</details>

## Development

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=Go&logoColor=FFFFFF)
![WebAssembly](https://img.shields.io/badge/WebAssembly-654FF0?style=for-the-badge&logo=WebAssembly&logoColor=FFFFFF)

**Why `/v0.2.0` instead of `/v0` or `/v0.2`?** \
Normally, yes, `/v0.2` would be more "proper" since semver states that v0.2 ➡ v0.3 can include breaking (API, behaviour, etc.) changes. Then `/v2`, `/v3`, and so on. However, that limits your precision when choosing a specific version of `wasi:io` to whatever the common denominator is between your `/v0.2` and all your dependencies' `/v0.2` requirements. That's not precise enough for `wasi:io@0.2.0` in WIT. This is because a module path (like `github.com/jcbhmr/go-wasi-io/v0.2`) can only resolve to **a single version** like `github.com/jcbhmr/go-wasi-io/v0.2@v0.2.7` which might not be what you want (what if I really want the old `@v0.2.0` version?). **TL;DR:** using `/vX.Y.Z` lets you, the library consumer, pin an exact version.

Go doesn't yet (that I know of) have a WebAssembly runtime that supports running WASM Components and has a WASI preview 2 environment. This project is tested using other toolchains to inspect and test WASM Components.
