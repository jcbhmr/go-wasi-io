# WASI IO bindings for Go

## Installation

```sh
go get github.com/jcbhmr/go-wasi-io/v0.2
```

## Usage

Instead of relying on your local copy of generated bindings from wit-bindgen-go you can use this package.

```sh
.
└── internal/
    ├── octocat/
    │   └── my-app/
    │       └── my-interface/
    │           └── ...
    └── wasi/
        └── io/ # 👈 Dedupe these wasi:io bindings from your project
            ├── error/
            │   ├── empty.s
            │   ├── error.wasm.go
            │   └── error.wit.go
            ├── poll/
            │   ├── empty.s
            │   ├── poll.wasm.go
            │   └── poll.wit.go
            └── streams/
                ├── empty.s
                ├── streams.wasm.go
                └── streams.wit.go
```

See [`_examples/hello-world`](https://github.com/jcbhmr/go-wasi-io/tree/main/_examples/hello-world) for an example project that uses this package.
