package profiler

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// CheckFunc encapsulates the file descriptor and exit code
// for standard output and standard error.
type CheckFunc func(t *testing.T, stdout, stderr string, exit int)

func TestProfilesEnabledExpectedOutput(t *testing.T) {
	storage, err := os.MkdirTemp("", "profiles")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(storage)
	tests := map[string]struct {
		source string
		checks []CheckFunc
	}{
		"cpu profiling works as expected": {
			source: `package main
import "github.com/symonk/profiler"

func main() {
	defer profiler.Start(profiler.WithCPUProfiler(), profiler.WithProfileFileLocation("` + storage + "\"" + `)).Stop()
}`,
			checks: []CheckFunc{
				exitedZero,
				emptyStdOut,
				stdErrOutMatchLines(
					".*profiling completed.  You can find the .*cpu.pprof.*",
					".*to view the profile, run.*cpu.pprof",
					"port can be any ephemeral port you wish to use",
					"Graph interpretation is outlined here.*graphical-reports",
				),
			},
		},
	}
	for name, tc := range tests {
		t.Log(name)
		t.Run(name, func(t *testing.T) {
			// Execute the program, capturing meta data
			stdout, stderr, exit := execute(t, tc.source)
			t.Log(stdout, stderr)
			// Assert the output is as expected
			for _, checkFunc := range tc.checks {
				checkFunc(t, stdout, stderr, exit)
			}
		})
	}
}

// execute executes the source code written earlier and captures
// its stdout, stderr and exit code for inspection later.
func execute(t *testing.T, source string) (string, string, int) {
	main, cleanup := createTempTestFile(t, source)
	defer cleanup()
	cmd := exec.Command("go", "run", main)
	var stdOut, stdErr bytes.Buffer
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return stdOut.String(), stdErr.String(), exitErr.ExitCode()
		}
	}
	return stdOut.String(), stdErr.String(), 0
}

// createTempTestFile creates a temporary test file with the source
// code ready for running.
func createTempTestFile(t *testing.T, source string) (string, func()) {
	fatal := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}
	dir, err := os.MkdirTemp("", "go-profiler")
	fatal(err)

	main := filepath.Join(dir, "main.go")

	// Create the source code file
	err = os.WriteFile(main, []byte(source), 0644)
	fatal(err)

	// Create the go mod file so we can actually execute!
	// TODO: This uses a hardcoded dependency version.
	mod := filepath.Join(dir, "go.mod")
	contents := `
	module github.com/ok/example
go 1.23.2
require github.com/symonk/profiler v0.0.0-20241021143805-788e1dbe92a9
`
	os.WriteFile(mod, []byte(contents), 0644)

	// create the appropriate mod file etc
	return main, func() { defer os.RemoveAll(dir) }
}

// Check function implementations for asserting against the responses
func exitedZero(t *testing.T, _, _ string, code int) {
	assert.Zero(t, code)
}

// patternMatchLines checks that the lines in stdout matched
func stdOutOutMatchLines(patterns ...string) CheckFunc {
	return func(t *testing.T, stdout, stderr string, exit int) {
		assert.NotEmpty(t, stdout)
		patternMatchLines(t, stdout, patterns...)
	}
}

// patternMatchLines checks that the lines in stderr matched
func stdErrOutMatchLines(patterns ...string) CheckFunc {
	return func(t *testing.T, stdout, stderr string, exit int) {
		assert.NotEmpty(t, stderr)
		patternMatchLines(t, stderr, patterns...)
	}
}

// patternMatchLines checks that the lines in either stdout/err matched
// the user provided regexp patterns.  No order is guarantee'd here and
// all are iterated for each pattern, this is not very performant and can
// be done in O(n) in future most likely.
func patternMatchLines(t *testing.T, input string, patterns ...string) bool {
	seen := make(map[string]struct{}, len(patterns))
	for _, p := range patterns {
		seen[p] = struct{}{}
	}

	lines := strings.Split(input, "\n")
	for i := 0; i < len(lines); i++ {
		for j := 0; j < len(patterns); j++ {
			if matched, _ := regexp.MatchString(patterns[j], lines[i]); matched {
				delete(seen, patterns[j])
			}
		}
	}

	if len(seen) == 0 {
		return true
	}

	t.Fatalf("expected all patterns to be matched, but the following were not: %v", seen)
	return false
}

func emptyStdErr(t *testing.T, _, stderr string, _ int) {
	assert.Empty(t, stderr)
}

func emptyStdOut(t *testing.T, stdout, _ string, _ int) {
	assert.Empty(t, stdout)
}
