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
	"path"
	"path/filepath"
	"strings"
)

var dir = filepath.Clean(".out/generate/bindings")

func main() {
	log.Printf("Removing %q", dir)
	err := os.RemoveAll(dir)
	if err != nil {
		err = fmt.Errorf("failed to remove %q: %w", dir, err)
		log.Fatal(err)
	}

	log.Printf("Creating %q", dir)
	err = os.MkdirAll(filepath.Dir(dir), 0755)
	if err != nil {
		err = fmt.Errorf("failed to create directory %q: %w", dir, err)
		log.Fatal(err)
	}

	cmd := exec.Command("go", "tool", "wit-bindgen-go", "generate", "--out", dir, "wit")
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	log.Printf("Running %q", cmd)
	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("failed to run command %q: %w", cmd, err)
		log.Fatal(err)
	}

	for _, name := range []string{"error", "poll", "streams"} {
		log.Printf("Removing %q", name)
		err = os.RemoveAll(name)
		if err != nil {
			err = fmt.Errorf("failed to remove %q: %w", name, err)
			log.Fatal(err)
		}

		log.Printf("Renaming %q", name)
		err = os.Rename(filepath.Join(dir, "wasi/io", name), name)
		if err != nil {
			err = fmt.Errorf("failed to rename %q: %w", name, err)
			log.Fatal(err)
		}

		log.Printf("Removing %q", filepath.Join(name, "empty.s"))
		err = os.Remove(filepath.Join(name, "empty.s"))
		if err != nil {
			err = fmt.Errorf("failed to remove %q: %w", filepath.Join(name, "empty.s"), err)
			log.Fatal(err)
		}

		err = filepath.WalkDir(name, func(path2 string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
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
			code = bytes.Replace(code, []byte("DO NOT EDIT.\n"), []byte("DO NOT EDIT.\n//go:build wasip2\n"), 1)
			code = bytes.ReplaceAll(code, []byte(path.Join(dir, "wasi/io")+"/"), nil)
			if filepath.Base(filepath.Dir(path2)) == "error" {
				code = bytes.ReplaceAll(code, []byte("package ioerror"), []byte("package error"))
			}
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

	log.Printf("Removing %q", dir)
	err = os.RemoveAll(dir)
	if err != nil {
		err = fmt.Errorf("failed to remove %q: %w", dir, err)
		log.Fatal(err)
	}
}
