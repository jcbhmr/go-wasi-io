package main_test

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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
	err := os.WriteFile(filepath.Join(d, "package.json"), []byte(`{"type":"module","dependencies":{"@bytecodealliance/preview2-shim":"^0.17.2"}}`), 0644)
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("npm", "install")
	cmd.Dir = d
	output, err := cmd.CombinedOutput()
	if len(output) > 0 {
		t.Logf("%s", output)
	}
	if err != nil {
		err = fmt.Errorf("failed to run command %q: %w", cmd, err)
		t.Fatal(err)
	}
	cmd = exec.Command("jco", "transpile", wasmPath, "--out-dir", d, "--name", "hello-world")
	t.Logf("Running %q", cmd)
	output, err = cmd.CombinedOutput()
	if len(output) > 0 {
		t.Logf("%s", output)
	}
	if err != nil {
		err = fmt.Errorf("failed to run command %q: %w", cmd, err)
		t.Fatal(err)
	}
	jsPathJSON, err := json.Marshal(filepath.Join(d, "hello-world.js"))
	if err != nil {
		t.Fatal(err)
	}
	cmd = exec.Command("node", "--eval", `import * as m from `+string(jsPathJSON)+`;
console.log(m);
m.greetings.sayHello("Alan Turing");
m.greetings.sayHi("Ada Lovelace");
m.greetings.sayHiInFrench("Charles Babbage");
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
	if !strings.Contains(string(output), "Alan Turing") {
		t.Errorf("expected to contain %q, got %q", "Alan Turing", string(output))
	}
	if !strings.Contains(string(output), "Ada Lovelace") {
		t.Errorf("expected to contain %q, got %q", "Ada Lovelace", string(output))
	}
	if !strings.Contains(string(output), "Charles Babbage") {
		t.Errorf("expected to contain %q, got %q", "Charles Babbage", string(output))
	}
}

// TODO: Move this build process to _scripts/build/main.go?
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
