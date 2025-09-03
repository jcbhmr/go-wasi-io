# `wasi:io` bindings for Go

📂 Centralized bindings to [`wasi:io`](https://github.com/WebAssembly/wasi-io) interfaces

<table align=center>
<td>

```
.
└── internal/
    ├── octocat/
    │   └── my-app/
    │       └── my-interface/
    │           └── ...
    └── wasi/
        ├── io/ 👈 Replaces this folder
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

</table>

✂️ Eliminate duplicate code; use a centralized pregenerated bindings package

## Installation

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=Go&logoColor=FFFFFF)

```sh
go get github.com/jcbhmr/go-wasi-io
```

⚠️ The latest version is v0.2.7. You probably want v0.2.0. Use @v0.2.0 to select it.

```sh
go get github.com/jcbhmr/go-wasi-io@v0.2.0
```

## Usage

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=Go&logoColor=FFFFFF)
![WebAssembly](https://img.shields.io/badge/WebAssembly-654FF0?style=for-the-badge&logo=WebAssembly&logoColor=FFFFFF)

```go
//go:generate go tool wit-bindgen-go generate --out ./internal/ ./wit/
//go:generate rm -rf ./internal/wasi/io/
//go:generate go tool jet -g "*.go" "<your-package-root>/internal/wasi/io/" "github.com/jcbhmr/go-wasi-io/" ./internal/
```

```json
{
    "go.buildTags": "wasip2"
}
```

## Development

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=Go&logoColor=FFFFFF)
![WebAssembly](https://img.shields.io/badge/WebAssembly-654FF0?style=for-the-badge&logo=WebAssembly&logoColor=FFFFFF)
