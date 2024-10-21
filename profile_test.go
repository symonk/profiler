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
		"cpu enabled successfully (verbose)": {
			source: `package main
import "github.com/symonk/profiler"

func main() {
	defer profiler.Start(profiler.WithCPUProfiler(), WithProfileFileLocation("` + storage + "\"" + `)).Stop()
}`,
			checks: []CheckFunc{
				exitedZero,
				emptyStdOut,
				stdErrMatchesLinesInOrder(
					".*setting up cpu profiler.*",
					".*profile completed.  You can find the .*cpu.pprof",
					".*to view the profile, run.*cpu.pprof",
				),
			},
		},
	}
	for name, tc := range tests {
		t.Log(name)
		t.Run(name, func(t *testing.T) {
			// Execute the program, capturing meta data
			stdout, stderr, exit := execute(t, tc.source)
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

// stdoutMatchesLinesInOrder checks that the stdout matches the
// provided patterns in order and that all the patterns are matched.
// a utility function for checking the output is as expected
func stdoutMatchesLinesInOrder(patterns ...string) CheckFunc {
	return func(t *testing.T, stdout, _ string, _ int) {
		hits := len(patterns)
		stdoutLines := strings.Split(stdout, "\n")
		p0, p1 := 0, 0
		for p0 < len(patterns) && p1 < len(stdoutLines) {
			if matched, _ := regexp.MatchString(patterns[p0], stdoutLines[p1]); matched {
				p0++
				hits--
			}
			p1++
		}
		assert.Zero(t, hits, "expected all the lines provided to match lines in the std output")
	}
}

func stdErrMatchesLinesInOrder(patterns ...string) CheckFunc {
	return func(t *testing.T, _, stderr string, _ int) {
		hits := len(patterns)
		stderrLines := strings.Split(stderr, "\n")
		p0, p1 := 0, 0
		for p0 < len(patterns) && p1 < len(stderrLines) {
			if matched, _ := regexp.MatchString(patterns[p0], stderrLines[p1]); matched {
				p0++
				hits--
			}
			p1++
		}
		assert.Zero(t, hits, "expected all the lines provided to match lines in the std error")
	}
}

func emptyStdErr(t *testing.T, _, stderr string, _ int) {
	assert.Empty(t, stderr)
}

func emptyStdOut(t *testing.T, stdout, _ string, _ int) {
	assert.Empty(t, stdout)
}
