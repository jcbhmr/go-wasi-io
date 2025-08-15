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
		err = fmt.Errorf("failed to remove directory %q: %w", "internal", err)
		log.Fatal(err)
	}

	cmd := exec.Command("go", "tool", "wit-bindgen-go", "generate", "--out", "internal", "wit")
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	log.Printf("Running %q", cmd)
	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("failed to run command %q: %w", cmd, err)
		log.Fatal(err)
	}

	log.Printf("Removing %q", "internal/wasi/io")
	err = os.RemoveAll("internal/wasi/io")
	if err != nil {
		log.Fatal(err)
	}

	err = filepath.WalkDir("internal", func(path2 string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk dir %q: %w", path2, err)
		}
		if d.IsDir() {
			return nil
		}
		if !d.Type().IsRegular() {
			log.Printf("Skipping non-regular file %q", path2)
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
		code = bytes.ReplaceAll(code, []byte("github.com/jcbhmr/go-wasi-io/v0.2/_examples/hello-world/internal/wasi/io/"), []byte("github.com/jcbhmr/go-wasi-io/v0.2/"))
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
