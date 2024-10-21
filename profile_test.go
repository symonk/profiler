package profiler

import (
	"bytes"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// CheckFunc encapsulates the file descriptor and exit code
// for standard output and standard error.
type CheckFunc func(t *testing.T, stdout, stderr string, exit int)

func TestProfilesEnabledExpectedOutput(t *testing.T) {
	tests := map[string]struct {
		source string
		checks []CheckFunc
	}{
		"cpu enabled successfully (verbose)": {
			source: `package main
import "github.com/symonk/profiler"

func main() {
	defer profiler.Start(profiler.WithCPUProfiler()).Stop()
}`,
			checks: []CheckFunc{
				exitedZero,
				emptyStdErr,
				stdoutMatchesLinesInOrder("CPU profiling enabled"),
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Create a temporary file with the source code
			f := createTempTestFile(t, tc.source)
			defer os.Remove(f.Name())
			// Execute the program, capturing meta data
			stdout, stderr, exit := executeCase(f)
			// Assert the output is as expected
			for _, checkFunc := range tc.checks {
				checkFunc(t, stdout, stderr, exit)
			}
		})
	}
}

// executeCase executes the source code written earlier and captures
// its stdout, stderr and exit code for inspection later.
func executeCase(f *os.File) (string, string, int) {
	cmd := exec.Command("go", "run", f.Name())
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
func createTempTestFile(t *testing.T, source string) *os.File {
	runnable, err := os.CreateTemp(os.TempDir(), "profiler_t")
	assert.NoError(t, err)
	_, err = runnable.Write([]byte(source))
	assert.NoError(t, err)
	return runnable
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

func emptyStdErr(t *testing.T, _, stderr string, _ int) {
	assert.Empty(t, stderr)
}
