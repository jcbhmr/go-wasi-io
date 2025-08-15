# WASI IO bindings for Go

## Installation

```sh
go get github.com/jcbhmr/go-wasi-io/v0.2.0-rc1
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
    └── wasi/ # Dedupe these wasi:io bindings from your project
        └── io/
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