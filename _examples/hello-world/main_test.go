package main_test

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
)

var wasmPath string

func TestInspectWIT(t *testing.T) {
	cmd := exec.Command("wasm-tools", "component", "wit", wasmPath)
	t.Logf("Running %q", cmd)
	output, err := cmd.CombinedOutput()
	if len(output) > 0 {
		t.Logf("%s", output)
	}
	if err != nil {
		err = fmt.Errorf("failed to run command %q: %w", cmd, err)
		t.Fatal(err)
	}
}

func TestRunJS(t *testing.T) {
	d := t.TempDir()
	cmd := exec.Command("jco", "--help", d)
	t.Logf("Running %q", cmd)
	output, err := cmd.CombinedOutput()
	if len(output) > 0 {
		t.Logf("%s", output)
	}
	if err != nil {
		err = fmt.Errorf("failed to run command %q: %w", cmd, err)
		t.Fatal(err)
	}
	cmd = exec.Command("node", "--input-type", "module", "--eval", `
		console.log(42)
	`)
	t.Logf("Running %q", cmd)
	output, err = cmd.CombinedOutput()
	if len(output) > 0 {
		t.Logf("%s", output)
	}
	if err != nil {
		err = fmt.Errorf("failed to run command %q: %w", cmd, err)
		t.Fatal(err)
	}
}

func TestMain(m *testing.M) {
	exe, err := os.Executable()
	if err != nil {
		err = fmt.Errorf("failed to get executable path: %w", err)
		log.Fatal(err)
	}
	witPath := strings.TrimSuffix(exe, ".exe") + ".wit.wasm"
	cmd := exec.Command("wkg", "wit", "build", "--output", witPath)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	log.Printf("Running %q", cmd)
	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("failed to run command %q: %w", cmd, err)
		log.Fatal(err)
	}
	wasmPath = strings.TrimSuffix(exe, ".exe") + ".wasm"
	cmd = exec.Command("tinygo", "build", "-target", "wasip2", "-o", wasmPath, "-wit-package", witPath, "-wit-world", "hello-world")
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	log.Printf("Running %q", cmd)
	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("failed to run command %q: %w", cmd, err)
		log.Fatal(err)
	}
	os.Exit(m.Run())
}
