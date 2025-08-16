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

var bindingsDir = filepath.Clean(".out/generate/bindings")

func main() {
	log.Printf("Removing %q", bindingsDir)
	err := os.RemoveAll(bindingsDir)
	if err != nil {
		log.Fatalf("failed to remove %q: %v", bindingsDir, err)
	}

	log.Printf("Creating %q", bindingsDir)
	err = os.MkdirAll(filepath.Dir(bindingsDir), 0755)
	if err != nil {
		log.Fatalf("failed to create directory %q: %v", bindingsDir, err)
	}

	cmd := exec.Command("go", "tool", "wit-bindgen-go", "generate", "--out", bindingsDir, "wit")
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	log.Printf("Running %q", cmd)
	err = cmd.Run()
	if err != nil {
		log.Fatalf("failed to run command %q: %v", cmd, err)
	}

	for _, name := range []string{"error", "poll", "streams"} {
		log.Printf("Removing %q", name)
		err = os.RemoveAll(name)
		if err != nil {
			log.Fatalf("failed to remove all %q: %v", name, err)
		}

		log.Printf("Renaming %q to %q", filepath.Join(bindingsDir, "wasi/io", name), name)
		err = os.Rename(filepath.Join(bindingsDir, "wasi/io", name), name)
		if err != nil {
			err = fmt.Errorf("failed to rename %q to %q: %w", filepath.Join(bindingsDir, "wasi/io", name), name, err)
			log.Fatal(err)
		}

		log.Printf("Removing %q", filepath.Join(name, "empty.s"))
		err = os.Remove(filepath.Join(name, "empty.s"))
		if err != nil {
			err = fmt.Errorf("failed to remove all %q: %w", filepath.Join(name, "empty.s"), err)
			log.Fatal(err)
		}

		err = filepath.WalkDir(name, func(path2 string, d fs.DirEntry, err error) error {
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
			code = bytes.Replace(code, []byte("DO NOT EDIT.\n"), []byte("DO NOT EDIT.\n\n//go:build wasip2\n"), 1)
			code = bytes.ReplaceAll(code, []byte(path.Join(bindingsDir, "wasi/io")+"/"), nil)
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

	log.Printf("Removing %q", bindingsDir)
	err = os.RemoveAll(bindingsDir)
	if err != nil {
		err = fmt.Errorf("failed to remove all %q: %w", bindingsDir, err)
		log.Fatal(err)
	}
}
