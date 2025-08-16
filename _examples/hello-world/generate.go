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
		code = bytes.ReplaceAll(code, []byte("github.com/jcbhmr/go-wasi-io/v0.2.0/_examples/hello-world/internal/wasi/io/"), []byte("github.com/jcbhmr/go-wasi-io/v0.2.0/"))
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